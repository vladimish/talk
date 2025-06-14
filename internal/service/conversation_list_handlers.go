package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/storage"
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

func (s *UpdateService) handleConversationListState(
	ctx context.Context,
	user *domain.User,
	update domain.Update,
) error {
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

	for _, conversation := range conversations {
		if update.MessageText == fmt.Sprintf("💬 %s", conversation.Name) {
			return s.selectConversation(ctx, user, conversation.ID)
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
	for _, conversation := range conversations {
		buttons = append(buttons, []domain.KeyboardButton{
			{Text: fmt.Sprintf("💬 %s", conversation.Name)},
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
	// Create new conversation
	conversationName := fmt.Sprintf("Conversation %s", time.Now().Format("Jan 2 15:04"))
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

	// Transition to conversation state (no reply since it's a new conversation)
	return s.transitionToConversation(ctx, user, nil)
}

func (s *UpdateService) selectConversation(ctx context.Context, user *domain.User, conversationID int64) error {
	// Update user's current conversation
	user.CurrentConversationID = &conversationID
	err := s.storage.UpdateUserCurrentConversationID(ctx, user.ID, &conversationID)
	if err != nil {
		return fmt.Errorf("can't update user conversation: %w", err)
	}

	// Get the latest message from the conversation to reply to
	var replyToMessageID *int64
	latestMessage, err := s.storage.GetLatestMessageByConversationID(ctx, conversationID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return fmt.Errorf("can't get latest message: %w", err)
	}
	if latestMessage != nil && latestMessage.ID <= int64(^uint32(0)>>1) {
		// Get the foreign (Telegram) message ID
		var foreignMessageID int32
		foreignMessageID, err = s.storage.GetForeignMessageByMessageID(ctx, int32(latestMessage.ID)) //nolint:gosec
		if err != nil && !errors.Is(err, storage.ErrNotFound) {
			s.logger.WarnContext(ctx, "failed to get foreign message ID",
				slog.Int64("message_id", latestMessage.ID),
				slog.String("error", err.Error()))
		} else if err == nil {
			// Convert int32 to int64 for the reply
			foreignID := int64(foreignMessageID)
			replyToMessageID = &foreignID
		}
	}

	// Transition to conversation state
	return s.transitionToConversation(ctx, user, replyToMessageID)
}
