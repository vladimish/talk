package tg

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/vladimish/talk/internal/port/formatter"
)

type Sender struct {
	bot       *bot.Bot
	formatter formatter.Formatter
	logger    *slog.Logger
}

func NewSender(bot *bot.Bot, formatter formatter.Formatter, logger *slog.Logger) *Sender {
	return &Sender{
		bot:       bot,
		formatter: formatter,
		logger:    logger,
	}
}

func (u *Sender) SendMessage(ctx context.Context, externalUserID string, text string) (string, error) {
	// Format the text for Telegram
	formattedText, err := u.formatter.FormatMarkdown(ctx, text)
	if err != nil {
		u.logger.WarnContext(ctx, "failed to format markdown, using raw text", slog.String("error", err.Error()))
		formattedText = text
	}

	msg, err := u.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    externalUserID,
		Text:      formattedText,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		return "", fmt.Errorf("can't send message: %w", err)
	}

	return strconv.Itoa(msg.ID), nil
}

func (u *Sender) UpdateMessage(
	ctx context.Context,
	externalUserID string,
	messageID string,
	text string,
) ([]string, error) {
	const maxMessageLength = 4096

	// Split text first, then format each chunk
	chunks := splitText(text, maxMessageLength)
	messageIDs := make([]string, 0, len(chunks))

	for i, chunk := range chunks {
		formattedChunk, formatErr := u.formatter.FormatMarkdown(ctx, chunk)
		if formatErr != nil {
			u.logger.WarnContext(
				ctx,
				"failed to format markdown for chunk, using raw text",
				slog.String("error", formatErr.Error()),
			)
			formattedChunk = chunk
		}

		var newMessageID string
		var err error

		if i == 0 {
			newMessageID, err = u.handleFirstMessageUpdate(ctx, externalUserID, messageID, formattedChunk)
		} else {
			newMessageID, err = u.sendOverflowMessage(ctx, externalUserID, formattedChunk)
		}

		if err != nil {
			return nil, err
		}
		messageIDs = append(messageIDs, newMessageID)
	}

	return messageIDs, nil
}

func (u *Sender) updateFirstMessage(ctx context.Context, externalUserID, messageID, text string) error {
	msgID, err := strconv.Atoi(messageID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	_, err = u.bot.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    externalUserID,
		MessageID: msgID,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		return fmt.Errorf("can't edit message: %w", err)
	}

	return nil
}

func (u *Sender) handleFirstMessageUpdate(ctx context.Context, externalUserID, messageID, text string) (string, error) {
	err := u.updateFirstMessage(ctx, externalUserID, messageID, text)
	if err != nil {
		if isMessageNotModified(err) {
			u.logger.DebugContext(ctx, "message content unchanged, skipping update")
			return messageID, nil
		}
		return "", err
	}
	return messageID, nil
}

func (u *Sender) sendOverflowMessage(ctx context.Context, externalUserID, text string) (string, error) {
	msg, err := u.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    externalUserID,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		return "", fmt.Errorf("can't send overflow message: %w", err)
	}

	return strconv.Itoa(msg.ID), nil
}

func (u *Sender) updateExistingMessage(ctx context.Context, externalUserID, messageID, text string) error {
	msgID, err := strconv.Atoi(messageID)
	if err != nil {
		return fmt.Errorf("invalid message ID %s: %w", messageID, err)
	}

	_, err = u.bot.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    externalUserID,
		MessageID: msgID,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		return fmt.Errorf("can't edit message %s: %w", messageID, err)
	}

	return nil
}

func (u *Sender) handleExistingMessageUpdate(
	ctx context.Context,
	externalUserID, messageID, text string,
) (string, error) {
	err := u.updateExistingMessage(ctx, externalUserID, messageID, text)
	if err != nil {
		if !isMessageNotModified(err) {
			return "", err
		}
		u.logger.DebugContext(
			ctx,
			"message content unchanged, skipping update",
			slog.String("message_id", messageID),
		)
	}
	return messageID, nil
}

func isMessageNotModified(err error) bool {
	return strings.Contains(err.Error(), "message is not modified")
}

func (u *Sender) UpdateMessages(
	ctx context.Context,
	externalUserID string,
	messageIDs []string,
	previousText, currentText string,
) ([]string, error) {
	const maxMessageLength = 4096

	// Split both previous and current text into chunks
	previousChunks := splitText(previousText, maxMessageLength)
	currentChunks := splitText(currentText, maxMessageLength)
	resultMessageIDs := make([]string, 0, len(currentChunks))

	for i, chunk := range currentChunks {
		// Format each chunk individually
		formattedChunk, formatErr := u.formatter.FormatMarkdown(ctx, chunk)
		if formatErr != nil {
			u.logger.WarnContext(
				ctx,
				"failed to format markdown for chunk, using raw text",
				slog.String("error", formatErr.Error()),
			)
			formattedChunk = chunk
		}

		var messageID string
		var err error

		if i < len(messageIDs) {
			// Check if we need to update this chunk
			shouldUpdate := true
			if i < len(previousChunks) {
				// Format previous chunk for comparison
				formattedPrevChunk, prevFormatErr := u.formatter.FormatMarkdown(ctx, previousChunks[i])
				if prevFormatErr == nil && formattedPrevChunk == formattedChunk {
					// Chunks are identical, no need to update
					shouldUpdate = false
					messageID = messageIDs[i]
					u.logger.DebugContext(ctx, "chunk unchanged, skipping update",
						slog.String("message_id", messageID),
						slog.Int("chunk_index", i))
				}
			}

			if shouldUpdate {
				// Update existing message
				messageID, err = u.handleExistingMessageUpdate(ctx, externalUserID, messageIDs[i], formattedChunk)
			}
		} else {
			// Send new message for overflow
			messageID, err = u.sendOverflowMessage(ctx, externalUserID, formattedChunk)
		}

		if err != nil {
			return nil, err
		}
		resultMessageIDs = append(resultMessageIDs, messageID)
	}

	// Note: We don't delete leftover messages if we have fewer chunks now
	// They will just remain as-is. This is simpler and avoids potential issues.

	return resultMessageIDs, nil
}

