package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/queue"
)

const (
	processingLockTimeout     = 5 * time.Minute
	maxConversationNameLength = 50
	minConversationNameLength = 2
	queueProcessingDelay      = 100 * time.Millisecond
)

// Conversation handlers.
func (s *UpdateService) handleConversationMessage(ctx context.Context, user *domain.User, update domain.Update) error {
	// Set processing lock with 5 minute timeout
	if err := s.queue.SetProcessing(ctx, user.ExternalID, processingLockTimeout); err != nil {
		if errors.Is(err, queue.ErrAlreadyProcessing) {
			// This shouldn't happen as we check before, but handle it gracefully
			return s.queue.EnqueueWithNotification(ctx, user.ExternalID, update, "")
		}
		s.logger.WarnContext(ctx, "failed to set processing lock",
			slog.String("error", err.Error()))
		// Continue without lock on error
	}

	// Ensure we clear the lock and process queued messages when done
	defer func() {
		if err := s.queue.ClearProcessing(ctx, user.ExternalID); err != nil {
			s.logger.ErrorContext(ctx, "failed to clear processing lock",
				slog.String("error", err.Error()))
		}

		// Process any queued messages
		go s.processQueuedMessages(context.Background(), user)
	}()
	// Check if this is the first message in a new conversation
	var isFirstMessage bool
	if user.CurrentConversationID != nil {
		messages, err := s.storage.GetMessagesByConversationID(ctx, *user.CurrentConversationID)
		if err != nil {
			return fmt.Errorf("can't check existing messages: %w", err)
		}
		isFirstMessage = len(messages) == 0
	}

	// Save incoming message from user
	userMessage, err := s.storage.CreateMessage(ctx, &domain.Message{
		UserID: user.ID,
		MessageType: domain.MessageType{
			Text: update.MessageText,
		},
		SentBy:         domain.MessageSenderUser,
		ConversationID: user.CurrentConversationID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		return fmt.Errorf("can't save user message: %w", err)
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

	// Send typing indicator
	err = s.sender.SendTyping(ctx, user.ExternalID)
	if err != nil {
		s.logger.WarnContext(ctx, "failed to send typing indicator", slog.String("error", err.Error()))
	}

	systemPrompt := "You are a helpful assistant."
	tokenStream, err := s.completion.CompleteStream(ctx, user.SelectedModel, systemPrompt, messages)
	if err != nil {
		return fmt.Errorf("can't get completion: %w", err)
	}

	// Send initial empty message to get message ID
	// If this is from a queued message, reply to the original
	var replyToMessageID *int64
	if update.ExternalMessageID > 0 {
		msgID := int64(update.ExternalMessageID)
		replyToMessageID = &msgID
	}

	initialContent := domain.MessageContent{
		Text:             "\\.\\.\\.",
		ReplyToMessageID: replyToMessageID,
	}
	messageID, err := s.sender.SendMessageWithContent(ctx, user.ExternalID, initialContent)
	if err != nil {
		return fmt.Errorf("can't send initial message: %w", err)
	}

	// Track all message IDs (starts with just one)
	messageIDs := []string{messageID}

	// Stream response tokens and update message
	var responseBuilder strings.Builder
	var previousContent string
	lastUpdate := time.Now()

	for token := range tokenStream {
		if token.Error != nil {
			return fmt.Errorf("completion stream error: %w", token.Error)
		}
		responseBuilder.WriteString(token.Content)

		// Update message at most once per second
		if time.Since(lastUpdate) >= time.Second/2 {
			currentContent := responseBuilder.String()

			// Send typing indicator to keep the conversation active
			err = s.sender.SendTyping(ctx, user.ExternalID)
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
				return fmt.Errorf("can't update messages: %w", updateErr)
			}
			// Update our tracked message IDs
			messageIDs = updatedMessageIDs
			previousContent = currentContent
			lastUpdate = time.Now()
		}
	}

	// Send final update if needed
	if time.Since(lastUpdate) > 0 {
		finalContent := responseBuilder.String()
		updatedMessageIDs, finalUpdateErr := s.sender.UpdateMessages(
			ctx,
			user.ExternalID,
			messageIDs,
			previousContent,
			finalContent,
		)
		if finalUpdateErr != nil {
			return fmt.Errorf("can't update final messages: %w", finalUpdateErr)
		}
		// Update our tracked message IDs for final state
		messageIDs = updatedMessageIDs
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

	return nil
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
		"google/gemini-2.5-flash-preview-05-20",
		systemPrompt,
		messages,
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

		// Process the update - it will reply to the original message automatically
		// since the update contains the original message ID
		if processErr := s.processUpdate(ctx, freshUser, queuedItem.Update); processErr != nil {
			s.logger.ErrorContext(ctx, "failed to process queued update",
				slog.String("error", processErr.Error()))
		}

		// Small delay between processing queued messages
		time.Sleep(queueProcessingDelay)
	}
}
