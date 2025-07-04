package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/storage"
	"github.com/vladimish/talk/pkg/i18n"
)

// HandleMenuState handles user interactions in the menu state.
func (s *UpdateService) HandleMenuState(ctx context.Context, user *domain.User, update domain.Update) error {
	// Check if user sent "start conversation" text
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonStartConversation) {
		return s.transitionToConversationList(ctx, user)
	}

	// Check if user sent "select model" text
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonModelSelect) {
		return s.transitionToModelSelect(ctx, user)
	}

	// Check if user sent "settings" text
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonSettings) {
		return s.transitionToSettings(ctx, user)
	}

	// Check if user sent "profile" text
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonProfile) {
		return s.transitionToProfile(ctx, user)
	}

	// Check if user sent "subscription" text
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonSubscription) {
		return s.handleSubscriptionBuyCallback(ctx, user)
	}

	// If message doesn't match any menu option, create new conversation and process the message
	if update.MessageText != "" {
		return s.createNewConversationFromMenu(ctx, user, update)
	}

	// Send menu with keyboard
	return s.sendMenu(ctx, user, i18n.GetString(user.Language, i18n.MenuWelcome))
}

func (s *UpdateService) sendMenu(ctx context.Context, user *domain.User, text string) error {
	content := domain.MessageContent{
		Text:         text,
		IsPersistent: true,
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: i18n.GetString(user.Language, i18n.ButtonStartConversation),
					},
				},
				{
					{
						Text: i18n.GetString(user.Language, i18n.ButtonModelSelect),
					},
					{
						Text: i18n.GetString(user.Language, i18n.ButtonProfile),
					},
				},
				{
					{
						Text: i18n.GetString(user.Language, i18n.ButtonSubscription),
					},
					{
						Text: i18n.GetString(user.Language, i18n.ButtonSettings),
					},
				},
			},
			Resize:  true,
			OneTime: true,
		},
	}

	_, err := s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
	return err
}

func (s *UpdateService) transitionToConversation(
	ctx context.Context,
	user *domain.User,
	replyToMessageID *int64,
) error {
	conversationState := domain.UserStateConversation
	user.CurrentStep = conversationState

	err := s.storage.UpdateUserCurrentStep(ctx, user.ID, conversationState)
	if err != nil {
		return fmt.Errorf("can't update user state: %w", err)
	}

	// Determine conversation start/resume message based on web search status
	text, err := s.getConversationStartMessage(ctx, user, replyToMessageID)
	if err != nil {
		return err
	}

	// Create keyboard buttons
	var buttons [][]domain.KeyboardButton

	// Add web search button if model supports it
	if s.shouldShowWebSearchButton(user) {
		webSearchButton, webSearchErr := s.getWebSearchButton(ctx, user)
		if webSearchErr != nil {
			return webSearchErr
		}
		buttons = append(buttons, []domain.KeyboardButton{webSearchButton})
	}

	// Add back to menu button
	buttons = append(buttons, []domain.KeyboardButton{
		{Text: i18n.GetString(user.Language, i18n.ButtonBackToMenu)},
	})

	content := domain.MessageContent{
		Text:             text,
		ReplyToMessageID: replyToMessageID,
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: buttons,
			Resize:  true,
			OneTime: false,
		},
	}

	_, err = s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
	return err
}

func (s *UpdateService) transitionToMenu(ctx context.Context, user *domain.User) error {
	menuState := domain.UserStateMenu
	user.CurrentStep = menuState

	err := s.storage.UpdateUserCurrentStep(ctx, user.ID, menuState)
	if err != nil {
		return fmt.Errorf("can't update user state: %w", err)
	}

	return s.sendMenu(ctx, user, i18n.GetString(user.Language, i18n.MenuBackToMain))
}

// shouldShowWebSearchButton checks if the web search button should be shown.
func (s *UpdateService) shouldShowWebSearchButton(user *domain.User) bool {
	currentModel := domain.GetModelByID(user.SelectedModel)
	return currentModel != nil && currentModel.WebSearch
}

