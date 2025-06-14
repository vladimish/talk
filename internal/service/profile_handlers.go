package service

import (
	"context"
	"fmt"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/pkg/i18n"
)

// Profile state handlers.
func (s *UpdateService) handleProfileState(ctx context.Context, user *domain.User, update domain.Update) error {
	// Check if user sent "back to menu" text
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonBackToMenu) {
		return s.transitionToMenu(ctx, user)
	}

	// Show profile
	return s.showProfile(ctx, user)
}

func (s *UpdateService) transitionToProfile(ctx context.Context, user *domain.User) error {
	profileState := domain.UserStateProfile
	user.CurrentStep = profileState

	err := s.storage.UpdateUserCurrentStep(ctx, user.ID, profileState)
	if err != nil {
		return fmt.Errorf("can't update user state: %w", err)
	}

	return s.showProfile(ctx, user)
}

func (s *UpdateService) showProfile(ctx context.Context, user *domain.User) error {
	// Get user token balance
	balance, err := s.storage.GetUserTokenBalance(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("can't get user token balance: %w", err)
	}

	// Format profile message
	title := i18n.GetString(user.Language, i18n.ProfileTitle)
	tokenBalanceText := i18n.GetString(user.Language, i18n.ProfileTokenBalance)
	premiumText := fmt.Sprintf(i18n.GetString(user.Language, i18n.ProfilePremiumTokens), balance.PremiumBalance)
	regularText := fmt.Sprintf(i18n.GetString(user.Language, i18n.ProfileRegularTokens), balance.RegularBalance)

	profileText := fmt.Sprintf("%s\n\n%s\n%s\n%s", title, tokenBalanceText, premiumText, regularText)

	content := domain.MessageContent{
		Text:         profileText,
		IsPersistent: true,
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{Text: i18n.GetString(user.Language, i18n.ButtonBackToMenu)},
				},
			},
			Resize:  true,
			OneTime: true,
		},
	}

	_, err = s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
	return err
}
