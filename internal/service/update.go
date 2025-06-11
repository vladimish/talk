package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"talk/internal/domain"
	"talk/internal/port/sender"
	"talk/internal/port/storage"
	"time"
)

type UpdateService struct {
	logger  *slog.Logger
	storage storage.Storage
	sender  sender.Sender
}

func NewUpdateService(logger *slog.Logger, storage storage.Storage, sender sender.Sender) *UpdateService {
	return &UpdateService{
		logger:  logger,
		storage: storage,
		sender:  sender,
	}
}

func (s *UpdateService) HandleUpdate(ctx context.Context, update domain.Update) error {
	user, err := s.storage.GetUserByExternalUserID(ctx, update.Message.ExternalUserID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			language := update.UserLanguage
			now := time.Now()
			newUser := &domain.User{
				ExternalID: update.Message.ExternalUserID,
				Language:   language,
				CreatedAt:  now,
				UpdatedAt:  now,
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

	_, err = s.sender.SendMessage(ctx, domain.Message{
		ExternalUserID: user.ExternalID,
		Content:        update.Message.Content,
	})
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}
