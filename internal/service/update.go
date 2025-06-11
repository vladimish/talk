package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"talk/internal/domain"
	"talk/internal/port/completion"
	"talk/internal/port/sender"
	"talk/internal/port/storage"
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

	// Get AI completion with streaming
	tokenStream, err := s.completion.CompleteStream(ctx, user.SelectedModel, messages)
	if err != nil {
		return fmt.Errorf("can't get completion: %w", err)
	}

	// Collect response tokens
	var responseBuilder strings.Builder
	for token := range tokenStream {
		if token.Error != nil {
			return fmt.Errorf("completion stream error: %w", token.Error)
		}
		responseBuilder.WriteString(token.Content)
		// TODO: In a real implementation, you might want to send partial responses
	}

	responseText := responseBuilder.String()

	_, err = s.sender.SendMessage(ctx, user.ExternalID, responseText)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

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
