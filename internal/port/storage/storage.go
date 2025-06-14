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
	UpdateUserLanguage(ctx context.Context, userID int64, language string) error

	CreateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error)
	GetMessagesByUserID(ctx context.Context, userID int64) ([]*domain.Message, error)
	GetMessagesByConversationID(ctx context.Context, conversationID int64) ([]*domain.Message, error)
	GetLatestMessageByConversationID(ctx context.Context, conversationID int64) (*domain.Message, error)

	CreateConversation(ctx context.Context, conversation *domain.Conversation) (*domain.Conversation, error)
	GetConversationsByUserID(ctx context.Context, userID int64) ([]*domain.Conversation, error)
	GetConversationByID(ctx context.Context, conversationID int64) (*domain.Conversation, error)
	UpdateConversationName(ctx context.Context, conversationID int64, name string) error

	UpdateUserCurrentConversationID(ctx context.Context, userID int64, conversationID *int64) error
	CreateForeignMessage(ctx context.Context, messageID int32, foreignMessageID int32) error
	GetForeignMessageByMessageID(ctx context.Context, messageID int32) (int32, error)
}

var ErrNotFound = errors.New("not found")
