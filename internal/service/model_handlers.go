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

	// Check if user selected a model
	for _, model := range domain.AvailableModels {
		if update.MessageText == model.DisplayName {
			return s.selectModel(ctx, user, model.ID)
		}
	}

	// If no valid selection, show model selection again
	return s.showModelSelection(ctx, user)
}

func (s *UpdateService) showModelSelection(ctx context.Context, user *domain.User) error {
	// Create buttons for each model
	var modelButtons [][]domain.KeyboardButton
	for _, model := range domain.AvailableModels {
		modelButtons = append(modelButtons, []domain.KeyboardButton{
			{Text: model.DisplayName},
		})
	}

	// Add back button
	modelButtons = append(modelButtons, []domain.KeyboardButton{
		{Text: i18n.GetString(user.Language, i18n.ButtonBackToMenu)},
	})

	titleText := i18n.GetString(user.Language, i18n.ModelSelectTitle)

	// Get current model display name
	currentModelDisplay := user.SelectedModel
	if currentModel := domain.GetModelByID(user.SelectedModel); currentModel != nil {
		currentModelDisplay = currentModel.DisplayName
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

	_, err := s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
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
		modelDisplayName = selectedModel.DisplayName
	}

	successMessage := fmt.Sprintf("%s %s", i18n.GetString(user.Language, i18n.ModelUpdateSuccess), modelDisplayName)
	return s.sendMenu(ctx, user, successMessage+"\n\n"+i18n.GetString(user.Language, i18n.MenuBackToMain))
}
