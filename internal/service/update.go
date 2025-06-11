package service

import (
	"context"
	"fmt"
	"log/slog"
	"talk/internal/domain"
	"talk/internal/port/sender"
	"talk/internal/port/storage"
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
		// todo: handle no rows
		return fmt.Errorf("can't get user: %w", err)
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
