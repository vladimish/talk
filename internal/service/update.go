package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/completion"
	"github.com/vladimish/talk/internal/port/sender"
	"github.com/vladimish/talk/internal/port/storage"
	"log/slog"
	"strings"
	"time"
)

type UpdateService struct {
	logger     *slog.Logger
	storage    storage.Storage
	sender     sender.Sender
	completion completion.Completion
}

func NewUpdateService(logger *slog.Logger, storage storage.Storage, sender sender.Sender, completion completion.Completion) *UpdateService {
	return &UpdateService{
		logger:     logger,
		storage:    storage,
		sender:     sender,
		completion: completion,
	}
}

func (s *UpdateService) HandleUpdate(ctx context.Context, update domain.Update) error {
	user, err := s.storage.GetUserByExternalUserID(ctx, update.ExternalUserID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			language := update.UserLanguage
			now := time.Now()
			newUser := &domain.User{
				ExternalID:    update.ExternalUserID,
				Language:      language,
				SelectedModel: "gpt-4o-mini",
				CreatedAt:     now,
				UpdatedAt:     now,
			}

			user, err = s.storage.CreateUser(ctx, newUser)
			if err != nil {
				return fmt.Errorf("can't create user: %w", err)
			}

			s.logger.InfoContext(ctx, "created new user",
				slog.String("external_id", user.ExternalID))
		} else {
			return fmt.Errorf("can't get user: %w", err)
		}
	}

	// Save incoming message from user
	_, err = s.storage.CreateMessage(ctx, &domain.Message{
		UserID: user.ID,
		MessageType: domain.MessageType{
			Text: update.MessageText,
		},
		SentBy:    domain.MessageSenderUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("can't save user message: %w", err)
	}

	// Get conversation history
	messages, err := s.storage.GetMessagesByUserID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("can't get messages: %w", err)
	}

	// Send typing indicator
	err = s.sender.SendTyping(ctx, user.ExternalID)
	if err != nil {
		s.logger.WarnContext(ctx, "failed to send typing indicator", slog.String("error", err.Error()))
	}

	// Get AI completion with streaming
	tokenStream, err := s.completion.CompleteStream(ctx, user.SelectedModel, messages)
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
	lastUpdate := time.Now()

	for token := range tokenStream {
		if token.Error != nil {
			return fmt.Errorf("completion stream error: %w", token.Error)
		}
		responseBuilder.WriteString(token.Content)

		// Update message at most once per second
		if time.Since(lastUpdate) >= time.Second/2 {
			// Send typing indicator to keep the conversation active
			err = s.sender.SendTyping(ctx, user.ExternalID)
			if err != nil {
				s.logger.WarnContext(ctx, "failed to send typing indicator", slog.String("error", err.Error()))
			}

			// Update all tracked messages
			updatedMessageIDs, err := s.sender.UpdateMessages(ctx, user.ExternalID, messageIDs, responseBuilder.String())
			if err != nil {
				return fmt.Errorf("can't update messages: %w", err)
			}
			// Update our tracked message IDs
			messageIDs = updatedMessageIDs
			lastUpdate = time.Now()
		}
	}

	// Send final update if needed
	if time.Since(lastUpdate) > 0 {
		updatedMessageIDs, err := s.sender.UpdateMessages(ctx, user.ExternalID, messageIDs, responseBuilder.String())
		if err != nil {
			return fmt.Errorf("can't update final messages: %w", err)
		}
		// Update our tracked message IDs for final state
		messageIDs = updatedMessageIDs
	}

	responseText := responseBuilder.String()

	_, err = s.storage.CreateMessage(ctx, &domain.Message{
		UserID: user.ID,
		MessageType: domain.MessageType{
			Text: responseText,
		},
		SentBy:    domain.MessageSenderBot,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("can't save bot message: %w", err)
	}

	return nil
}
