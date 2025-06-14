package service

import (
	"context"
	"fmt"

	"github.com/vladimish/talk/internal/domain"
)

// Settings state handlers.
func (s *UpdateService) handleSettingsState(ctx context.Context, user *domain.User, update domain.Update) error {
	// Check if user sent "back to menu" text
	if update.MessageText == domain.ButtonBackToMenu {
		return s.transitionToMenu(ctx, user)
	}

	// Check if user sent "language" text
	if update.MessageText == domain.ButtonLanguage {
		return s.transitionToLanguageSelect(ctx, user)
	}

	// Send settings with keyboard
	return s.sendSettings(ctx, user, "⚙️ Settings. Choose an option:")
}

func (s *UpdateService) handleLanguageSelectState(ctx context.Context, user *domain.User, update domain.Update) error {
	// Check if user sent "back to menu" text
	if update.MessageText == domain.ButtonBackToMenu {
		return s.transitionToSettings(ctx, user)
	}

	// Check if user selected a language
	if languageCode, exists := domain.SupportedLanguages[update.MessageText]; exists {
		return s.updateUserLanguage(ctx, user, languageCode)
	}

	// Send language selection with keyboard
	return s.sendLanguageSelection(ctx, user, "🌐 Select your language:")
}

func (s *UpdateService) sendSettings(ctx context.Context, user *domain.User, text string) error {
	content := domain.MessageContent{
		Text: text,
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: domain.ButtonLanguage,
					},
				},
				{
					{
						Text: domain.ButtonBackToMenu,
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

func (s *UpdateService) sendLanguageSelection(ctx context.Context, user *domain.User, text string) error {
	// Create buttons in 3 rows: English & Russian first, then others
	buttons := [][]domain.KeyboardButton{
		// First row: English and Russian
		{
			{Text: "🇺🇸 English"},
			{Text: "🇷🇺 Русский"},
		},
		// Second row: European languages
		{
			{Text: "🇪🇸 Español"},
			{Text: "🇫🇷 Français"},
			{Text: "🇩🇪 Deutsch"},
			{Text: "🇮🇹 Italiano"},
		},
		// Third row: Asian languages and Portuguese
		{
			{Text: "🇨🇳 中文"},
			{Text: "🇯🇵 日本語"},
			{Text: "🇰🇷 한국어"},
			{Text: "🇵🇹 Português"},
		},
		// Back button
		{
			{Text: domain.ButtonBackToMenu},
		},
	}

	content := domain.MessageContent{
		Text: text,
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: buttons,
			Resize:  true,
			OneTime: true,
		},
	}

	_, err := s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
	return err
}

func (s *UpdateService) transitionToSettings(ctx context.Context, user *domain.User) error {
	settingsState := domain.UserStateSettings
	user.CurrentStep = settingsState

	err := s.storage.UpdateUserCurrentStep(ctx, user.ID, settingsState)
	if err != nil {
		return fmt.Errorf("can't update user state: %w", err)
	}

	return s.sendSettings(ctx, user, "⚙️ Settings. Choose an option:")
}

func (s *UpdateService) transitionToLanguageSelect(ctx context.Context, user *domain.User) error {
	languageSelectState := domain.UserStateLanguageSelect
	user.CurrentStep = languageSelectState

	err := s.storage.UpdateUserCurrentStep(ctx, user.ID, languageSelectState)
	if err != nil {
		return fmt.Errorf("can't update user state: %w", err)
	}

	return s.sendLanguageSelection(ctx, user, "🌐 Select your language:")
}

func (s *UpdateService) updateUserLanguage(ctx context.Context, user *domain.User, languageCode string) error {
	// Update language in database
	err := s.storage.UpdateUserLanguage(ctx, user.ID, languageCode)
	if err != nil {
		return fmt.Errorf("can't update user language: %w", err)
	}

	// Update user object
	user.Language = languageCode

	// Send confirmation and return to settings
	confirmationText := "✅ Language updated successfully!"
	_, err = s.sender.SendMessage(ctx, user.ExternalID, confirmationText)
	if err != nil {
		return err
	}

	// Return to settings menu
	return s.transitionToSettings(ctx, user)
}
