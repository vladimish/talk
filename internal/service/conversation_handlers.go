package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/vladimish/talk/internal/domain"
)

// Conversation handlers.
func (s *UpdateService) handleConversationMessage(ctx context.Context, user *domain.User, update domain.Update) error {
	// Save incoming message from user
	userMessage, err := s.storage.CreateMessage(ctx, &domain.Message{
		UserID: user.ID,
		MessageType: domain.MessageType{
			Text: update.MessageText,
		},
		SentBy:         domain.MessageSenderUser,
		ConversationID: user.CurrentConversationID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		return fmt.Errorf("can't save user message: %w", err)
	}

	if update.ExternalMessageID > 0 {
		if userMessage.ID > int64(^uint32(0)>>1) || int64(update.ExternalMessageID) > int64(^uint32(0)>>1) {
			s.logger.WarnContext(ctx, "message ID too large for foreign message mapping")
		} else {
			err = s.storage.CreateForeignMessage(ctx, int32(userMessage.ID), int32(update.ExternalMessageID)) //nolint:gosec
			if err != nil {
				s.logger.WarnContext(ctx, "failed to save foreign message mapping", slog.String("error", err.Error()))
			}
		}
	}

	// Get conversation history
	var messages []*domain.Message
	if user.CurrentConversationID != nil {
		messages, err = s.storage.GetMessagesByConversationID(ctx, *user.CurrentConversationID)
		if err != nil {
			return fmt.Errorf("can't get messages: %w", err)
		}
	}

	// Send typing indicator
	err = s.sender.SendTyping(ctx, user.ExternalID)
	if err != nil {
		s.logger.WarnContext(ctx, "failed to send typing indicator", slog.String("error", err.Error()))
	}

	systemPrompt := "You are a helpful assistant."
	tokenStream, err := s.completion.CompleteStream(ctx, user.SelectedModel, systemPrompt, messages)
	if err != nil {
		return fmt.Errorf("can't get completion: %w", err)
	}

	// Send initial empty message to get message ID
	messageID, err := s.sender.SendMessage(ctx, user.ExternalID, "\\.\\.\\.")
	if err != nil {
		return fmt.Errorf("can't send initial message: %w", err)
	}

	// Track all message IDs (starts with just one)
	messageIDs := []string{messageID}

	// Stream response tokens and update message
	var responseBuilder strings.Builder
	var previousContent string
	lastUpdate := time.Now()

	for token := range tokenStream {
		if token.Error != nil {
			return fmt.Errorf("completion stream error: %w", token.Error)
		}
		responseBuilder.WriteString(token.Content)

		// Update message at most once per second
		if time.Since(lastUpdate) >= time.Second/2 {
			currentContent := responseBuilder.String()

			// Send typing indicator to keep the conversation active
			err = s.sender.SendTyping(ctx, user.ExternalID)
			if err != nil {
				s.logger.WarnContext(ctx, "failed to send typing indicator", slog.String("error", err.Error()))
			}

			// Update all tracked messages with optimization
			updatedMessageIDs, updateErr := s.sender.UpdateMessages(
				ctx,
				user.ExternalID,
				messageIDs,
				previousContent,
				currentContent,
			)
			if updateErr != nil {
				return fmt.Errorf("can't update messages: %w", updateErr)
			}
			// Update our tracked message IDs
			messageIDs = updatedMessageIDs
			previousContent = currentContent
			lastUpdate = time.Now()
		}
	}

	// Send final update if needed
	if time.Since(lastUpdate) > 0 {
		finalContent := responseBuilder.String()
		updatedMessageIDs, finalUpdateErr := s.sender.UpdateMessages(
			ctx,
			user.ExternalID,
			messageIDs,
			previousContent,
			finalContent,
		)
		if finalUpdateErr != nil {
			return fmt.Errorf("can't update final messages: %w", finalUpdateErr)
		}
		// Update our tracked message IDs for final state
		messageIDs = updatedMessageIDs
	}

	responseText := responseBuilder.String()

	botMessage, err := s.storage.CreateMessage(ctx, &domain.Message{
		UserID: user.ID,
		MessageType: domain.MessageType{
			Text: responseText,
		},
		SentBy:         domain.MessageSenderBot,
		ConversationID: user.CurrentConversationID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		return fmt.Errorf("can't save bot message: %w", err)
	}

	// Save foreign message mappings for all bot messages
	for _, msgIDStr := range messageIDs {
		if msgIDStr == "" {
			continue
		}

		// Parse message ID from string to int
		msgID, parseErr := strconv.ParseInt(msgIDStr, 10, 32)
		if parseErr != nil {
			s.logger.WarnContext(ctx, "failed to parse message ID",
				slog.String("message_id", msgIDStr),
				slog.String("error", parseErr.Error()))
			continue
		}

		// Create foreign message mapping
		if botMessage.ID > int64(^uint32(0)>>1) || msgID > int64(^uint32(0)>>1) {
			s.logger.WarnContext(ctx, "message ID too large for foreign message mapping")
			continue
		}
		err = s.storage.CreateForeignMessage(ctx, int32(botMessage.ID), int32(msgID)) //nolint:gosec
		if err != nil {
			s.logger.WarnContext(ctx, "failed to save foreign message mapping for bot",
				slog.String("telegram_message_id", msgIDStr),
				slog.Int64("db_message_id", botMessage.ID),
				slog.String("error", err.Error()))
		} else {
			s.logger.DebugContext(ctx, "saved foreign message mapping",
				slog.String("telegram_message_id", msgIDStr),
				slog.Int64("db_message_id", botMessage.ID))
		}
	}

	return nil
}