func splitText(text string, maxLength int) []string {
	if len(text) <= maxLength {
		return []string{text}
	}

	return smartSplitWithCodeBlocks(text, maxLength)
}

func smartSplitWithCodeBlocks(text string, maxLength int) []string {
	var chunks []string
	runes := []rune(text)
	inCodeBlock := false

	for len(runes) > 0 {
		// Account for potential code block markers we'll add
		adjustedMaxLength := maxLength
		if inCodeBlock {
			adjustedMaxLength -= 4 // "```\n"
		}

		chunkSize := adjustedMaxLength
		if len(runes) < chunkSize {
			chunkSize = len(runes)
		}

		// Try to break at a word boundary
		if chunkSize < len(runes) && chunkSize > adjustedMaxLength/2 {
			for i := chunkSize - 1; i >= adjustedMaxLength/2; i-- {
				if runes[i] == ' ' || runes[i] == '\n' {
					chunkSize = i + 1
					break
				}
			}
		}

		chunk := string(runes[:chunkSize])
		originalChunk := chunk

		// Count ``` occurrences in the original chunk (before we add any)
		codeBlockCount := strings.Count(originalChunk, "```")

		// Determine if we'll be in a code block after this chunk
		newInCodeBlock := inCodeBlock
		if inCodeBlock {
			// We started in a code block, odd number of ``` closes it
			newInCodeBlock = (codeBlockCount % 2) == 0
		} else {
			// We started outside, odd number of ``` opens it
			newInCodeBlock = (codeBlockCount % 2) == 1
		}

		// If we'll be in a code block and there's more text, we need to handle this carefully
		if newInCodeBlock && chunkSize < len(runes) {
			// Find the start of the unclosed code block in this chunk
			lastCodeBlockStart := -1
			chunkText := originalChunk
			for i := len(chunkText) - 3; i >= 0; i-- {
				if i+2 < len(chunkText) &&
					chunkText[i] == '`' && chunkText[i+1] == '`' && chunkText[i+2] == '`' {
					// Check if this is an opening ``` (odd total count up to this point)
					prefixCount := strings.Count(chunkText[:i], "```")
					totalCount := prefixCount
					if inCodeBlock {
						totalCount++ // We started in a code block
					}
					if totalCount%2 == 0 {
						// This is an opening ```, so it's the start of our unclosed block
						lastCodeBlockStart = i
						break
					}
				}
			}

			if lastCodeBlockStart > 0 { // Only move if there's content before the code block
				// Move the entire code block to the next chunk
				chunk = originalChunk[:lastCodeBlockStart]
				chunkSize = lastCodeBlockStart
				newInCodeBlock = inCodeBlock // Reset since we removed the code block

				// Update counts after removing the code block
				codeBlockCount = strings.Count(chunk, "```")
				if inCodeBlock {
					newInCodeBlock = (codeBlockCount % 2) == 0
				} else {
					newInCodeBlock = (codeBlockCount % 2) == 1
				}
			} else if lastCodeBlockStart == 0 {
				// Look ahead to find the closing ``` in the remaining text
				remainingText := string(runes)
				closingPos := -1
				searchStart := chunkSize + 3 // Start after the opening ```
				for i := searchStart; i <= len(remainingText)-3; i++ {
					if remainingText[i] == '`' && remainingText[i+1] == '`' && remainingText[i+2] == '`' {
						closingPos = i
						break
					}
				}

				if closingPos != -1 && closingPos+3 <= maxLength {
					// The entire code block fits in one chunk, use it all
					chunkSize = closingPos + 3
					chunk = string(runes[:chunkSize])
					newInCodeBlock = false // Code block is now closed
				}
				// If code block is too large, we continue it in the next chunk without artificial closing
			}
			// If we can't find the code block start, we continue it in the next chunk without artificial closing
		}

		// If we're continuing a code block from the previous chunk, prepend ```
		if inCodeBlock {
			chunk = "```\n" + chunk
		}

		chunks = append(chunks, chunk)
		runes = runes[chunkSize:]
		inCodeBlock = newInCodeBlock
	}

	return chunks
}

func (u *Sender) SendTyping(ctx context.Context, externalUserID string) error {
	_, err := u.bot.SendChatAction(ctx, &bot.SendChatActionParams{
		ChatID: externalUserID,
		Action: models.ChatActionTyping,
	})
	return err
}
