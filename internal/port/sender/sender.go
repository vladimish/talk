package sender

import (
	"context"
	"talk/internal/domain"
)

type Sender interface {
	SendMessage(ctx context.Context, message domain.Message) (*domain.Message, error)
}
