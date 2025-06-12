package sender

import (
	"context"
)

type Sender interface {
	SendMessage(ctx context.Context, externalUserID string, text string) (string, error)
	UpdateMessage(ctx context.Context, externalUserID string, messageID string, text string) ([]string, error)
	UpdateMessages(ctx context.Context, externalUserID string, messageIDs []string, previousText, currentText string) ([]string, error)
	SendTyping(ctx context.Context, externalUserID string) error
}
