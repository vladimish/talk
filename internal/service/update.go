package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sort"
	"time"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/completion"
	"github.com/vladimish/talk/internal/port/filestorage"
	"github.com/vladimish/talk/internal/port/queue"
	"github.com/vladimish/talk/internal/port/sender"
	"github.com/vladimish/talk/internal/port/storage"
	"github.com/vladimish/talk/pkg/i18n"
)

const (
	initialTokenGrant             = 20
	messageConcatenationWindow    = time.Second      // Time window to wait for additional messages
	messageConcatenationTTLFactor = 2                // Factor for TTL calculation
	generationLockDuration        = 10 * time.Minute // Duration for generation lock
)

type UpdateService struct {
	logger      *slog.Logger
	storage     storage.Storage
	sender      sender.Sender
	completion  completion.Completion
	queue       queue.Queue
	fileStorage filestorage.FileStorage
}

func NewUpdateService(
	logger *slog.Logger,
	storage storage.Storage,
	sender sender.Sender,
	completion completion.Completion,
	queue queue.Queue,
	fileStorage filestorage.FileStorage,
) *UpdateService {
	return &UpdateService{
		logger:      logger,
		storage:     storage,
		sender:      sender,
		completion:  completion,
		queue:       queue,
		fileStorage: fileStorage,
	}
}

func (s *UpdateService) HandleUpdate(ctx context.Context, update domain.Update) (err error) {
	// Add panic recovery to prevent crashes during user request handling
	defer func() {
		if r := recover(); r != nil {
			stackTrace := debug.Stack()
			s.logger.ErrorContext(ctx, "Panic occurred while handling user update",
				"panic", r,
				"stack_trace", string(stackTrace),
				"user_id", update.ExternalUserID,
				"message_text", update.MessageText)

			// Convert panic to error
			err = fmt.Errorf("panic occurred while handling update: %v", r)
		}
	}()

	// Set received timestamp for message ordering
	if update.ReceivedAt.IsZero() {
		update.ReceivedAt = time.Now()
	}

	// Check for /start command first - always redirect to menu regardless of current state
	if update.MessageText == "/start" {
		user, userErr := s.getOrCreateUser(ctx, update)
		if userErr != nil {
			return userErr
		}
		return s.transitionToMenu(ctx, user)
	}

	// Check if this user is in conversation state and might need queueing
	user, err := s.getOrCreateUser(ctx, update)
	if err != nil {
		return err
	}

	// Only queue messages in conversation state
	if s.shouldQueueMessage(user, update) {
		queued, queueErr := s.handleMessageQueueing(ctx, user, update)
		if queueErr != nil {
			return queueErr
		}
		if queued {
			return nil
		}
	}

	// Process the update normally
	return s.processUpdate(ctx, user, update)
}

func (s *UpdateService) processUpdate(ctx context.Context, user *domain.User, update domain.Update) error {
	// Handle based on current user state
	currentState := user.CurrentStep

	var err error
	switch currentState {
	case domain.UserStateMenu:
		err = s.HandleMenuState(ctx, user, update)
	case domain.UserStateConversation:
		err = s.HandleConversationState(ctx, user, update)
	case domain.UserStateModelSelect:
		err = s.HandleModelSelectState(ctx, user, update)
	case domain.UserStateConversationList:
		err = s.HandleConversationListState(ctx, user, update)
	case domain.UserStateSettings:
		err = s.HandleSettingsState(ctx, user, update)
	case domain.UserStateLanguageSelect:
		err = s.HandleLanguageSelectState(ctx, user, update)
	case domain.UserStateProfile:
		err = s.HandleProfileState(ctx, user, update)
	default:
		// Default to menu state for unknown states
		err = s.HandleMenuState(ctx, user, update)
	}

	// If an error occurred during handler execution, send user-friendly message
	if err != nil {
		// Log the error for debugging
		s.logger.ErrorContext(ctx, "error processing update",
			slog.String("error", err.Error()),
			slog.String("user_id", user.ExternalID),
			slog.String("state", currentState))

		// Send localized error message to user
		errorMsg := i18n.GetString(user.Language, i18n.ErrorResponseGeneration)
		_, sendErr := s.sender.SendMessage(ctx, user.ExternalID, errorMsg)
		if sendErr != nil {
			s.logger.ErrorContext(ctx, "failed to send error message to user",
				slog.String("send_error", sendErr.Error()),
				slog.String("original_error", err.Error()))
		}
	}

	return err
}

