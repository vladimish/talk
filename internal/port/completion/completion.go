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

type FileAttachment struct {
	URL      string
	MimeType string
	FileName string
	Data     []byte // Raw file data for direct processing
}

type Completion interface {
	CompleteStream(
		ctx context.Context,
		model string,
		systemPrompt string,
		messages []*domain.Message,
		currentImageURL string, // Empty string if no image
		webSearchEnabled bool, // Whether to enable web search plugin
	) (<-chan StreamToken, error)

	CompleteStreamWithAttachments(
		ctx context.Context,
		model string,
		systemPrompt string,
		messages []*domain.Message,
		imageAttachment *FileAttachment, // nil if no image
		pdfAttachment *FileAttachment, // nil if no PDF
		webSearchEnabled bool, // Whether to enable web search plugin
	) (<-chan StreamToken, error)
}
