package service

import (
	"context"
	"fmt"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/pkg/i18n"
)

// Menu state handlers.
func (s *UpdateService) handleMenuState(ctx context.Context, user *domain.User, update domain.Update) error {
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

	// Send menu with keyboard
	return s.sendMenu(ctx, user, i18n.GetString(user.Language, i18n.MenuWelcome))
}

func (s *UpdateService) sendMenu(ctx context.Context, user *domain.User, text string) error {
	content := domain.MessageContent{
		Text: text,
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

	text := i18n.GetString(user.Language, i18n.ConversationStarted)
	if replyToMessageID != nil {
		text = i18n.GetString(user.Language, i18n.ConversationResumed)
	}

	content := domain.MessageContent{
		Text:             text,
		ReplyToMessageID: replyToMessageID,
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
