package storage

import (
	"context"
	"errors"
	"talk/internal/domain"
)

type Storage interface {
	GetUserByExternalUserID(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUserCurrentStep(ctx context.Context, userID int64, currentStep string) error

	CreateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error)
	GetMessagesByUserID(ctx context.Context, userID int64) ([]*domain.Message, error)
}

var ErrNotFound = errors.New("not found")
