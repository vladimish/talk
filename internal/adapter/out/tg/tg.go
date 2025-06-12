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
	text string,
) ([]string, error) {
	const maxMessageLength = 4096

	// Split text first, then format each chunk
	chunks := splitText(text, maxMessageLength)
	resultMessageIDs := make([]string, 0, len(chunks))

	for i, chunk := range chunks {
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
			// Update existing message
			messageID, err = u.handleExistingMessageUpdate(ctx, externalUserID, messageIDs[i], formattedChunk)
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

	for len(runes) > 0 {
		chunkSize := maxLength
		if len(runes) < chunkSize {
			chunkSize = len(runes)
		}

		// Try to break at a word boundary
		if chunkSize < len(runes) {
			for i := chunkSize - 1; i >= maxLength/2; i-- {
				if runes[i] == ' ' || runes[i] == '\n' {
					chunkSize = i + 1
					break
				}
			}
		}

		chunk := string(runes[:chunkSize])

		// Check if chunk ends with an unclosed code block
		if chunkSize < len(runes) { // Only if there's more text
			adjustedChunk, newChunkSize := handleCodeBlockInChunk(chunk, runes, chunkSize, maxLength)
			if newChunkSize != chunkSize {
				chunk = adjustedChunk
				chunkSize = newChunkSize
			}
		}

		chunks = append(chunks, chunk)
		runes = runes[chunkSize:]
	}

	return chunks
}

func handleCodeBlockInChunk(chunk string, allRunes []rune, chunkSize int, maxLength int) (string, int) {
	codeBlockCount := strings.Count(chunk, "```")

	// If even number (or zero), chunk is fine
	if codeBlockCount%2 == 0 {
		return chunk, chunkSize
	}

	// We have an unclosed code block - try to move it to next chunk
	chunkRunes := []rune(chunk)

	// Find the last ``` in the chunk
	const codeBlockMarkerLen = 3
	lastCodeBlockStart := -1
	for i := len(chunkRunes) - codeBlockMarkerLen; i >= 0; i-- {
		if i+2 < len(chunkRunes) &&
			chunkRunes[i] == '`' && chunkRunes[i+1] == '`' && chunkRunes[i+2] == '`' {
			lastCodeBlockStart = i
			break
		}
	}

	if lastCodeBlockStart == -1 {
		// Couldn't find the code block start, use fallback
		return addCodeBlockBoundaries(chunk, true), chunkSize
	}

	// Check if moving the code block to next chunk would make it fit
	beforeCodeBlock := string(chunkRunes[:lastCodeBlockStart])

	// Look ahead to see if the code block will be complete and fit in maxLength
	remainingText := string(allRunes[chunkSize:])
	codeBlockPart := string(chunkRunes[lastCodeBlockStart:])

	// Try to find where this code block ends in the remaining text
	completeCodeBlock := findCompleteCodeBlock(codeBlockPart + remainingText)

	// If the complete code block fits in maxLength, move it to next chunk
	if len([]rune(completeCodeBlock)) <= maxLength {
		// Move the code block to next chunk
		return beforeCodeBlock, lastCodeBlockStart
	}

	// Code block is too large, use fallback splitting
	return addCodeBlockBoundaries(chunk, true), chunkSize
}

func findCompleteCodeBlock(text string) string {
	runes := []rune(text)
	if len(runes) < 3 || runes[0] != '`' || runes[1] != '`' || runes[2] != '`' {
		return text
	}

	// Find the closing ```
	for i := 3; i < len(runes)-2; i++ {
		if runes[i] == '`' && runes[i+1] == '`' && runes[i+2] == '`' {
			// Found closing, include it
			return string(runes[:i+3])
		}
	}

	// No closing found, return the whole text (unclosed code block)
	return text
}

func addCodeBlockBoundaries(chunk string, isUnclosed bool) string {
	if isUnclosed {
		return chunk + "\n```"
	}
	return "```\n" + chunk
}

func (u *Sender) SendTyping(ctx context.Context, externalUserID string) error {
	_, err := u.bot.SendChatAction(ctx, &bot.SendChatActionParams{
		ChatID: externalUserID,
		Action: models.ChatActionTyping,
	})
	return err
}
