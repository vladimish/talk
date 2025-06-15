package openai

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sashabaranov/go-openai"
	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/completion"
)

type Completion struct {
	client *openai.Client
	logger *slog.Logger
}

func NewOpenAICompletion(apiKey string) *Completion {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1"

	return &Completion{
		client: openai.NewClientWithConfig(config),
		logger: slog.Default(),
	}
}

func (o *Completion) CompleteStream(
	ctx context.Context,
	model string,
	systemPrompt string,
	messages []*domain.Message,
	currentImageURL string,
) (<-chan completion.StreamToken, error) {
	// Convert domain messages to OpenAI format
	openaiMessages := make([]openai.ChatCompletionMessage, 0, len(messages)+1)

	// Add system prompt if provided
	if systemPrompt != "" {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		})
	}

	// Convert domain messages
	for _, msg := range messages {
		role := openai.ChatMessageRoleUser
		if msg.SentBy == domain.MessageSenderBot {
			role = openai.ChatMessageRoleAssistant
		}

		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.MessageType.Text,
		})
	}

	// If there's an image URL, modify the last user message to include it
	if currentImageURL != "" && len(openaiMessages) > 0 {
		// Find the last user message and convert it to multipart content
		for i := len(openaiMessages) - 1; i >= 0; i-- {
			if openaiMessages[i].Role == openai.ChatMessageRoleUser {
				// Convert to multipart content
				openaiMessages[i].MultiContent = []openai.ChatMessagePart{
					{
						Type: openai.ChatMessagePartTypeText,
						Text: openaiMessages[i].Content,
					},
					{
						Type: openai.ChatMessagePartTypeImageURL,
						ImageURL: &openai.ChatMessageImageURL{
							URL: currentImageURL,
						},
					},
				}
				openaiMessages[i].Content = "" // Clear content when using MultiContent
				break
			}
		}
	}

	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: openaiMessages,
		Stream:   true,
	}

	stream, err := o.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create completion stream: %w", err)
	}

	tokenChan := make(chan completion.StreamToken)

	go func() {
		defer close(tokenChan)
		defer stream.Close()

		for {
			response, recvErr := stream.Recv()
			if recvErr != nil {
				if recvErr.Error() == "EOF" {
					// Stream ended normally
					return
				}
				select {
				case tokenChan <- completion.StreamToken{Error: fmt.Errorf("stream error: %w", recvErr)}:
				case <-ctx.Done():
				}
				return
			}

			// Send content if available
			if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
				select {
				case tokenChan <- completion.StreamToken{Content: response.Choices[0].Delta.Content}:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return tokenChan, nil
}
