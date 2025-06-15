package service

import (
	"context"
	"fmt"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/pkg/i18n"
)

// Model selection handlers.
func (s *UpdateService) transitionToModelSelect(ctx context.Context, user *domain.User) error {
	modelSelectState := domain.UserStateModelSelect
	user.CurrentStep = modelSelectState

	err := s.storage.UpdateUserCurrentStep(ctx, user.ID, modelSelectState)
	if err != nil {
		return fmt.Errorf("can't update user state: %w", err)
	}

	return s.showModelSelection(ctx, user)
}

func (s *UpdateService) HandleModelSelectState(ctx context.Context, user *domain.User, update domain.Update) error {
	// Check if user sent "back to menu" text
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonBackToMenu) {
		return s.transitionToMenu(ctx, user)
	}

	// Get user balance and subscription to compare with displayed names
	balance, err := s.storage.GetUserTokenBalance(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("can't get user token balance: %w", err)
	}

	activeSubscription, err := s.storage.GetActiveSubscriptionByUserID(ctx, user.ID)
	hasActiveSubscription := err == nil && activeSubscription != nil

	// Check if user selected a model by comparing with display names
	for _, model := range domain.AvailableModels {
		displayName := model.GetDisplayNameForUser(user.Language, hasActiveSubscription, *balance)
		if update.MessageText == displayName {
			// Check if user can actually use this model
			if !domain.CanUserUseModel(&model, hasActiveSubscription, *balance) {
				// Send error message for unavailable model
				insufficientText := i18n.GetString(user.Language, i18n.ProfileInsufficientTokens)
				_, sendErr := s.sender.SendMessageWithContent(ctx, user.ExternalID, domain.MessageContent{
					Text: insufficientText,
				})
				if sendErr != nil {
					return sendErr
				}
				return s.showModelSelection(ctx, user)
			}
			return s.selectModel(ctx, user, model.ID)
		}
	}

	// If no valid selection, show model selection again
	return s.showModelSelection(ctx, user)
}

func (s *UpdateService) showModelSelection(ctx context.Context, user *domain.User) error {
	// Get user balance and subscription status
	balance, err := s.storage.GetUserTokenBalance(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("can't get user token balance: %w", err)
	}

	activeSubscription, err := s.storage.GetActiveSubscriptionByUserID(ctx, user.ID)
	hasActiveSubscription := err == nil && activeSubscription != nil

	// Create buttons for each model with user-specific display names
	var modelButtons [][]domain.KeyboardButton
	for _, model := range domain.AvailableModels {
		displayName := model.GetDisplayNameForUser(user.Language, hasActiveSubscription, *balance)
		modelButtons = append(modelButtons, []domain.KeyboardButton{
			{Text: displayName},
		})
	}

	// Add back button
	modelButtons = append(modelButtons, []domain.KeyboardButton{
		{Text: i18n.GetString(user.Language, i18n.ButtonBackToMenu)},
	})

	titleText := i18n.GetString(user.Language, i18n.ModelSelectTitle)

	// Get current model display name with emojis
	currentModelDisplay := user.SelectedModel
	if currentModel := domain.GetModelByID(user.SelectedModel); currentModel != nil {
		currentModelDisplay = currentModel.GetDisplayNameWithEmojis(user.Language)
	}

	content := domain.MessageContent{
		Text:         fmt.Sprintf("%s\n\nCurrent model: %s\n\nChoose a model:", titleText, currentModelDisplay),
		IsPersistent: true,
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: modelButtons,
			Resize:  true,
			OneTime: true,
		},
	}

	_, err = s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
	return err
}

func (s *UpdateService) selectModel(ctx context.Context, user *domain.User, model string) error {
	// Update user's selected model in memory
	user.SelectedModel = model

	// Update in database
	err := s.storage.UpdateUserSelectedModel(ctx, user.ID, model)
	if err != nil {
		return fmt.Errorf("can't update user model: %w", err)
	}

	// Transition back to menu
	menuState := domain.UserStateMenu
	user.CurrentStep = menuState

	err = s.storage.UpdateUserCurrentStep(ctx, user.ID, menuState)
	if err != nil {
		return fmt.Errorf("can't update user state: %w", err)
	}

	// Get model display name for success message
	modelDisplayName := model
	if selectedModel := domain.GetModelByID(model); selectedModel != nil {
		modelDisplayName = selectedModel.GetDisplayNameWithEmojis(user.Language)
	}

	successMessage := fmt.Sprintf("%s %s", i18n.GetString(user.Language, i18n.ModelUpdateSuccess), modelDisplayName)
	return s.sendMenu(ctx, user, successMessage+"\n\n"+i18n.GetString(user.Language, i18n.MenuBackToMain))
}
