package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vladimish/talk/internal/domain"
)

// Conversation list handlers.
func (s *UpdateService) transitionToConversationList(ctx context.Context, user *domain.User) error {
	conversationListState := domain.UserStateConversationList
	user.CurrentStep = conversationListState

	err := s.storage.UpdateUserCurrentStep(ctx, user.ID, conversationListState)
	if err != nil {
		return fmt.Errorf("can't update user state: %w", err)
	}

	return s.showConversationList(ctx, user)
}

func (s *UpdateService) handleConversationListState(ctx context.Context, user *domain.User, update domain.Update) error {
	// Check if user sent "back to menu" text
	if update.MessageText == domain.ButtonBackToMenu {
		return s.transitionToMenu(ctx, user)
	}

	// Check if user sent "new conversation" text
	if update.MessageText == domain.ButtonNewConversation {
		return s.createNewConversation(ctx, user)
	}

	// Check if user selected an existing conversation
	conversations, err := s.storage.GetConversationsByUserID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("can't get conversations: %w", err)
	}

	for _, conversationID := range conversations {
		if update.MessageText == fmt.Sprintf("💬 %s", conversationID[:8]) {
			return s.selectConversation(ctx, user, conversationID)
		}
	}

	// If no valid selection, show conversation list again
	return s.showConversationList(ctx, user)
}

func (s *UpdateService) showConversationList(ctx context.Context, user *domain.User) error {
	conversations, err := s.storage.GetConversationsByUserID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("can't get conversations: %w", err)
	}

	var buttons [][]domain.KeyboardButton

	// Add "New Conversation" button at the top
	buttons = append(buttons, []domain.KeyboardButton{
		{Text: domain.ButtonNewConversation},
	})

	// Add existing conversations
	for _, conversationID := range conversations {
		buttons = append(buttons, []domain.KeyboardButton{
			{Text: fmt.Sprintf("💬 %s", conversationID[:8])},
		})
	}

	// Add back button
	buttons = append(buttons, []domain.KeyboardButton{
		{Text: domain.ButtonBackToMenu},
	})

	text := "💬 Select a conversation:"
	if len(conversations) == 0 {
		text = "💬 No previous conversations. Start a new one:"
	}

	content := domain.MessageContent{
		Text: text,
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: buttons,
			Resize:  true,
			OneTime: true,
		},
	}

	_, err = s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
	return err
}

func (s *UpdateService) createNewConversation(ctx context.Context, user *domain.User) error {
	// Generate new conversation UUID
	conversationID := uuid.New().String()

	// Update user's current conversation
	user.CurrentConversation = &conversationID
	err := s.storage.UpdateUserCurrentConversation(ctx, user.ID, &conversationID)
	if err != nil {
		return fmt.Errorf("can't update user conversation: %w", err)
	}

	// Transition to conversation state
	return s.transitionToConversation(ctx, user)
}

func (s *UpdateService) selectConversation(ctx context.Context, user *domain.User, conversationID string) error {
	// Update user's current conversation
	user.CurrentConversation = &conversationID
	err := s.storage.UpdateUserCurrentConversation(ctx, user.ID, &conversationID)
	if err != nil {
		return fmt.Errorf("can't update user conversation: %w", err)
	}

	// Transition to conversation state
	return s.transitionToConversation(ctx, user)
}
