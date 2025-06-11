package sender

import (
	"context"
)

type Sender interface {
	SendMessage(ctx context.Context, externalUserID string, text string) (string, error)
}