func (s *UpdateService) getOrCreateUser(ctx context.Context, update domain.Update) (*domain.User, error) {
	user, err := s.storage.GetUserByExternalUserID(ctx, update.ExternalUserID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return s.createNewUserWithTokens(ctx, update)
		}
		return nil, fmt.Errorf("can't get user: %w", err)
	}
	return user, nil
}

// handleConversationMessageWithConcatenation handles message concatenation for split messages.
func (s *UpdateService) handleConversationMessageWithConcatenation(
	ctx context.Context,
	user *domain.User,
	update domain.Update,
) error {
	// Check if there's an ongoing generation that we need to cancel
	isGenerating, err := s.queue.IsGenerating(ctx, user.ExternalID)
	if err != nil {
		s.logger.WarnContext(ctx, "failed to check generation status",
			slog.String("error", err.Error()))
		// Continue anyway
	}

	// If there's an ongoing generation, we'll cancel it by clearing the lock and adding to pending
	if isGenerating {
		if clearErr := s.queue.ClearGenerationLock(ctx, user.ExternalID); clearErr != nil {
			s.logger.WarnContext(ctx, "failed to clear generation lock",
				slog.String("error", clearErr.Error()))
		}
		s.logger.InfoContext(ctx, "cancelled ongoing generation for message concatenation",
			slog.String("user_id", user.ExternalID))
	}

	// Check if there are already pending messages
	pending, err := s.queue.GetPendingMessages(ctx, user.ExternalID)
	if err != nil && !errors.Is(err, queue.ErrEmptyQueue) {
		s.logger.WarnContext(ctx, "failed to get pending messages",
			slog.String("error", err.Error()))
		// Continue with just this message
		pending = nil
	} else if errors.Is(err, queue.ErrEmptyQueue) {
		pending = nil
	}

	now := time.Now()
	var messages []domain.Update

	if pending != nil {
		// Add to existing pending messages
		messages = make([]domain.Update, len(pending.Messages), len(pending.Messages)+1)
		copy(messages, pending.Messages)
		messages = append(messages, update)
		s.logger.InfoContext(ctx, "added message to concatenation batch",
			slog.String("user_id", user.ExternalID),
			slog.Int("total_messages", len(messages)),
			slog.Int("new_message_id", update.ExternalMessageID))
	} else {
		// Start new batch
		messages = []domain.Update{update}
		s.logger.InfoContext(ctx, "started new message concatenation batch",
			slog.String("user_id", user.ExternalID),
			slog.Int("message_id", update.ExternalMessageID))
	}

	// Store updated pending messages with timeout
	newPending := &queue.PendingMessages{
		Messages:  messages,
		Timestamp: now,
	}

	if setErr := s.queue.SetPendingMessages(ctx, user.ExternalID, newPending, messageConcatenationWindow*messageConcatenationTTLFactor); setErr != nil {
		s.logger.ErrorContext(ctx, "failed to set pending messages",
			slog.String("error", setErr.Error()))
		// Fall back to immediate processing
		return s.handleConversationMessage(ctx, user, update)
	}

	// Schedule processing after concatenation window
	go func() {
		time.Sleep(messageConcatenationWindow)
		s.processPendingMessages(context.Background(), user)
	}()

	return nil
}

