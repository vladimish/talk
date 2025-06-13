package storage

import (
	"context"
	"errors"

	"github.com/vladimish/talk/internal/domain"
)

type Storage interface {
	GetUserByExternalUserID(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUserCurrentStep(ctx context.Context, userID int64, currentStep string) error
	UpdateUserSelectedModel(ctx context.Context, userID int64, selectedModel string) error

	CreateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error)
	GetMessagesByUserID(ctx context.Context, userID int64) ([]*domain.Message, error)
	GetMessagesByConversation(ctx context.Context, conversationID string) ([]*domain.Message, error)
	GetConversationsByUserID(ctx context.Context, userID int64) ([]string, error)

	UpdateUserCurrentConversation(ctx context.Context, userID int64, conversationID *string) error
	CreateForeignMessage(ctx context.Context, messageID int32, foreignMessageID int32) error
}

var ErrNotFound = errors.New("not found")
