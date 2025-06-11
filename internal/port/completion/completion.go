package completion

import (
	"context"
	"github.com/vladimish/talk/internal/domain"
)

type StreamToken struct {
	Content string
	Error   error
}

type Completion interface {
	CompleteStream(ctx context.Context, model string, messages []*domain.Message) (<-chan StreamToken, error)
}
