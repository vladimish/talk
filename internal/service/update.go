package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/completion"
	"github.com/vladimish/talk/internal/port/filestorage"
	"github.com/vladimish/talk/internal/port/queue"
	"github.com/vladimish/talk/internal/port/sender"
	"github.com/vladimish/talk/internal/port/storage"
	"github.com/vladimish/talk/pkg/i18n"
)

const initialTokenGrant = 20

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

func (s *UpdateService) HandleUpdate(ctx context.Context, update domain.Update) error {
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

	switch currentState {
	case domain.UserStateMenu:
		return s.HandleMenuState(ctx, user, update)
	case domain.UserStateConversation:
		return s.HandleConversationState(ctx, user, update)
	case domain.UserStateModelSelect:
		return s.HandleModelSelectState(ctx, user, update)
	case domain.UserStateConversationList:
		return s.HandleConversationListState(ctx, user, update)
	case domain.UserStateSettings:
		return s.HandleSettingsState(ctx, user, update)
	case domain.UserStateLanguageSelect:
		return s.HandleLanguageSelectState(ctx, user, update)
	case domain.UserStateProfile:
		return s.HandleProfileState(ctx, user, update)
	default:
		// Default to menu state for unknown states
		return s.HandleMenuState(ctx, user, update)
	}
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

	// Handle regular conversation message
	if update.MessageText != "" && update.MessageText != i18n.GetString(user.Language, i18n.ButtonBackToMenu) {
		return s.handleConversationMessage(ctx, user, update)
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

func (s *UpdateService) HandleCallbackQuery(ctx context.Context, callbackQuery domain.CallbackQuery) error {
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

func (s *UpdateService) HandlePreCheckoutQuery(ctx context.Context, query domain.PreCheckoutQuery) error {
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

func (s *UpdateService) HandleSuccessfulPayment(ctx context.Context, payment domain.SuccessfulPayment) error {
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
