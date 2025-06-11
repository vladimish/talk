package formatter

import "context"

type Formatter interface {
	FormatMarkdown(ctx context.Context, text string) (string, error)
}
