package openai

import (
	"context"
	"errors"
	"fmt"
	"io"

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
	messages []*domain.Message,
) (<-chan completion.StreamToken, error) {
	// Convert domain messages to OpenAI messages
	chatMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))

	for _, msg := range messages {
		switch msg.SentBy {
		case domain.MessageSenderUser:
			chatMessages = append(chatMessages, openai.UserMessage(msg.MessageType.Text))
		case domain.MessageSenderBot:
			chatMessages = append(chatMessages, openai.AssistantMessage(msg.MessageType.Text))
		}
	}

	// Create streaming request
	stream := o.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    model,
		Messages: chatMessages,
	})

	// Create channel for streaming tokens
	tokenChan := make(chan completion.StreamToken)

	// Start goroutine to read from stream
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

		// Check for stream error
		if err := stream.Err(); err != nil {
			select {
			case tokenChan <- completion.StreamToken{Error: fmt.Errorf("stream error: %w", err)}:
			case <-ctx.Done():
			}
		}
	}()

	return tokenChan, nil
}

// CompleteStreamWithSystemPrompt adds a system prompt before the messages.
func (o *Completion) CompleteStreamWithSystemPrompt(
	ctx context.Context,
	model string,
	systemPrompt string,
	messages []*domain.Message,
) (<-chan completion.StreamToken, error) {
	// Convert domain messages to OpenAI messages
	chatMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages)+1)

	// Add system prompt
	if systemPrompt != "" {
		chatMessages = append(chatMessages, openai.SystemMessage(systemPrompt))
	}

	// Add conversation messages
	for _, msg := range messages {
		switch msg.SentBy {
		case domain.MessageSenderUser:
			chatMessages = append(chatMessages, openai.UserMessage(msg.MessageType.Text))
		case domain.MessageSenderBot:
			chatMessages = append(chatMessages, openai.AssistantMessage(msg.MessageType.Text))
		}
	}

	// Create streaming request
	stream := o.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    model,
		Messages: chatMessages,
	})

	// Create channel for streaming tokens
	tokenChan := make(chan completion.StreamToken)

	// Start goroutine to read from stream
	go func() {
		defer close(tokenChan)

		// Stream reader for OpenAI's streaming responses
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

		// Check for any errors after streaming completes
		if err := stream.Err(); err != nil && !errors.Is(err, io.EOF) {
			select {
			case tokenChan <- completion.StreamToken{Error: fmt.Errorf("stream error: %w", err)}:
			case <-ctx.Done():
			}
		}
	}()

	return tokenChan, nil
}