// processPendingMessages processes all pending messages for a user.
func (s *UpdateService) processPendingMessages(ctx context.Context, user *domain.User) {
	// Get pending messages
	pending, err := s.queue.GetPendingMessages(ctx, user.ExternalID)
	if err != nil && !errors.Is(err, queue.ErrEmptyQueue) {
		s.logger.ErrorContext(ctx, "failed to get pending messages for processing",
			slog.String("error", err.Error()))
		return
	} else if errors.Is(err, queue.ErrEmptyQueue) {
		return // No messages to process
	}

	if pending == nil || len(pending.Messages) == 0 {
		return // No messages to process
	}

	// Clear pending messages
	if clearErr := s.queue.ClearPendingMessages(ctx, user.ExternalID); clearErr != nil {
		s.logger.WarnContext(ctx, "failed to clear pending messages",
			slog.String("error", clearErr.Error()))
	}

	s.logger.InfoContext(ctx, "processing pending concatenated messages",
		slog.String("user_id", user.ExternalID),
		slog.Int("message_count", len(pending.Messages)))

	// Combine all messages
	combinedUpdate := s.combineMessages(pending.Messages)

	// Set generation lock to prevent new concatenations during processing
	if lockErr := s.queue.SetGenerationLock(ctx, user.ExternalID, generationLockDuration); lockErr != nil {
		s.logger.WarnContext(ctx, "failed to set generation lock",
			slog.String("error", lockErr.Error()))
	}

	// Process the combined message
	// For concatenated messages, don't reply - this is just regular conversation flow
	err = s.handleConversationMessage(ctx, user, combinedUpdate)

	// Clear generation lock
	if clearErr := s.queue.ClearGenerationLock(ctx, user.ExternalID); clearErr != nil {
		s.logger.WarnContext(ctx, "failed to clear generation lock after processing",
			slog.String("error", clearErr.Error()))
	}

	if err != nil {
		s.logger.ErrorContext(ctx, "failed to process concatenated messages",
			slog.String("error", err.Error()),
			slog.String("user_id", user.ExternalID))

		// Send error message to user
		errorMsg := i18n.GetString(user.Language, i18n.ErrorResponseGeneration)
		_, sendErr := s.sender.SendMessage(ctx, user.ExternalID, errorMsg)
		if sendErr != nil {
			s.logger.ErrorContext(ctx, "failed to send error message to user",
				slog.String("send_error", sendErr.Error()),
				slog.String("original_error", err.Error()))
		}
	}
}

// combineMessages combines multiple domain.Update messages into a single one in correct order.
func (s *UpdateService) combineMessages(messages []domain.Update) domain.Update {
	if len(messages) == 0 {
		return domain.Update{}
	}

	if len(messages) == 1 {
		return messages[0]
	}

	// Sort messages by external message ID first, then by received timestamp
	// This ensures messages are combined in the correct order
	sortedMessages := make([]domain.Update, len(messages))
	copy(sortedMessages, messages)

	sort.Slice(sortedMessages, func(i, j int) bool {
		// Primary sort: by external message ID (Telegram message order)
		if sortedMessages[i].ExternalMessageID != sortedMessages[j].ExternalMessageID {
			return sortedMessages[i].ExternalMessageID < sortedMessages[j].ExternalMessageID
		}
		// Secondary sort: by received timestamp for messages with same ID (should be rare)
		return sortedMessages[i].ReceivedAt.Before(sortedMessages[j].ReceivedAt)
	})

	s.logger.DebugContext(context.Background(), "sorted messages for concatenation",
		slog.Int("message_count", len(sortedMessages)),
		slog.Any("message_ids", extractMessageIDs(sortedMessages)))

	// Combine all message texts in sorted order
	combinedText := ""
	var combinedImageData []byte
	var combinedImageMimeType string
	var combinedPDFData []byte
	var combinedPDFMimeType string
	var combinedPDFFileName string
	var lastExternalMessageID int64

	for i, msg := range sortedMessages {
		if i > 0 {
			combinedText += " " // Add space between messages
		}
		combinedText += msg.MessageText

		// Use attachments from the last message that has them (in sorted order)
		if len(msg.ImageData) > 0 {
			combinedImageData = msg.ImageData
			combinedImageMimeType = msg.ImageMimeType
		}
		if len(msg.PDFData) > 0 {
			combinedPDFData = msg.PDFData
			combinedPDFMimeType = msg.PDFMimeType
			combinedPDFFileName = msg.PDFFileName
		}
		if msg.ExternalMessageID > 0 {
			lastExternalMessageID = int64(msg.ExternalMessageID)
		}
	}

	// Create combined update using the first message as base
	baseMsg := sortedMessages[0]
	return domain.Update{
		ExternalUserID:    baseMsg.ExternalUserID,
		UserLanguage:      baseMsg.UserLanguage,
		MessageText:       combinedText,
		ImageData:         combinedImageData,
		ImageMimeType:     combinedImageMimeType,
		PDFData:           combinedPDFData,
		PDFMimeType:       combinedPDFMimeType,
		PDFFileName:       combinedPDFFileName,
		ExternalMessageID: int(lastExternalMessageID),
		ReceivedAt:        baseMsg.ReceivedAt, // Use the first message's timestamp
	}
}

