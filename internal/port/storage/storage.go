package storage

import (
	"context"
	"errors"

	"github.com/vladimish/talk/internal/domain"
)

//go:generate go tool mockgen -source=storage.go -destination=../../../mocks/mock_storage.go -package=mocks

type Storage interface {
	GetUserByExternalUserID(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUserCurrentStep(ctx context.Context, userID int64, currentStep string) error
	UpdateUserSelectedModel(ctx context.Context, userID int64, selectedModel string) error
	UpdateUserLanguage(ctx context.Context, userID int64, language string) error
	UpdateUserConversationListOffset(ctx context.Context, userID int64, offset int) error

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

	// Transaction methods
	CreateTransaction(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	GetUserTokenBalance(ctx context.Context, userID int64) (*domain.TokenBalance, error)
	GetUserTokenBalanceByType(ctx context.Context, userID int64, tokenType domain.TokenType) (int64, error)

	// Payment methods
	CreatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	GetPaymentByInvoicePayload(ctx context.Context, invoicePayload string) (*domain.Payment, error)
	UpdatePaymentStatus(
		ctx context.Context,
		paymentID int64,
		status domain.PaymentStatus,
		telegramChargeID, providerChargeID *string,
	) (*domain.Payment, error)
	UpdatePaymentWithInvoice(
		ctx context.Context,
		paymentID int64,
		invoiceLink, invoicePayload, messageID string,
	) (*domain.Payment, error)

	// Subscription methods
	CreateSubscription(ctx context.Context, subscription *domain.Subscription) (*domain.Subscription, error)
	GetActiveSubscriptionByUserID(ctx context.Context, userID int64) (*domain.Subscription, error)

	// Attachment methods
	CreateAttachment(ctx context.Context, attachment *domain.Attachment) (*domain.Attachment, error)
}

var ErrNotFound = errors.New("not found")
