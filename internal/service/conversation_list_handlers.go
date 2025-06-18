package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/storage"
	"github.com/vladimish/talk/pkg/i18n"
)

// Conversation list handlers.
func (s *UpdateService) transitionToConversationList(ctx context.Context, user *domain.User) error {
	conversationListState := domain.UserStateConversationList
	user.CurrentStep = conversationListState

	err := s.storage.UpdateUserCurrentStep(ctx, user.ID, conversationListState)
	if err != nil {
		return fmt.Errorf("can't update user state: %w", err)
	}

	// Reset pagination offset when entering conversation list
	user.ConversationListOffset = 0
	err = s.storage.UpdateUserConversationListOffset(ctx, user.ID, 0)
	if err != nil {
		return fmt.Errorf("can't update conversation list offset: %w", err)
	}

	return s.showConversationList(ctx, user)
}

func (s *UpdateService) HandleConversationListState(
	ctx context.Context,
	user *domain.User,
	update domain.Update,
) error {
	// Check if user sent "back to menu" text
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonBackToMenu) {
		return s.transitionToMenu(ctx, user)
	}

	// Check if user sent "new conversation" text
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonNewConversation) {
		return s.createNewConversation(ctx, user)
	}

	// Check navigation buttons
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonNextPage) {
		return s.handleNextPage(ctx, user)
	}
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonPrevPage) {
		return s.handlePrevPage(ctx, user)
	}

	// Check if user selected an existing conversation
	conversations, err := s.storage.GetConversationsByUserID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("can't get conversations: %w", err)
	}

	for _, conversation := range conversations {
		lastUpdate := s.formatDateTime(conversation.UpdatedAt)
		if update.MessageText == fmt.Sprintf("💬 %s, %s", conversation.Name, lastUpdate) {
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

	const pageSize = 5
	totalConversations := len(conversations)
	offset := user.ConversationListOffset

	// Validate and adjust offset
	if offset < 0 {
		offset = 0
	}
	if offset >= totalConversations && totalConversations > 0 {
		offset = ((totalConversations - 1) / pageSize) * pageSize
	}

	var buttons [][]domain.KeyboardButton

	// Add "New Conversation" button at the top
	buttons = append(buttons, []domain.KeyboardButton{
		{Text: i18n.GetString(user.Language, i18n.ButtonNewConversation)},
	})

	// Add conversations for current page
	start := offset
	end := offset + pageSize
	if end > totalConversations {
		end = totalConversations
	}

	for i := start; i < end; i++ {
		conversation := conversations[i]
		// Format the last update time according to user's language
		lastUpdate := s.formatDateTime(conversation.UpdatedAt)
		buttons = append(buttons, []domain.KeyboardButton{
			{Text: fmt.Sprintf("💬 %s, %s", conversation.Name, lastUpdate)},
		})
	}

	// Add navigation buttons if needed
	if totalConversations > pageSize {
		var navButtons []domain.KeyboardButton

		// Add previous button if not on first page
		if offset > 0 {
			prevText := i18n.GetString(user.Language, i18n.ButtonPrevPage)
			navButtons = append(navButtons, domain.KeyboardButton{Text: prevText})
		}

		// Add next button if not on last page
		if end < totalConversations {
			nextText := i18n.GetString(user.Language, i18n.ButtonNextPage)
			navButtons = append(navButtons, domain.KeyboardButton{Text: nextText})
		}

		if len(navButtons) > 0 {
			buttons = append(buttons, navButtons)
		}
	}

	// Add back button
	buttons = append(buttons, []domain.KeyboardButton{
		{Text: i18n.GetString(user.Language, i18n.ButtonBackToMenu)},
	})

	text := i18n.GetString(user.Language, i18n.ConversationListSelect)
	if totalConversations == 0 {
		text = i18n.GetString(user.Language, i18n.ConversationListEmpty)
	} else if totalConversations > pageSize {
		currentPage := (offset / pageSize) + 1
		totalPages := ((totalConversations - 1) / pageSize) + 1
		text = fmt.Sprintf(i18n.GetString(user.Language, i18n.ConversationListPageInfo), currentPage, totalPages)
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

	_, err = s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
	return err
}

func (s *UpdateService) createNewConversation(ctx context.Context, user *domain.User) error {
	// Create new conversation with temporary name
	conversationName := "New conversation"
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

func (s *UpdateService) handleNextPage(ctx context.Context, user *domain.User) error {
	conversations, err := s.storage.GetConversationsByUserID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("can't get conversations: %w", err)
	}

	const pageSize = 5
	totalConversations := len(conversations)
	newOffset := user.ConversationListOffset + pageSize

	// Don't go past the last page
	if newOffset >= totalConversations {
		return s.showConversationList(ctx, user)
	}

	user.ConversationListOffset = newOffset
	err = s.storage.UpdateUserConversationListOffset(ctx, user.ID, newOffset)
	if err != nil {
		return fmt.Errorf("can't update conversation list offset: %w", err)
	}

	return s.showConversationList(ctx, user)
}

func (s *UpdateService) handlePrevPage(ctx context.Context, user *domain.User) error {
	const pageSize = 5
	newOffset := user.ConversationListOffset - pageSize

	// Don't go before the first page
	if newOffset < 0 {
		newOffset = 0
	}

	user.ConversationListOffset = newOffset
	err := s.storage.UpdateUserConversationListOffset(ctx, user.ID, newOffset)
	if err != nil {
		return fmt.Errorf("can't update conversation list offset: %w", err)
	}

	return s.showConversationList(ctx, user)
}

// formatDateTime formats a timestamp in a unified dd.mm hh:mm format.
func (s *UpdateService) formatDateTime(t time.Time) string {
	return t.Format("02.01 15:04")
}
