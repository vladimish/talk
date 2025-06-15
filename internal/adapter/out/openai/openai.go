package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/completion"
)

const httpErrorThreshold = 400

type Completion struct {
	client *openai.Client
	logger *slog.Logger
	apiKey string
}

func NewOpenAICompletion(apiKey string) *Completion {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1"

	return &Completion{
		client: openai.NewClientWithConfig(config),
		logger: slog.Default(),
		apiKey: apiKey,
	}
}

func (o *Completion) CompleteStream(
	ctx context.Context,
	model string,
	systemPrompt string,
	messages []*domain.Message,
	currentImageURL string,
	webSearchEnabled bool,
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

	// Handle web search plugin with custom request
	if webSearchEnabled {
		return o.createStreamWithPlugins(ctx, openaiMessages, systemPrompt)
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

// CustomRequest represents a request with plugins support.
type CustomRequest struct {
	Model    string                         `json:"model"`
	Messages []openai.ChatCompletionMessage `json:"messages"`
	Stream   bool                           `json:"stream"`
	Plugins  []map[string]interface{}       `json:"plugins"`
}

// CustomStreamResponse represents the streaming response from OpenRouter.
type CustomStreamResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func (o *Completion) createStreamWithPlugins(
	ctx context.Context,
	messages []openai.ChatCompletionMessage,
	_ string, // systemPrompt is unused in this implementation
) (<-chan completion.StreamToken, error) {
	// Create custom request with plugins
	reqBody := CustomRequest{
		Model:    "openrouter/auto",
		Messages: messages,
		Stream:   true,
		Plugins: []map[string]interface{}{
			{"id": "web"},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://openrouter.ai/api/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req) //nolint:bodyclose // body is closed in the goroutine defer
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= httpErrorThreshold {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	tokenChan := make(chan completion.StreamToken)

	go func() {
		defer close(tokenChan)
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				o.logger.Error("failed to close response body", "error", closeErr)
			}
		}()

		o.processStreamResponse(ctx, resp, tokenChan)
	}()

	return tokenChan, nil
}

func (o *Completion) processStreamResponse(
	ctx context.Context,
	resp *http.Response,
	tokenChan chan<- completion.StreamToken,
) {
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			return
		}

		var streamResp CustomStreamResponse
		if unmarshalErr := json.Unmarshal([]byte(data), &streamResp); unmarshalErr != nil {
			continue // Skip malformed lines
		}

		if len(streamResp.Choices) > 0 {
			content := streamResp.Choices[0].Delta.Content
			if content != "" {
				select {
				case tokenChan <- completion.StreamToken{Content: content}:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		select {
		case tokenChan <- completion.StreamToken{Error: fmt.Errorf("scanner error: %w", scanErr)}:
		case <-ctx.Done():
		}
	}
}