// getWebSearchButton returns the web search button for the current model.
func (s *UpdateService) getWebSearchButton(ctx context.Context, user *domain.User) (domain.KeyboardButton, error) {
	// Check if user has active subscription (for button state)
	activeSubscription, subErr := s.storage.GetActiveSubscriptionByUserID(ctx, user.ID)
	if subErr != nil && !errors.Is(subErr, storage.ErrNotFound) {
		return domain.KeyboardButton{}, fmt.Errorf("failed to check subscription status: %w", subErr)
	}
	hasActiveSubscription := activeSubscription != nil

	// Always show web search toggle button for models that support it
	var webSearchButtonText string
	if hasActiveSubscription {
		// User has subscription - show normal toggle
		if user.WebSearchEnabled {
			webSearchButtonText = i18n.GetString(user.Language, i18n.ButtonWebSearchOn)
		} else {
			webSearchButtonText = i18n.GetString(user.Language, i18n.ButtonWebSearchOff)
		}
	} else {
		// User has no subscription - show with red cross
		webSearchButtonText = "❌ " + i18n.GetString(user.Language, i18n.ButtonWebSearchOff)
	}

	return domain.KeyboardButton{Text: webSearchButtonText}, nil
}

// getConversationStartMessage returns the appropriate conversation message based on web search status.
func (s *UpdateService) getConversationStartMessage(
	ctx context.Context,
	user *domain.User,
	replyToMessageID *int64,
) (string, error) {
	isResuming := replyToMessageID != nil
	webSearchActive, err := s.isWebSearchActive(ctx, user)
	if err != nil {
		return "", err
	}

	if webSearchActive {
		return s.getWebSearchOnMessage(user.Language, isResuming), nil
	}
	return s.getWebSearchOffMessage(user.Language, isResuming), nil
}

// isWebSearchActive checks if web search is active for the user.
func (s *UpdateService) isWebSearchActive(ctx context.Context, user *domain.User) (bool, error) {
	currentModel := domain.GetModelByID(user.SelectedModel)
	if currentModel == nil || !currentModel.WebSearch || !user.WebSearchEnabled {
		return false, nil
	}

	activeSubscription, err := s.storage.GetActiveSubscriptionByUserID(ctx, user.ID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return false, fmt.Errorf("failed to check subscription status: %w", err)
	}
	return activeSubscription != nil, nil
}

// getWebSearchOnMessage returns the message when web search is on.
func (s *UpdateService) getWebSearchOnMessage(language string, isResuming bool) string {
	if isResuming {
		return i18n.GetString(language, i18n.ConversationResumedWebSearchOn)
	}
	return i18n.GetString(language, i18n.ConversationStartedWebSearchOn)
}

// getWebSearchOffMessage returns the message when web search is off.
func (s *UpdateService) getWebSearchOffMessage(language string, isResuming bool) string {
	if isResuming {
		return i18n.GetString(language, i18n.ConversationResumedWebSearchOff)
	}
	return i18n.GetString(language, i18n.ConversationStartedWebSearchOff)
}

// createNewConversationFromMenu creates a new conversation from menu and processes the message.
func (s *UpdateService) createNewConversationFromMenu(
	ctx context.Context,
	user *domain.User,
	update domain.Update,
) error {
	// Create new conversation with temporary name
	conversationName := defaultConversationName
	conversation, err := s.storage.CreateConversation(ctx, &domain.Conversation{
		Name:      conversationName,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("can't create conversation: %w", err)
	}

	// Update user's current conversation
	user.CurrentConversationID = &conversation.ID
	err = s.storage.UpdateUserCurrentConversationID(ctx, user.ID, &conversation.ID)
	if err != nil {
		return fmt.Errorf("can't update user conversation: %w", err)
	}

	// Transition to conversation state with the conversation keyboard
	if transitionErr := s.transitionToConversation(ctx, user, nil); transitionErr != nil {
		return fmt.Errorf("can't transition to conversation: %w", transitionErr)
	}

	// Process the message in the new conversation
	return s.handleConversationMessage(ctx, user, update)
}
