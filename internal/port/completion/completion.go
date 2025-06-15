package completion

import (
	"context"

	"github.com/vladimish/talk/internal/domain"
)

//go:generate go tool mockgen -source=completion.go -destination=../../../mocks/mock_completion.go -package=mocks

type StreamToken struct {
	Content string
	Error   error
}

type Completion interface {
	CompleteStream(
		ctx context.Context,
		model string,
		systemPrompt string,
		messages []*domain.Message,
		currentImageURL string, // Empty string if no image
	) (<-chan StreamToken, error)
}
