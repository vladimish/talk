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
		return s.transitionToConversationList(ctx, user)
	}

	// Check if user sent "select model" text
	if update.MessageText == domain.ButtonModelSelect {
		return s.transitionToModelSelect(ctx, user)
	}

	// Send menu with keyboard
	return s.sendMenu(ctx, user, "Welcome! Choose an option:")
}

func (s *UpdateService) sendMenu(ctx context.Context, user *domain.User, text string) error {
	content := domain.MessageContent{
		Text: text,
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: domain.ButtonStartConversation,
					},
				},
				{
					{
						Text: domain.ButtonModelSelect,
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

	text := "🗣️ Conversation started! Send me a message and I'll respond. You can always go back to the menu."
	if replyToMessageID != nil {
		text = "🗣️ Conversation resumed! Send me a message and I'll respond. You can always go back to the menu."
	}

	content := domain.MessageContent{
		Text:             text,
		ReplyToMessageID: replyToMessageID,
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

	return s.sendMenu(ctx, user, "🏠 Back to main menu. Choose an option:")
}