// extractMessageIDs extracts message IDs for logging purposes.
func extractMessageIDs(messages []domain.Update) []int {
	ids := make([]int, len(messages))
	for i, msg := range messages {
		ids[i] = msg.ExternalMessageID
	}
	return ids
}

func (s *UpdateService) createNewUserWithTokens(ctx context.Context, update domain.Update) (*domain.User, error) {
	language := update.UserLanguage
	now := time.Now()
	newUser := &domain.User{
		ExternalID:             update.ExternalUserID,
		Language:               language,
		CurrentStep:            domain.UserStateMenu,
		SelectedModel:          "google/gemini-2.5-flash-preview-05-20",
		ConversationListOffset: 0,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	user, err := s.storage.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("can't create user: %w", err)
	}

	// Grant initial regular tokens to new user
	initialTransaction := &domain.Transaction{
		UserID:          user.ID,
		TokenType:       domain.TokenTypeRegular,
		Amount:          initialTokenGrant,
		TransactionType: domain.TransactionTypeInitialCredit,
		Description:     stringPtr("Initial welcome tokens"),
		CreatedAt:       now,
	}

	_, err = s.storage.CreateTransaction(ctx, initialTransaction)
	if err != nil {
		s.logger.WarnContext(ctx, "failed to grant initial tokens to new user",
			slog.String("external_id", user.ExternalID),
			slog.String("error", err.Error()))
	}

	s.logger.InfoContext(ctx, "created new user with initial tokens",
		slog.String("external_id", user.ExternalID),
		slog.Int64("initial_tokens", initialTokenGrant))

	return user, nil
}

func (s *UpdateService) HandleConversationState(ctx context.Context, user *domain.User, update domain.Update) error {
	// Check if user sent "back to menu" text
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonBackToMenu) {
		// Clear any pending messages when going back to menu
		if clearErr := s.queue.ClearPendingMessages(ctx, user.ExternalID); clearErr != nil {
			s.logger.WarnContext(ctx, "failed to clear pending messages on menu transition",
				slog.String("error", clearErr.Error()))
		}
		return s.transitionToMenu(ctx, user)
	}

	// Check if user clicked web search toggle button (with or without red cross)
	webSearchOnText := i18n.GetString(user.Language, i18n.ButtonWebSearchOn)
	webSearchOffText := i18n.GetString(user.Language, i18n.ButtonWebSearchOff)
	redCrossOffText := "❌ " + webSearchOffText

	if update.MessageText == webSearchOnText ||
		update.MessageText == webSearchOffText ||
		update.MessageText == redCrossOffText {
		return s.handleWebSearchToggle(ctx, user)
	}

	// Check if user clicked subscription button
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonSubscription) {
		return s.handleSubscriptionBuyCallback(ctx, user)
	}

	// Handle regular conversation message with concatenation
	if update.MessageText != "" && update.MessageText != i18n.GetString(user.Language, i18n.ButtonBackToMenu) {
		return s.handleConversationMessageWithConcatenation(ctx, user, update)
	}

	// Send back to menu button if no message text
	content := domain.MessageContent{
		Text: i18n.GetString(user.Language, i18n.ConversationModePrompt),
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: i18n.GetString(user.Language, i18n.ButtonBackToMenu),
					},
				},
			},
			Resize:  true,
			OneTime: false,
		},
	}

	_, err := s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
	return err
}

func (s *UpdateService) handleWebSearchToggle(ctx context.Context, user *domain.User) error {
	// Check if user has active subscription (required for web search)
	activeSubscription, err := s.storage.GetActiveSubscriptionByUserID(ctx, user.ID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return fmt.Errorf("failed to check subscription status: %w", err)
	}

	// Only allow web search toggle if user has active subscription
	if activeSubscription == nil {
		// Send subscription requirement message with subscribe button
		subscriptionMessage := i18n.GetString(user.Language, i18n.WebSearchSubscriptionRequired)
		content := domain.MessageContent{
			Text: subscriptionMessage,
			ReplyKeyboard: &domain.ReplyKeyboard{
				Buttons: [][]domain.KeyboardButton{
					{
						{Text: i18n.GetString(user.Language, i18n.ButtonSubscription)},
					},
					{
						{Text: i18n.GetString(user.Language, i18n.ButtonBackToMenu)},
					},
				},
				Resize:  true,
				OneTime: false,
			},
		}
		_, sendErr := s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
		return sendErr
	}

	// Toggle the web search enabled state
	user.WebSearchEnabled = !user.WebSearchEnabled

	// Update in database
	err = s.storage.UpdateUserWebSearchEnabled(ctx, user.ID, user.WebSearchEnabled)
	if err != nil {
		return fmt.Errorf("can't update web search state: %w", err)
	}

	// Send updated conversation view with new button state - the unified message will include web search status
	return s.transitionToConversation(ctx, user, nil)
}

