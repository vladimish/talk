package service

import (
	"context"
	"fmt"

	"github.com/vladimish/talk/internal/domain"
)

// Menu state handlers.
func (s *UpdateService) handleMenuState(ctx context.Context, user *domain.User, update domain.Update) error {
	// Check if user sent "start conversation" text
	if update.MessageText == domain.ButtonStartConversation {
		return s.transitionToConversation(ctx, user)
	}

	// Send menu with keyboard
	content := domain.MessageContent{
		Text: "Welcome! Choose an option:",
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: domain.ButtonStartConversation,
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

func (s *UpdateService) transitionToConversation(ctx context.Context, user *domain.User) error {
	conversationState := domain.UserStateConversation
	user.CurrentStep = conversationState

	err := s.storage.UpdateUserCurrentStep(ctx, user.ID, conversationState)
	if err != nil {
		return fmt.Errorf("can't update user state: %w", err)
	}

	content := domain.MessageContent{
		Text: "🗣️ Conversation started! Send me a message and I'll respond. You can always go back to the menu.",
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: domain.ButtonBackToMenu,
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

	content := domain.MessageContent{
		Text: "🏠 Back to main menu. Choose an option:",
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: domain.ButtonStartConversation,
					},
				},
			},
			Resize:  true,
			OneTime: true,
		},
	}

	_, err = s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
	return err
}
