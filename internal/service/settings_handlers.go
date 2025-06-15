package service

import (
	"context"
	"fmt"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/pkg/i18n"
)

// HandleSettingsState handles user interactions in the settings state.
func (s *UpdateService) HandleSettingsState(ctx context.Context, user *domain.User, update domain.Update) error {
	// Check if user sent "back to menu" text
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonBackToMenu) {
		return s.transitionToMenu(ctx, user)
	}

	// Check if user sent "language" text
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonLanguage) {
		return s.transitionToLanguageSelect(ctx, user)
	}

	// Send settings with keyboard
	return s.sendSettings(ctx, user, i18n.GetString(user.Language, i18n.SettingsTitle))
}

func (s *UpdateService) HandleLanguageSelectState(ctx context.Context, user *domain.User, update domain.Update) error {
	// Check if user sent "back to menu" text
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonBackToMenu) {
		return s.transitionToSettings(ctx, user)
	}

	// Check if user selected a language
	if languageCode, exists := i18n.LanguageCodes[update.MessageText]; exists {
		return s.updateUserLanguage(ctx, user, languageCode)
	}

	// Send language selection with keyboard
	return s.sendLanguageSelection(ctx, user, i18n.GetString(user.Language, i18n.LanguageSelectTitle))
}

func (s *UpdateService) sendSettings(ctx context.Context, user *domain.User, text string) error {
	content := domain.MessageContent{
		Text:         text,
		IsPersistent: true,
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: i18n.GetString(user.Language, i18n.ButtonLanguage),
					},
				},
				{
					{
						Text: i18n.GetString(user.Language, i18n.ButtonBackToMenu),
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
	// Create buttons in 5 rows: English, Russian, Spanish, Hindi in first row, then others
	buttons := [][]domain.KeyboardButton{
		// First row: English, Russian, Spanish, Hindi
		{
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangEnglish)},
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangRussian)},
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangSpanish)},
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangHindi)},
		},
		// Second row: European languages
		{
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangFrench)},
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangGerman)},
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangItalian)},
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangPortuguese)},
		},
		// Third row: Asian languages
		{
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangChinese)},
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangJapanese)},
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangKorean)},
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangArabic)},
		},
		// Fourth row: Additional languages
		{
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangArmenian)},
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangUkrainian)},
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangKazakh)},
			{Text: i18n.GetLanguageDisplayName(user.Language, i18n.LangKyrgyz)},
		},
		// Back button
		{
			{Text: i18n.GetString(user.Language, i18n.ButtonBackToMenu)},
		},
	}

	content := domain.MessageContent{
		Text:         text,
		IsPersistent: true,
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

	return s.sendSettings(ctx, user, i18n.GetString(user.Language, i18n.SettingsTitle))
}

func (s *UpdateService) transitionToLanguageSelect(ctx context.Context, user *domain.User) error {
	languageSelectState := domain.UserStateLanguageSelect
	user.CurrentStep = languageSelectState

	err := s.storage.UpdateUserCurrentStep(ctx, user.ID, languageSelectState)
	if err != nil {
		return fmt.Errorf("can't update user state: %w", err)
	}

	return s.sendLanguageSelection(ctx, user, i18n.GetString(user.Language, i18n.LanguageSelectTitle))
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
	confirmationText := i18n.GetString(user.Language, i18n.LanguageUpdateSuccess)
	_, err = s.sender.SendMessage(ctx, user.ExternalID, confirmationText)
	if err != nil {
		return err
	}

	// Return to settings menu
	return s.transitionToSettings(ctx, user)
}
