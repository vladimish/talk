package sender

import (
	"context"

	"github.com/vladimish/talk/internal/domain"
)

//go:generate go tool mockgen -source=sender.go -destination=../../../mocks/mock_sender.go -package=mocks

type Sender interface {
	SendMessage(ctx context.Context, externalUserID string, text string) (string, error)
	SendMessageWithContent(ctx context.Context, externalUserID string, content domain.MessageContent) (string, error)
	UpdateMessage(ctx context.Context, externalUserID string, messageID string, text string) ([]string, error)
	UpdateMessages(
		ctx context.Context,
		externalUserID string,
		messageIDs []string,
		previousText, currentText string,
	) ([]string, error)
	SendTyping(ctx context.Context, externalUserID string) error
	DeleteMessage(ctx context.Context, externalUserID string, messageID string) error
	CreateInvoiceLink(ctx context.Context, params domain.CreateInvoiceLinkParams) (string, error)
	AnswerPreCheckoutQuery(ctx context.Context, preCheckoutQueryID string, ok bool, errorMessage string) error
}
