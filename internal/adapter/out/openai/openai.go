package openai

import (
	"context"
	"fmt"
	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/completion"
	"io"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAICompletion struct {
	client *openai.Client
}

func NewOpenAICompletion(apiKey string) *OpenAICompletion {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://openrouter.ai/api/v1"),
	)

	return &OpenAICompletion{
		client: &client,
	}
}

func (o *OpenAICompletion) CompleteStream(ctx context.Context, model string, messages []*domain.Message) (<-chan completion.StreamToken, error) {
	// Convert domain messages to OpenAI messages
	chatMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))

	for _, msg := range messages {
		if msg.SentBy == domain.MessageSenderUser {
			chatMessages = append(chatMessages, openai.UserMessage(msg.MessageType.Text))
		} else if msg.SentBy == domain.MessageSenderBot {
			chatMessages = append(chatMessages, openai.AssistantMessage(msg.MessageType.Text))
		}
	}

	// Create streaming request
	stream := o.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(model),
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

// CompleteStreamWithSystemPrompt adds a system prompt before the messages
func (o *OpenAICompletion) CompleteStreamWithSystemPrompt(ctx context.Context, model string, systemPrompt string, messages []*domain.Message) (<-chan completion.StreamToken, error) {
	// Convert domain messages to OpenAI messages
	chatMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages)+1)

	// Add system prompt
	if systemPrompt != "" {
		chatMessages = append(chatMessages, openai.SystemMessage(systemPrompt))
	}

	// Add conversation messages
	for _, msg := range messages {
		if msg.SentBy == domain.MessageSenderUser {
			chatMessages = append(chatMessages, openai.UserMessage(msg.MessageType.Text))
		} else if msg.SentBy == domain.MessageSenderBot {
			chatMessages = append(chatMessages, openai.AssistantMessage(msg.MessageType.Text))
		}
	}

	// Create streaming request
	stream := o.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(model),
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
		if err := stream.Err(); err != nil && err != io.EOF {
			select {
			case tokenChan <- completion.StreamToken{Error: fmt.Errorf("stream error: %w", err)}:
			case <-ctx.Done():
			}
		}
	}()

	return tokenChan, nil
}
