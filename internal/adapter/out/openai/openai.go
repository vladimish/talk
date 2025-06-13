package openai

import (
	"context"
	"fmt"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/completion"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Completion struct {
	client *openai.Client
}

func NewOpenAICompletion(apiKey string) *Completion {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://openrouter.ai/api/v1"),
	)

	return &Completion{
		client: &client,
	}
}

func (o *Completion) CompleteStream(
	ctx context.Context,
	model string,
	systemPrompt string,
	messages []*domain.Message,
) (<-chan completion.StreamToken, error) {
	chatMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))

	if systemPrompt != "" {
		chatMessages = append(chatMessages, openai.SystemMessage(systemPrompt))
	}

	for _, msg := range messages {
		switch msg.SentBy {
		case domain.MessageSenderUser:
			chatMessages = append(chatMessages, openai.UserMessage(msg.MessageType.Text))
		case domain.MessageSenderBot:
			chatMessages = append(chatMessages, openai.AssistantMessage(msg.MessageType.Text))
		}
	}

	stream := o.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    model,
		Messages: chatMessages,
	})

	tokenChan := make(chan completion.StreamToken)

	go func() {
		defer close(tokenChan)

		for stream.Next() {
			chunk := stream.Current()
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				select {
				case tokenChan <- completion.StreamToken{Content: chunk.Choices[0].Delta.Content}:
				case <-ctx.Done():
					return
				}
			}
		}

		if err := stream.Err(); err != nil {
			select {
			case tokenChan <- completion.StreamToken{Error: fmt.Errorf("stream error: %w", err)}:
			case <-ctx.Done():
			}
		}
	}()

	return tokenChan, nil
}