func (s *UpdateService) shouldQueueMessage(user *domain.User, update domain.Update) bool {
	return user.CurrentStep == domain.UserStateConversation &&
		update.MessageText != "" &&
		update.MessageText != i18n.GetString(user.Language, i18n.ButtonBackToMenu)
}

func (s *UpdateService) handleMessageQueueing(
	ctx context.Context,
	user *domain.User,
	update domain.Update,
) (bool, error) {
	// Check if user is currently processing
	processing, processingErr := s.queue.IsProcessing(ctx, user.ExternalID)
	if processingErr != nil {
		s.logger.WarnContext(ctx, "failed to check processing status",
			slog.String("error", processingErr.Error()))
		// Continue without queueing on error
		return false, nil
	}

	if !processing {
		return false, nil
	}

	// Notify user that message is queued
	queueLength, _ := s.queue.GetQueueLength(ctx, user.ExternalID)
	queueMessage := fmt.Sprintf(
		i18n.GetString(user.Language, i18n.QueueMessageQueued),
		queueLength+1,
	)
	notificationID, sendErr := s.sender.SendMessage(ctx, user.ExternalID, queueMessage)
	if sendErr != nil {
		s.logger.WarnContext(ctx, "failed to send queue notification",
			slog.String("error", sendErr.Error()))
		notificationID = "" // Continue without notification ID
	}

	// User is processing, queue the update with notification ID
	if enqueueErr := s.queue.EnqueueWithNotification(ctx, user.ExternalID, update, notificationID); enqueueErr != nil {
		s.logger.ErrorContext(ctx, "failed to enqueue update",
			slog.String("error", enqueueErr.Error()))
		return false, fmt.Errorf("failed to queue message: %w", enqueueErr)
	}

	return true, nil
}

func (s *UpdateService) HandleCallbackQuery(ctx context.Context, callbackQuery domain.CallbackQuery) (err error) {
	// Add panic recovery to prevent crashes during callback handling
	defer func() {
		if r := recover(); r != nil {
			stackTrace := debug.Stack()
			s.logger.ErrorContext(ctx, "Panic occurred while handling callback query",
				"panic", r,
				"stack_trace", string(stackTrace),
				"user_id", callbackQuery.ExternalUserID,
				"callback_data", callbackQuery.Data)

			// Convert panic to error
			err = fmt.Errorf("panic occurred while handling callback query: %v", r)
		}
	}()

	// Get user for callback (don't create if not exists)
	user, err := s.storage.GetUserByExternalUserID(ctx, callbackQuery.ExternalUserID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			s.logger.ErrorContext(ctx, "user not found for callback query",
				slog.String("external_user_id", callbackQuery.ExternalUserID),
				slog.String("callback_data", callbackQuery.Data))
			return nil
		}
		return fmt.Errorf("can't get user for callback: %w", err)
	}

	// Handle callback based on data
	switch callbackQuery.Data {
	case "subscription_buy_monthly":
		return s.handleSubscriptionBuyCallback(ctx, user)
	case "back_to_menu":
		return s.transitionToMenu(ctx, user)
	default:
		s.logger.WarnContext(ctx, "unknown callback data", slog.String("data", callbackQuery.Data))
		return nil
	}
}

