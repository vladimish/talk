package formatter

import "context"

//go:generate go tool mockgen -source=formatter.go -destination=../../../mocks/mock_formatter.go -package=mocks

type Formatter interface {
	FormatMarkdown(ctx context.Context, text string) (string, error)
}
