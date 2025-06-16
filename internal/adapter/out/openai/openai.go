package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/completion"
)

const (
	httpErrorThreshold = 400
	maxFileNameLength  = 100
)

// sanitizeFilename removes problematic characters from filenames for OpenRouter.
func sanitizeFilename(filename string) string {
	// Replace non-ASCII characters, spaces, and special characters with underscores
	reg := regexp.MustCompile(`[^\w\-.]`)
	sanitized := reg.ReplaceAllString(filename, "_")

	// Limit filename length to maxFileNameLength characters
	if len(sanitized) > maxFileNameLength {
		ext := ""
		if idx := strings.LastIndex(sanitized, "."); idx > 0 && idx < len(sanitized)-1 {
			ext = sanitized[idx:]
			sanitized = sanitized[:idx]
		}
		sanitized = sanitized[:maxFileNameLength-len(ext)] + ext
	}

	// Ensure it has a .pdf extension
	if !strings.HasSuffix(strings.ToLower(sanitized), ".pdf") {
		sanitized += ".pdf"
	}

	return sanitized
}

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

func (o *Completion) CompleteStreamWithAttachments(
	ctx context.Context,
	model string,
	systemPrompt string,
	messages []*domain.Message,
	imageAttachment *completion.FileAttachment,
	pdfAttachment *completion.FileAttachment,
	webSearchEnabled bool,
) (<-chan completion.StreamToken, error) {
	// For backwards compatibility, if there's no PDF attachment, use the old method
	if pdfAttachment == nil {
		imageURL := ""
		if imageAttachment != nil {
			imageURL = imageAttachment.URL
		}
		return o.CompleteStream(ctx, model, systemPrompt, messages, imageURL, webSearchEnabled)
	}

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

	// If there's a PDF, we need to use custom request format for OpenRouter
	if len(openaiMessages) > 0 {
		// Find the last user message and convert it to multipart content
		for i := len(openaiMessages) - 1; i >= 0; i-- {
			if openaiMessages[i].Role == openai.ChatMessageRoleUser {
				// For PDF, we'll need to use the custom request format
				return o.createStreamWithPDF(ctx, openaiMessages, pdfAttachment, webSearchEnabled)
			}
		}
	}

	// If no PDF but has image, use standard multipart format
	if imageAttachment != nil && len(openaiMessages) > 0 {
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
							URL: imageAttachment.URL,
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
		// Log the full error for debugging, especially HTTP errors
		o.logger.ErrorContext(ctx, "Failed to create completion stream",
			"error", err.Error(),
			"model", model,
			"messages_count", len(openaiMessages))
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
		// Read response body for error details
		bodyBytes, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		var errorMsg string
		if readErr != nil {
			errorMsg = fmt.Sprintf("HTTP error: %d (failed to read response body: %v)", resp.StatusCode, readErr)
		} else {
			errorMsg = fmt.Sprintf("HTTP error: %d - Response: %s", resp.StatusCode, string(bodyBytes))
		}

		return nil, fmt.Errorf("%s", errorMsg)
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

// CustomMessageContent represents custom message content for PDF support.
type CustomMessageContent struct {
	Type string                 `json:"type"`
	Text string                 `json:"text,omitempty"`
	File map[string]interface{} `json:"file,omitempty"`
}

// CustomMessage represents a message with custom content format.
type CustomMessage struct {
	Role    string                 `json:"role"`
	Content []CustomMessageContent `json:"content"`
}

// CustomRequestWithPDF represents a request with PDF file support.
type CustomRequestWithPDF struct {
	Model    string                   `json:"model"`
	Messages []CustomMessage          `json:"messages"`
	Stream   bool                     `json:"stream"`
	Plugins  []map[string]interface{} `json:"plugins,omitempty"`
}

func (o *Completion) createStreamWithPDF(
	ctx context.Context,
	messages []openai.ChatCompletionMessage,
	pdfAttachment *completion.FileAttachment,
	webSearchEnabled bool,
) (<-chan completion.StreamToken, error) {
	// Convert OpenAI messages to custom format with PDF support
	customMessages := make([]CustomMessage, 0, len(messages))

	for i, msg := range messages {
		// For the last user message, add PDF attachment
		if i == len(messages)-1 && msg.Role == openai.ChatMessageRoleUser && pdfAttachment != nil {
			// Convert PDF data to base64 data URL
			base64Data := base64.StdEncoding.EncodeToString(pdfAttachment.Data)
			dataURL := fmt.Sprintf("data:application/pdf;base64,%s", base64Data)

			// Sanitize filename for OpenRouter
			sanitizedFilename := sanitizeFilename(pdfAttachment.FileName)

			customMessages = append(customMessages, CustomMessage{
				Role: msg.Role,
				Content: []CustomMessageContent{
					{
						Type: "text",
						Text: msg.Content,
					},
					{
						Type: "file",
						File: map[string]interface{}{
							"filename":  sanitizedFilename,
							"file_data": dataURL,
						},
					},
				},
			})
		} else {
			// Regular message
			customMessages = append(customMessages, CustomMessage{
				Role: msg.Role,
				Content: []CustomMessageContent{
					{
						Type: "text",
						Text: msg.Content,
					},
				},
			})
		}
	}

	// Create request body
	reqBody := CustomRequestWithPDF{
		Model:    "openrouter/auto", // Use auto-routing for PDF support
		Messages: customMessages,
		Stream:   true,
	}

	// Add plugins if needed
	if webSearchEnabled {
		reqBody.Plugins = []map[string]interface{}{
			{"id": "web"},
		}
	}

	// Add file parser plugin for PDF
	reqBody.Plugins = append(reqBody.Plugins, map[string]interface{}{
		"id": "file-parser",
		"pdf": map[string]string{
			"engine": "native", // Use text extraction engine
		},
	})

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal PDF request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://openrouter.ai/api/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req) //nolint:bodyclose // body is closed in the goroutine defer
	if err != nil {
		return nil, fmt.Errorf("failed to send PDF request: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= httpErrorThreshold {
		// Read response body for error details
		bodyBytes, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		var errorMsg string
		if readErr != nil {
			errorMsg = fmt.Sprintf("HTTP error: %d (failed to read response body: %v)", resp.StatusCode, readErr)
		} else {
			errorMsg = fmt.Sprintf("HTTP error: %d - Response: %s", resp.StatusCode, string(bodyBytes))
		}

		return nil, fmt.Errorf("%s", errorMsg)
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