func (s *UpdateService) HandlePreCheckoutQuery(ctx context.Context, query domain.PreCheckoutQuery) (err error) {
	// Add panic recovery to prevent crashes during pre-checkout handling
	defer func() {
		if r := recover(); r != nil {
			stackTrace := debug.Stack()
			s.logger.ErrorContext(ctx, "Panic occurred while handling pre-checkout query",
				"panic", r,
				"stack_trace", string(stackTrace),
				"user_id", query.ExternalUserID,
				"invoice_payload", query.InvoicePayload)

			// Convert panic to error
			err = fmt.Errorf("panic occurred while handling pre-checkout query: %v", r)
		}
	}()

	// Get payment record by invoice payload
	dbPayment, err := s.storage.GetPaymentByInvoicePayload(ctx, query.InvoicePayload)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to find payment for pre-checkout",
			slog.String("invoice_payload", query.InvoicePayload),
			slog.String("error", err.Error()))
		return s.sender.AnswerPreCheckoutQuery(ctx, query.ID, false, "Payment not found")
	}

	// Check if user already has an active subscription
	subscription, err := s.storage.GetActiveSubscriptionByUserID(ctx, dbPayment.UserID)
	if err == nil && subscription.IsActive() {
		s.logger.InfoContext(ctx, "user already has active subscription, declining pre-checkout",
			slog.Int64("payment_id", dbPayment.ID),
			slog.String("user_id", query.ExternalUserID),
			slog.Int64("subscription_id", subscription.ID))
		return s.sender.AnswerPreCheckoutQuery(ctx, query.ID, false, "You already have an active subscription")
	}

	// Update payment status to pre_checkout
	_, err = s.storage.UpdatePaymentStatus(ctx, dbPayment.ID, domain.PaymentStatusPreCheckout, nil, nil)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to update payment status to pre_checkout",
			slog.Int64("payment_id", dbPayment.ID),
			slog.String("error", err.Error()))
	}

	// Always approve the pre-checkout if no active subscription
	s.logger.InfoContext(ctx, "pre-checkout query approved",
		slog.Int64("payment_id", dbPayment.ID),
		slog.String("user_id", query.ExternalUserID))

	return s.sender.AnswerPreCheckoutQuery(ctx, query.ID, true, "")
}

func (s *UpdateService) HandleSuccessfulPayment(ctx context.Context, payment domain.SuccessfulPayment) (err error) {
	// Add panic recovery to prevent crashes during payment handling
	defer func() {
		if r := recover(); r != nil {
			stackTrace := debug.Stack()
			s.logger.ErrorContext(ctx, "Panic occurred while handling successful payment",
				"panic", r,
				"stack_trace", string(stackTrace),
				"user_id", payment.ExternalUserID,
				"invoice_payload", payment.InvoicePayload,
				"telegram_charge_id", payment.TelegramChargeID)

			// Convert panic to error
			err = fmt.Errorf("panic occurred while handling successful payment: %v", r)
		}
	}()

	// Get user
	user, err := s.storage.GetUserByExternalUserID(ctx, payment.ExternalUserID)
	if err != nil {
		return fmt.Errorf("can't get user for payment: %w", err)
	}

	// Get payment record by invoice payload
	dbPayment, err := s.storage.GetPaymentByInvoicePayload(ctx, payment.InvoicePayload)
	if err != nil {
		return fmt.Errorf("can't find payment with invoice payload %s: %w", payment.InvoicePayload, err)
	}

	// Update payment status to paid with Telegram charge IDs
	_, err = s.storage.UpdatePaymentStatus(
		ctx,
		dbPayment.ID,
		domain.PaymentStatusPaid,
		&payment.TelegramChargeID,
		&payment.ProviderChargeID,
	)
	if err != nil {
		return fmt.Errorf("can't update payment status: %w", err)
	}

	// Grant tokens for subscription
	err = s.grantSubscriptionTokens(ctx, dbPayment)
	if err != nil {
		return fmt.Errorf("can't grant subscription tokens: %w", err)
	}

	// Create subscription record
	now := time.Now()
	subscription := &domain.Subscription{
		UserID:           dbPayment.UserID,
		PaymentID:        dbPayment.ID,
		SubscriptionType: dbPayment.SubscriptionType,
		ValidFrom:        now,
		ValidTo:          now.AddDate(0, 1, 0), // Add 1 month
		Status:           domain.SubscriptionStatusActive,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	_, err = s.storage.CreateSubscription(ctx, subscription)
	if err != nil {
		return fmt.Errorf("can't create subscription record: %w", err)
	}

	// Send success message
	successMsg := i18n.GetString(user.Language, i18n.SubscriptionSuccess)
	_, err = s.sender.SendMessage(ctx, user.ExternalID, successMsg)
	if err != nil {
		s.logger.WarnContext(ctx, "failed to send payment success message", slog.String("error", err.Error()))
	}

	return nil
}
