package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/vladimish/talk/pkg/pointer"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/completion"
	"github.com/vladimish/talk/internal/port/queue"
	"github.com/vladimish/talk/internal/port/storage"
	"github.com/vladimish/talk/pkg/i18n"
)

const (
	processingLockTimeout     = 5 * time.Minute
	maxConversationNameLength = 50
	minConversationNameLength = 2
	queueProcessingDelay      = 100 * time.Millisecond
	typingIndicatorInterval   = 4500 * time.Millisecond // 4.5 seconds
)

// Conversation handlers.
func (s *UpdateService) handleConversationMessage(ctx context.Context, user *domain.User, update domain.Update) error {
	return s.handleConversationMessageWithReply(ctx, user, update, nil)
}

func (s *UpdateService) handleConversationMessageWithReply(
	ctx context.Context,
	user *domain.User,
	update domain.Update,
	replyToMessageID *int64,
) error {
	// Check if an image is provided and validate model support
	if len(update.ImageData) > 0 || update.ImageMimeType != "" {
		modelInfo := domain.GetModelByID(user.SelectedModel)
		if modelInfo == nil || !modelInfo.ImageSupport {
			errorMsg := i18n.GetString(user.Language, i18n.ModelImageNotSupported)
			_, sendErr := s.sender.SendMessage(ctx, user.ExternalID, errorMsg)
			if sendErr != nil {
				s.logger.WarnContext(
					ctx,
					"failed to send image not supported message",
					slog.String("error", sendErr.Error()),
				)
			}
			return nil // Don't process the message further
		}
	}

	// Get model info for cost calculation
	currentModel := domain.GetModelByID(user.SelectedModel)
	if currentModel == nil {
		return fmt.Errorf("model not found: %s", user.SelectedModel)
	}

	// Check if user has active subscription (required for web search)
	activeSubscription, err := s.storage.GetActiveSubscriptionByUserID(ctx, user.ID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return fmt.Errorf("failed to check subscription status: %w", err)
	}
	hasActiveSubscription := activeSubscription != nil

	// Check if web search is enabled and user has subscription
	webSearchEnabled := currentModel.WebSearch && user.WebSearchEnabled && hasActiveSubscription

	// Check base model token balance
	baseBalance, err := s.storage.GetUserTokenBalanceByType(ctx, user.ID, currentModel.TokenType)
	if err != nil {
		return fmt.Errorf("failed to get user base token balance: %w", err)
	}

	if baseBalance < currentModel.Cost {
		tokenTypeName := "regular"
		if currentModel.TokenType == domain.TokenTypePremium {
			tokenTypeName = "premium"
		}
		insufficientTokensMsg := fmt.Sprintf(
			i18n.GetString(user.Language, i18n.ProfileInsufficientTokens),
			currentModel.Cost,
			tokenTypeName,
		)
		_, sendErr := s.sender.SendMessage(ctx, user.ExternalID, insufficientTokensMsg)
		return sendErr
	}

	// Check search token balance if web search is enabled
	if webSearchEnabled && currentModel.SearchCost != nil && currentModel.SearchTokenType != nil {
		searchBalance, searchErr := s.storage.GetUserTokenBalanceByType(ctx, user.ID, *currentModel.SearchTokenType)
		if searchErr != nil {
			return fmt.Errorf("failed to get user search token balance: %w", searchErr)
		}

		if searchBalance < *currentModel.SearchCost {
			searchTokenTypeName := "regular"
			if *currentModel.SearchTokenType == domain.TokenTypePremium {
				searchTokenTypeName = "premium"
			}
			insufficientSearchTokensMsg := fmt.Sprintf(
				i18n.GetString(user.Language, i18n.ProfileInsufficientTokens),
				*currentModel.SearchCost,
				searchTokenTypeName,
			)
			_, sendErr := s.sender.SendMessage(ctx, user.ExternalID, insufficientSearchTokensMsg)
			return sendErr
		}
	}

	// Set processing lock with 5 minute timeout
	if lockErr := s.queue.SetProcessing(ctx, user.ExternalID, processingLockTimeout); lockErr != nil {
		if errors.Is(lockErr, queue.ErrAlreadyProcessing) {
			// This shouldn't happen as we check before, but handle it gracefully
			return s.queue.EnqueueWithNotification(ctx, user.ExternalID, update, "")
		}
		s.logger.WarnContext(ctx, "failed to set processing lock",
			slog.String("error", lockErr.Error()))
		// Continue without lock on error
	}

	// Ensure we clear the lock and process queued messages when done
	defer func() {
		if clearErr := s.queue.ClearProcessing(ctx, user.ExternalID); clearErr != nil {
			s.logger.ErrorContext(ctx, "failed to clear processing lock",
				slog.String("error", clearErr.Error()))
		}

		// Process any queued messages
		go s.processQueuedMessages(context.Background(), user)
	}()
	// Check if this is the first message in a new conversation
	var isFirstMessage bool
	if user.CurrentConversationID != nil {
		messages, msgErr := s.storage.GetMessagesByConversationID(ctx, *user.CurrentConversationID)
		if msgErr != nil {
			return fmt.Errorf("can't check existing messages: %w", msgErr)
		}
		isFirstMessage = len(messages) == 0
	}

	// Save incoming message from user
	userMessage, err := s.storage.CreateMessage(ctx, &domain.Message{
		UserID: user.ID,
		MessageType: domain.MessageType{
			Text:          update.MessageText,
			ImageData:     update.ImageData,
			ImageMimeType: update.ImageMimeType,
			PDFData:       update.PDFData,
			PDFMimeType:   update.PDFMimeType,
			PDFFileName:   update.PDFFileName,
		},
		SentBy:         domain.MessageSenderUser,
		ConversationID: user.CurrentConversationID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		return fmt.Errorf("can't save user message: %w", err)
	}

	// Update conversation timestamp when user message is created
	if user.CurrentConversationID != nil {
		if timestampErr := s.storage.UpdateConversationTimestamp(ctx, *user.CurrentConversationID); timestampErr != nil {
			s.logger.WarnContext(ctx, "failed to update conversation timestamp",
				slog.String("error", timestampErr.Error()),
				slog.Int64("conversation_id", *user.CurrentConversationID))
		}
	}

	// Generate conversation name if this is the first message
	if isFirstMessage && user.CurrentConversationID != nil {
		go s.generateAndUpdateConversationName(context.Background(), *user.CurrentConversationID, update.MessageText)
	}

	if update.ExternalMessageID > 0 {
		if userMessage.ID > int64(^uint32(0)>>1) || int64(update.ExternalMessageID) > int64(^uint32(0)>>1) {
			s.logger.WarnContext(ctx, "message ID too large for foreign message mapping")
		} else {
			err = s.storage.CreateForeignMessage(ctx, int32(userMessage.ID), int32(update.ExternalMessageID)) //nolint:gosec
			if err != nil {
				s.logger.WarnContext(ctx, "failed to save foreign message mapping", slog.String("error", err.Error()))
			}
		}
	}

	// Get conversation history
	var messages []*domain.Message
	if user.CurrentConversationID != nil {
		messages, err = s.storage.GetMessagesByConversationID(ctx, *user.CurrentConversationID)
		if err != nil {
			return fmt.Errorf("can't get messages: %w", err)
		}
	}

	// Send typing indicator and start periodic typing
	err = s.sender.SendTyping(ctx, user.ExternalID)
	if err != nil {
		s.logger.WarnContext(ctx, "failed to send typing indicator", slog.String("error", err.Error()))
	}

	// Start periodic typing indicator for the entire duration
	typingDone := make(chan struct{})
	defer close(typingDone)

	go s.sendPeriodicTyping(ctx, user.ExternalID, typingDone)

	// Handle image upload if present
	uploadResult := s.handleImageUpload(ctx, update.ImageData, update.ImageMimeType)
	currentImageURL := ""
	if uploadResult != nil {
		currentImageURL = uploadResult.imageURL
	}

	// Handle PDF upload if present
	pdfUploadResult := s.handlePDFUpload(ctx, update.PDFData, update.PDFMimeType, update.PDFFileName)
	currentPDFURL := ""
	currentPDFFileName := ""
	if pdfUploadResult != nil {
		currentPDFURL = pdfUploadResult.pdfURL
		currentPDFFileName = update.PDFFileName
	}

	// Create attachment record if image was uploaded
	if uploadResult != nil {
		_, attachmentErr := s.storage.CreateAttachment(ctx, &domain.Attachment{
			MessageID:   userMessage.ID,
			S3Name:      uploadResult.objectName,
			ContentType: uploadResult.contentType,
			Size:        uploadResult.size,
		})
		if attachmentErr != nil {
			s.logger.ErrorContext(
				ctx,
				"failed to create image attachment record",
				slog.String("error", attachmentErr.Error()),
			)
		}
	}

	// Create attachment record if PDF was uploaded
	if pdfUploadResult != nil {
		_, attachmentErr := s.storage.CreateAttachment(ctx, &domain.Attachment{
			MessageID:   userMessage.ID,
			S3Name:      pdfUploadResult.objectName,
			ContentType: pdfUploadResult.contentType,
			Size:        pdfUploadResult.size,
		})
		if attachmentErr != nil {
			s.logger.ErrorContext(
				ctx,
				"failed to create PDF attachment record",
				slog.String("error", attachmentErr.Error()),
			)
		}
	}

	// Check if PDF is provided but model doesn't support it
	if pdfUploadResult != nil && !currentModel.PDFSupport {
		notSupportedMsg := i18n.GetString(user.Language, i18n.ModelPDFNotSupported)
		_, sendErr := s.sender.SendMessage(ctx, user.ExternalID, notSupportedMsg)
		if sendErr != nil {
			s.logger.WarnContext(ctx, "failed to send PDF not supported message",
				slog.String("error", sendErr.Error()))
		}
		return nil
	}

	systemPrompt := "You are a helpful assistant."

	// Web search is already calculated above with subscription check
	// Use the webSearchEnabled variable from the cost calculation

	// Prepare file attachments if any
	var imageAttachment *completion.FileAttachment
	if currentImageURL != "" {
		imageAttachment = &completion.FileAttachment{
			URL:      currentImageURL,
			MimeType: "image/jpeg",
		}
	}

	var pdfAttachment *completion.FileAttachment
	if len(update.PDFData) > 0 {
		pdfAttachment = &completion.FileAttachment{
			URL:      currentPDFURL, // Keep the MinIO URL for other purposes if needed
			MimeType: "application/pdf",
			FileName: currentPDFFileName,
			Data:     update.PDFData, // Pass the raw PDF data for OpenRouter
		}
	}

	tokenStream, err := s.completion.CompleteStreamWithAttachments(
		ctx,
		user.SelectedModel,
		systemPrompt,
		messages,
		imageAttachment,
		pdfAttachment,
		webSearchEnabled,
	)
	if err != nil {
		return fmt.Errorf("can't get completion: %w", err)
	}

	// Stream response tokens and handle messaging
	var responseBuilder strings.Builder
	var reasoningBuilder strings.Builder
	var messageIDs []string
	var reasoningMessageIDs []string
	var previousContent string
	var previousReasoningContent string
	var hasReasoningMessage bool
	var hasMainMessage bool
	lastUpdate := time.Now()

	for token := range tokenStream {
		if token.Error != nil {
			return fmt.Errorf("completion stream error: %w", token.Error)
		}

		// Handle reasoning tokens - accumulate in reasoning builder
		if token.Reasoning != "" {
			reasoningBuilder.WriteString(token.Reasoning)
		}

		// Handle regular content
		if token.Content != "" {
			responseBuilder.WriteString(token.Content)
		}

		// Update messages at most once per second
		if time.Since(lastUpdate) >= time.Second {
			// Handle periodic updates
			shouldContinue, updateErr := s.handlePeriodicUpdates(
				ctx,
				user,
				&reasoningBuilder,
				&responseBuilder,
				&reasoningMessageIDs,
				&messageIDs,
				&previousReasoningContent,
				&previousContent,
				&hasReasoningMessage,
				&hasMainMessage,
				replyToMessageID,
			)
			if updateErr != nil {
				return updateErr
			}
			if shouldContinue {
				continue
			}

			lastUpdate = time.Now()
		}
	}

	// Send final updates if needed
	if time.Since(lastUpdate) > 0 {
		messageIDs, err = s.sendFinalUpdates(
			ctx,
			user,
			&responseBuilder,
			&reasoningBuilder,
			messageIDs,
			reasoningMessageIDs,
			previousContent,
			previousReasoningContent,
			hasReasoningMessage,
			hasMainMessage,
			replyToMessageID,
		)
		if err != nil {
			return err
		}
	}

	responseText := responseBuilder.String()

	botMessage, err := s.storage.CreateMessage(ctx, &domain.Message{
		UserID: user.ID,
		MessageType: domain.MessageType{
			Text: responseText,
		},
		SentBy:         domain.MessageSenderBot,
		ConversationID: user.CurrentConversationID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		return fmt.Errorf("can't save bot message: %w", err)
	}

	// Update conversation timestamp when bot message is created
	if user.CurrentConversationID != nil {
		if timestampErr := s.storage.UpdateConversationTimestamp(ctx, *user.CurrentConversationID); timestampErr != nil {
			s.logger.WarnContext(ctx, "failed to update conversation timestamp",
				slog.String("error", timestampErr.Error()),
				slog.Int64("conversation_id", *user.CurrentConversationID))
		}
	}

	// Save foreign message mappings for all bot messages
	for _, msgIDStr := range messageIDs {
		if msgIDStr == "" {
			continue
		}

		// Parse message ID from string to int
		msgID, parseErr := strconv.ParseInt(msgIDStr, 10, 32)
		if parseErr != nil {
			s.logger.WarnContext(ctx, "failed to parse message ID",
				slog.String("message_id", msgIDStr),
				slog.String("error", parseErr.Error()))
			continue
		}

		// Create foreign message mapping
		if botMessage.ID > int64(^uint32(0)>>1) || msgID > int64(^uint32(0)>>1) {
			s.logger.WarnContext(ctx, "message ID too large for foreign message mapping")
			continue
		}
		err = s.storage.CreateForeignMessage(ctx, int32(botMessage.ID), int32(msgID)) //nolint:gosec
		if err != nil {
			s.logger.WarnContext(ctx, "failed to save foreign message mapping for bot",
				slog.String("telegram_message_id", msgIDStr),
				slog.Int64("db_message_id", botMessage.ID),
				slog.String("error", err.Error()))
		} else {
			s.logger.DebugContext(ctx, "saved foreign message mapping",
				slog.String("telegram_message_id", msgIDStr),
				slog.Int64("db_message_id", botMessage.ID))
		}
	}

	// Deduct base model tokens after successful completion
	baseTransaction := &domain.Transaction{
		UserID:          user.ID,
		TokenType:       currentModel.TokenType,
		Amount:          -currentModel.Cost, // Negative for debit
		TransactionType: domain.TransactionTypeMessageCost,
		ModelUsed:       &currentModel.ID,
		Description:     nil,
		CreatedAt:       time.Now(),
	}

	_, err = s.storage.CreateTransaction(ctx, baseTransaction)
	if err != nil {
		// Log error but don't fail the entire operation
		s.logger.ErrorContext(ctx, "failed to create base token deduction transaction",
			slog.String("error", err.Error()),
			slog.String("model", currentModel.ID),
			slog.Int("cost", int(currentModel.Cost)))
	}

	// Deduct search tokens if web search was used
	if webSearchEnabled && currentModel.SearchCost != nil && currentModel.SearchTokenType != nil {
		searchTransaction := &domain.Transaction{
			UserID:          user.ID,
			TokenType:       *currentModel.SearchTokenType,
			Amount:          -*currentModel.SearchCost, // Negative for debit
			TransactionType: domain.TransactionTypeMessageCost,
			ModelUsed:       &currentModel.ID,
			Description:     pointer.To("Web search cost"),
			CreatedAt:       time.Now(),
		}

		_, err = s.storage.CreateTransaction(ctx, searchTransaction)
		if err != nil {
			// Log error but don't fail the entire operation
			s.logger.ErrorContext(ctx, "failed to create search token deduction transaction",
				slog.String("error", err.Error()),
				slog.String("model", currentModel.ID),
				slog.Int("search_cost", int(*currentModel.SearchCost)))
		}
	}

	return nil
}

// handleReasoningUpdate manages reasoning message creation and updates.
func (s *UpdateService) handleReasoningUpdate(
	ctx context.Context,
	user *domain.User,
	reasoningBuilder *strings.Builder,
	reasoningMessageIDs []string,
	previousReasoningContent string,
	hasReasoningMessage bool,
) ([]string, string, bool) {
	if reasoningBuilder.Len() == 0 {
		return reasoningMessageIDs, previousReasoningContent, hasReasoningMessage
	}

	currentReasoningContent := "```reasoning\n" + reasoningBuilder.String() + "\n```"

	// Send initial reasoning message if not yet sent
	if !hasReasoningMessage {
		reasoningMsgID, reasoningErr := s.sender.SendMessage(ctx, user.ExternalID, currentReasoningContent)
		if reasoningErr != nil {
			s.logger.WarnContext(ctx, "failed to send reasoning message",
				slog.String("error", reasoningErr.Error()))
			return reasoningMessageIDs, previousReasoningContent, hasReasoningMessage
		}
		return []string{reasoningMsgID}, currentReasoningContent, true
	}

	// Update existing reasoning message if content changed
	if currentReasoningContent != previousReasoningContent {
		updatedReasoningIDs, updateErr := s.sender.UpdateMessages(
			ctx,
			user.ExternalID,
			reasoningMessageIDs,
			previousReasoningContent,
			currentReasoningContent,
		)
		if updateErr != nil {
			// Don't fail on rate limit errors - we'll retry next time
			if strings.Contains(updateErr.Error(), "Too Many Requests") ||
				strings.Contains(updateErr.Error(), "too many requests") {
				s.logger.WarnContext(ctx, "rate limited while updating reasoning messages, will retry",
					slog.String("error", updateErr.Error()))
				return reasoningMessageIDs, previousReasoningContent, hasReasoningMessage
			}
			// Log other errors but continue
			s.logger.WarnContext(ctx, "failed to update reasoning messages",
				slog.String("error", updateErr.Error()))
			return reasoningMessageIDs, previousReasoningContent, hasReasoningMessage
		}
		return updatedReasoningIDs, currentReasoningContent, hasReasoningMessage
	}

	return reasoningMessageIDs, previousReasoningContent, hasReasoningMessage
}

// handleMainContentUpdate manages main content message creation and updates.
func (s *UpdateService) handleMainContentUpdate(
	ctx context.Context,
	user *domain.User,
	responseBuilder *strings.Builder,
	messageIDs []string,
	previousContent string,
	hasMainMessage bool,
	replyToMessageID *int64,
) ([]string, string, bool, error) {
	currentContent := responseBuilder.String()
	if currentContent == previousContent {
		return messageIDs, previousContent, hasMainMessage, nil
	}

	// Send initial main message if not yet sent
	if !hasMainMessage {
		initialContent := domain.MessageContent{
			Text:             currentContent,
			ReplyToMessageID: replyToMessageID,
		}
		messageID, err := s.sender.SendMessageWithContent(ctx, user.ExternalID, initialContent)
		if err != nil {
			return messageIDs, previousContent, hasMainMessage, fmt.Errorf("can't send main message: %w", err)
		}
		return []string{messageID}, currentContent, true, nil
	}

	// Send typing indicator to keep the conversation active
	err := s.sender.SendTyping(ctx, user.ExternalID)
	if err != nil {
		s.logger.WarnContext(ctx, "failed to send typing indicator", slog.String("error", err.Error()))
	}

	// Update all tracked messages with optimization
	updatedMessageIDs, updateErr := s.sender.UpdateMessages(
		ctx,
		user.ExternalID,
		messageIDs,
		previousContent,
		currentContent,
	)
	if updateErr != nil {
		// Check if it's a rate limit error
		if strings.Contains(updateErr.Error(), "Too Many Requests") ||
			strings.Contains(updateErr.Error(), "too many requests") {
			s.logger.WarnContext(ctx, "rate limited while updating main content, will retry",
				slog.String("error", updateErr.Error()))
			return messageIDs, previousContent, hasMainMessage, nil
		}
		// Log error but don't send another error message to user to avoid spam
		s.logger.ErrorContext(ctx, "failed to update messages during streaming",
			slog.String("error", updateErr.Error()),
			slog.String("user_id", user.ExternalID))
		return messageIDs, previousContent, hasMainMessage, fmt.Errorf("can't update messages: %w", updateErr)
	}

	return updatedMessageIDs, currentContent, hasMainMessage, nil
}

// sendFinalUpdates handles final message updates at the end of streaming.
func (s *UpdateService) sendFinalUpdates(
	ctx context.Context,
	user *domain.User,
	responseBuilder *strings.Builder,
	reasoningBuilder *strings.Builder,
	messageIDs []string,
	reasoningMessageIDs []string,
	previousContent string,
	previousReasoningContent string,
	hasReasoningMessage bool,
	hasMainMessage bool,
	replyToMessageID *int64,
) ([]string, error) {
	// Final reasoning update
	s.handleFinalReasoningUpdate(ctx, user, reasoningBuilder, reasoningMessageIDs,
		previousReasoningContent, hasReasoningMessage)

	// Final content update - create message if needed
	return s.handleFinalContentUpdate(ctx, user, responseBuilder, messageIDs,
		previousContent, hasMainMessage, replyToMessageID)
}

// handlePeriodicUpdates manages periodic message updates during streaming.
func (s *UpdateService) handlePeriodicUpdates(
	ctx context.Context,
	user *domain.User,
	reasoningBuilder *strings.Builder,
	responseBuilder *strings.Builder,
	reasoningMessageIDs *[]string,
	messageIDs *[]string,
	previousReasoningContent *string,
	previousContent *string,
	hasReasoningMessage *bool,
	hasMainMessage *bool,
	replyToMessageID *int64,
) (bool, error) {
	// Handle reasoning message updates (create reasoning message first if we have reasoning)
	if reasoningBuilder.Len() > 0 {
		*reasoningMessageIDs, *previousReasoningContent, *hasReasoningMessage = s.handleReasoningUpdate(
			ctx,
			user,
			reasoningBuilder,
			*reasoningMessageIDs,
			*previousReasoningContent,
			*hasReasoningMessage,
		)
	}

	// Handle main content updates (create main message after reasoning)
	if responseBuilder.Len() > 0 {
		updatedMessageIDs, updatedContent, updatedHasMessage, err := s.handleMainContentUpdate(
			ctx,
			user,
			responseBuilder,
			*messageIDs,
			*previousContent,
			*hasMainMessage,
			replyToMessageID,
		)
		if err != nil {
			// Check if it's a rate limit error
			if strings.Contains(err.Error(), "Too Many Requests") ||
				strings.Contains(err.Error(), "too many requests") {
				s.logger.WarnContext(ctx, "rate limited by Telegram, will retry in next iteration",
					slog.String("error", err.Error()))
				// Skip this update and continue - don't update lastUpdate so we retry sooner
				return true, nil // shouldContinue = true
			}
			return false, err
		}
		*messageIDs = updatedMessageIDs
		*previousContent = updatedContent
		*hasMainMessage = updatedHasMessage
	}

	return false, nil // shouldContinue = false
}

// handleFinalReasoningUpdate handles the final reasoning message update.
func (s *UpdateService) handleFinalReasoningUpdate(
	ctx context.Context,
	user *domain.User,
	reasoningBuilder *strings.Builder,
	reasoningMessageIDs []string,
	previousReasoningContent string,
	hasReasoningMessage bool,
) {
	if hasReasoningMessage && reasoningBuilder.Len() > 0 {
		finalReasoningContent := "```reasoning\n" + reasoningBuilder.String() + "\n```"
		if finalReasoningContent != previousReasoningContent {
			s.updateFinalReasoningContent(ctx, user, reasoningMessageIDs,
				previousReasoningContent, finalReasoningContent)
		}
	}
}

// updateFinalReasoningContent updates the final reasoning content.
func (s *UpdateService) updateFinalReasoningContent(
	ctx context.Context,
	user *domain.User,
	reasoningMessageIDs []string,
	previousReasoningContent string,
	finalReasoningContent string,
) {
	_, updateErr := s.sender.UpdateMessages(
		ctx,
		user.ExternalID,
		reasoningMessageIDs,
		previousReasoningContent,
		finalReasoningContent,
	)
	if updateErr != nil {
		// Don't fail on rate limit errors - we'll retry next time
		if strings.Contains(updateErr.Error(), "Too Many Requests") ||
			strings.Contains(updateErr.Error(), "too many requests") {
			s.logger.WarnContext(ctx, "rate limited while updating final reasoning messages",
				slog.String("error", updateErr.Error()))
		} else {
			s.logger.WarnContext(ctx, "failed to update final reasoning messages",
				slog.String("error", updateErr.Error()))
		}
	}
}

// handleFinalContentUpdate handles the final content message update or creation.
func (s *UpdateService) handleFinalContentUpdate(
	ctx context.Context,
	user *domain.User,
	responseBuilder *strings.Builder,
	messageIDs []string,
	previousContent string,
	hasMainMessage bool,
	replyToMessageID *int64,
) ([]string, error) {
	finalContent := responseBuilder.String()

	if finalContent != "" && !hasMainMessage {
		return s.createFinalMainMessage(ctx, user, finalContent, messageIDs, replyToMessageID)
	} else if finalContent != previousContent && hasMainMessage {
		return s.updateFinalMainMessage(ctx, user, messageIDs, previousContent, finalContent)
	}

	return messageIDs, nil
}

// createFinalMainMessage creates a main message if we never created one.
func (s *UpdateService) createFinalMainMessage(
	ctx context.Context,
	user *domain.User,
	finalContent string,
	messageIDs []string,
	replyToMessageID *int64,
) ([]string, error) {
	initialContent := domain.MessageContent{
		Text:             finalContent,
		ReplyToMessageID: replyToMessageID,
	}
	messageID, err := s.sender.SendMessageWithContent(ctx, user.ExternalID, initialContent)
	if err != nil {
		// Check if it's a rate limit error
		if strings.Contains(err.Error(), "Too Many Requests") ||
			strings.Contains(err.Error(), "too many requests") {
			s.logger.WarnContext(ctx, "rate limited while sending final main message",
				slog.String("error", err.Error()))
			return messageIDs, nil
		}
		return messageIDs, fmt.Errorf("can't send final main message: %w", err)
	}
	return []string{messageID}, nil
}

// updateFinalMainMessage updates the existing main message.
func (s *UpdateService) updateFinalMainMessage(
	ctx context.Context,
	user *domain.User,
	messageIDs []string,
	previousContent string,
	finalContent string,
) ([]string, error) {
	updatedMessageIDs, finalUpdateErr := s.sender.UpdateMessages(
		ctx,
		user.ExternalID,
		messageIDs,
		previousContent,
		finalContent,
	)
	if finalUpdateErr != nil {
		// Check if it's a rate limit error
		if strings.Contains(finalUpdateErr.Error(), "Too Many Requests") ||
			strings.Contains(finalUpdateErr.Error(), "too many requests") {
			s.logger.WarnContext(ctx, "rate limited while updating final main content",
				slog.String("error", finalUpdateErr.Error()))
			return messageIDs, nil
		}
		s.logger.ErrorContext(ctx, "failed to update final messages",
			slog.String("error", finalUpdateErr.Error()),
			slog.String("user_id", user.ExternalID))
		return messageIDs, fmt.Errorf("can't update final messages: %w", finalUpdateErr)
	}
	return updatedMessageIDs, nil
}

func (s *UpdateService) generateAndUpdateConversationName(
	ctx context.Context,
	conversationID int64,
	firstMessage string,
) {
	const systemPrompt = `Generate a short, descriptive name for a chat conversation based on the user's first message. 
The name should be 2-5 words max and capture the main topic or intent.
Return ONLY the conversation name, nothing else.
Examples:
- User: "How do I cook pasta?" -> "Cooking Pasta Guide"
- User: "Tell me about quantum physics" -> "Quantum Physics Discussion"
- User: "I need help with Python code" -> "Python Code Help"`

	messages := []*domain.Message{
		{
			MessageType: domain.MessageType{Text: firstMessage},
			SentBy:      domain.MessageSenderUser,
		},
	}

	// Use Gemini model for generating conversation name
	tokenStream, err := s.completion.CompleteStream(
		ctx,
		"google/gemini-2.5-flash",
		systemPrompt,
		messages,
		"",    // no image for conversation name generation
		false, // no web search for conversation name generation
	)
	if err != nil {
		s.logger.WarnContext(ctx, "failed to generate conversation name with LLM, using fallback",
			slog.String("error", err.Error()))
		// Fallback to date-based naming
		fallbackName := fmt.Sprintf("Conversation %s", time.Now().Format("Jan 2 15:04"))
		if updateErr := s.storage.UpdateConversationName(ctx, conversationID, fallbackName); updateErr != nil {
			s.logger.ErrorContext(ctx, "failed to update conversation name with fallback",
				slog.String("error", updateErr.Error()))
		}
		return
	}

	var nameBuilder strings.Builder
	for token := range tokenStream {
		if token.Error != nil {
			s.logger.WarnContext(ctx, "error in conversation name generation stream",
				slog.String("error", token.Error.Error()))
			// Fallback to date-based naming
			fallbackName := fmt.Sprintf("Conversation %s", time.Now().Format("Jan 2 15:04"))
			if updateErr := s.storage.UpdateConversationName(ctx, conversationID, fallbackName); updateErr != nil {
				s.logger.ErrorContext(ctx, "failed to update conversation name with fallback",
					slog.String("error", updateErr.Error()))
			}
			return
		}
		nameBuilder.WriteString(token.Content)
	}

	generatedName := strings.TrimSpace(nameBuilder.String())
	// Ensure the name is not too long
	if len(generatedName) > maxConversationNameLength {
		generatedName = generatedName[:maxConversationNameLength-3] + "..."
	}

	// If generated name is empty or too short, use fallback
	if len(generatedName) < minConversationNameLength {
		generatedName = fmt.Sprintf("Conversation %s", time.Now().Format("Jan 2 15:04"))
	}

	// Update the conversation name
	if updateErr := s.storage.UpdateConversationName(ctx, conversationID, generatedName); updateErr != nil {
		s.logger.ErrorContext(ctx, "failed to update conversation name",
			slog.String("error", updateErr.Error()))
	} else {
		s.logger.InfoContext(ctx, "conversation name generated successfully",
			slog.String("name", generatedName))
	}
}

func (s *UpdateService) processQueuedMessages(ctx context.Context, user *domain.User) {
	for {
		// Get next queued update
		queuedItem, err := s.queue.DequeueWithMetadata(ctx, user.ExternalID)
		if err != nil {
			if errors.Is(err, queue.ErrEmptyQueue) {
				// No more queued messages
				return
			}
			s.logger.ErrorContext(ctx, "failed to dequeue update",
				slog.String("error", err.Error()))
			return
		}

		// Delete the queue notification message if exists
		if queuedItem.QueueNotificationID != "" {
			if deleteErr := s.sender.DeleteMessage(ctx, user.ExternalID, queuedItem.QueueNotificationID); deleteErr != nil {
				s.logger.WarnContext(ctx, "failed to delete queue notification",
					slog.String("message_id", queuedItem.QueueNotificationID),
					slog.String("error", deleteErr.Error()))
			}
		}

		// Process the queued update
		s.logger.InfoContext(ctx, "processing queued message",
			slog.String("user_id", user.ExternalID),
			slog.String("message", queuedItem.Update.MessageText))

		// Re-fetch user to ensure we have latest state
		freshUser, err := s.storage.GetUserByExternalUserID(ctx, user.ExternalID)
		if err != nil {
			s.logger.ErrorContext(ctx, "failed to get user for queued message",
				slog.String("error", err.Error()))
			continue
		}

		// Process the update with reply to the original message
		// We need to handle conversation messages specially to add reply_to_message_id
		if freshUser.CurrentStep == domain.UserStateConversation && queuedItem.Update.MessageText != "" {
			// Convert external message ID to int64 for reply
			var replyToMessageID *int64
			if queuedItem.Update.ExternalMessageID > 0 {
				msgID := int64(queuedItem.Update.ExternalMessageID)
				replyToMessageID = &msgID
			}

			// Process with reply to original message
			if processErr := s.handleConversationMessageWithReply(ctx, freshUser, queuedItem.Update, replyToMessageID); processErr != nil {
				s.logger.ErrorContext(ctx, "failed to process queued conversation update",
					slog.String("error", processErr.Error()))
			}
		} else {
			// Process normally for non-conversation states
			if processErr := s.processUpdate(ctx, freshUser, queuedItem.Update); processErr != nil {
				s.logger.ErrorContext(ctx, "failed to process queued update",
					slog.String("error", processErr.Error()))
			}
		}

		// Small delay between processing queued messages
		time.Sleep(queueProcessingDelay)
	}
}

type imageUploadResult struct {
	imageURL    string
	objectName  string
	contentType string
	size        int64
}

func (s *UpdateService) handleImageUpload(
	ctx context.Context,
	imageData []byte,
	imageMimeType string,
) *imageUploadResult {
	if len(imageData) == 0 || imageMimeType == "" || s.fileStorage == nil {
		return nil
	}

	// Upload image to S3 and get pre-signed URL
	objectName, uploadErr := s.fileStorage.Upload(ctx, imageData, imageMimeType)
	if uploadErr != nil {
		s.logger.ErrorContext(ctx, "failed to upload image", slog.String("error", uploadErr.Error()))
		return nil
	}

	// Generate pre-signed URL valid for 1 hour
	imageURL, urlErr := s.fileStorage.GetPreSignedURL(ctx, objectName, time.Hour)
	if urlErr != nil {
		s.logger.ErrorContext(ctx, "failed to generate pre-signed URL", slog.String("error", urlErr.Error()))
		return nil
	}

	s.logger.InfoContext(ctx, "Image uploaded successfully",
		slog.String("object_name", objectName),
		slog.String("image_url", imageURL))

	return &imageUploadResult{
		imageURL:    imageURL,
		objectName:  objectName,
		contentType: imageMimeType,
		size:        int64(len(imageData)),
	}
}

type pdfUploadResult struct {
	pdfURL      string
	objectName  string
	contentType string
	size        int64
}

func (s *UpdateService) handlePDFUpload(
	ctx context.Context,
	pdfData []byte,
	pdfMimeType string,
	pdfFileName string,
) *pdfUploadResult {
	if len(pdfData) == 0 || pdfMimeType == "" || s.fileStorage == nil {
		return nil
	}

	// Upload PDF to S3 and get pre-signed URL
	objectName, uploadErr := s.fileStorage.Upload(ctx, pdfData, pdfMimeType)
	if uploadErr != nil {
		s.logger.ErrorContext(ctx, "failed to upload PDF", slog.String("error", uploadErr.Error()))
		return nil
	}

	// Generate pre-signed URL valid for 1 hour
	pdfURL, urlErr := s.fileStorage.GetPreSignedURL(ctx, objectName, time.Hour)
	if urlErr != nil {
		s.logger.ErrorContext(ctx, "failed to generate pre-signed URL for PDF", slog.String("error", urlErr.Error()))
		return nil
	}

	s.logger.InfoContext(ctx, "PDF uploaded successfully",
		slog.String("object_name", objectName),
		slog.String("pdf_url", pdfURL),
		slog.String("file_name", pdfFileName))

	return &pdfUploadResult{
		pdfURL:      pdfURL,
		objectName:  objectName,
		contentType: pdfMimeType,
		size:        int64(len(pdfData)),
	}
}

// sendPeriodicTyping sends typing indicators every 4.5 seconds during generation.
func (s *UpdateService) sendPeriodicTyping(ctx context.Context, userExternalID string, done <-chan struct{}) {
	ticker := time.NewTicker(typingIndicatorInterval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.sender.SendTyping(ctx, userExternalID); err != nil {
				s.logger.WarnContext(ctx, "failed to send periodic typing indicator",
					slog.String("error", err.Error()),
					slog.String("user_id", userExternalID))
			}
		}
	}
}
