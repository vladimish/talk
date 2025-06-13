package tg

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/vladimish/talk/internal/domain"
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

func (u *Sender) SendMessageWithContent(
	ctx context.Context,
	externalUserID string,
	content domain.MessageContent,
) (string, error) {
	// Format the text for Telegram
	formattedText, err := u.formatter.FormatMarkdown(ctx, content.Text)
	if err != nil {
		u.logger.WarnContext(ctx, "failed to format markdown, using raw text", slog.String("error", err.Error()))
		formattedText = content.Text
	}

	params := &bot.SendMessageParams{
		ChatID:    externalUserID,
		Text:      formattedText,
		ParseMode: models.ParseModeMarkdown,
	}

	// Add reply keyboard if present
	if content.ReplyKeyboard != nil {
		keyboard := u.buildReplyKeyboard(content.ReplyKeyboard)
		params.ReplyMarkup = keyboard
	}

	msg, err := u.bot.SendMessage(ctx, params)
	if err != nil {
		return "", fmt.Errorf("can't send message: %w", err)
	}

	return strconv.Itoa(msg.ID), nil
}

func (u *Sender) buildReplyKeyboard(keyboard *domain.ReplyKeyboard) *models.ReplyKeyboardMarkup {
	var rows [][]models.KeyboardButton

	for _, buttonRow := range keyboard.Buttons {
		var tgRow []models.KeyboardButton
		for _, button := range buttonRow {
			tgButton := models.KeyboardButton{
				Text: button.Text,
			}
			tgRow = append(tgRow, tgButton)
		}
		rows = append(rows, tgRow)
	}

	return &models.ReplyKeyboardMarkup{
		Keyboard:              rows,
		ResizeKeyboard:        keyboard.Resize,
		OneTimeKeyboard:       keyboard.OneTime,
		InputFieldPlaceholder: keyboard.Placeholder,
	}
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

func (u *Sender) handleMessageUpdate(
	ctx context.Context,
	externalUserID string,
	messageIDs []string,
	previousChunks []string,
	formattedChunk string,
	chunkIndex int,
) (string, error) {
	if chunkIndex < len(messageIDs) {
		shouldUpdate := u.shouldUpdateChunk(ctx, previousChunks, formattedChunk, chunkIndex)
		if !shouldUpdate {
			return messageIDs[chunkIndex], nil
		}
		return u.handleExistingMessageUpdate(ctx, externalUserID, messageIDs[chunkIndex], formattedChunk)
	}
	return u.sendOverflowMessage(ctx, externalUserID, formattedChunk)
}

func (u *Sender) shouldUpdateChunk(
	ctx context.Context,
	previousChunks []string,
	formattedChunk string,
	chunkIndex int,
) bool {
	if chunkIndex >= len(previousChunks) {
		return true
	}

	formattedPrevChunk, prevFormatErr := u.formatter.FormatMarkdown(ctx, previousChunks[chunkIndex])
	if prevFormatErr != nil || formattedPrevChunk != formattedChunk {
		return true
	}

	u.logger.DebugContext(ctx, "chunk unchanged, skipping update",
		slog.String("message_id", ""), // Will be filled by caller
		slog.Int("chunk_index", chunkIndex))
	return false
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

		messageID, err = u.handleMessageUpdate(ctx, externalUserID, messageIDs, previousChunks, formattedChunk, i)

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

func calculateAdjustedMaxLength(maxLength int, inCodeBlock bool) int {
	adjustedMaxLength := maxLength
	if inCodeBlock {
		adjustedMaxLength -= 4 // "```\n"
	}
	return adjustedMaxLength
}

func calculateInitialChunkSize(runes []rune, adjustedMaxLength int) int {
	chunkSize := adjustedMaxLength
	if len(runes) < chunkSize {
		chunkSize = len(runes)
	}
	return chunkSize
}

func findWordBoundary(runes []rune, chunkSize, adjustedMaxLength int) int {
	if chunkSize < len(runes) && chunkSize > adjustedMaxLength/2 {
		for i := chunkSize - 1; i >= adjustedMaxLength/2; i-- {
			if runes[i] == ' ' || runes[i] == '\n' {
				return i + 1
			}
		}
	}
	return chunkSize
}

const (
	codeBlockMarkerLength = 3
	codeBlockDivisor      = 2
)

func calculateNewCodeBlockState(inCodeBlock bool, codeBlockCount int) bool {
	if inCodeBlock {
		return (codeBlockCount % codeBlockDivisor) == 0
	}
	return (codeBlockCount % codeBlockDivisor) == 1
}

func findLastCodeBlockStart(chunkText string, inCodeBlock bool) int {
	for i := len(chunkText) - codeBlockMarkerLength; i >= 0; i-- {
		if i+2 < len(chunkText) &&
			chunkText[i] == '`' && chunkText[i+1] == '`' && chunkText[i+2] == '`' {
			prefixCount := strings.Count(chunkText[:i], "```")
			totalCount := prefixCount
			if inCodeBlock {
				totalCount++
			}
			if totalCount%2 == 0 {
				return i
			}
		}
	}
	return -1
}

func findCodeBlockClosing(runes []rune, chunkSize int) int {
	remainingText := string(runes)
	searchStart := chunkSize + codeBlockMarkerLength
	for i := searchStart; i <= len(remainingText)-codeBlockMarkerLength; i++ {
		if remainingText[i] == '`' && remainingText[i+1] == '`' && remainingText[i+2] == '`' {
			return i
		}
	}
	return -1
}

type codeBlockResult struct {
	chunk          string
	chunkSize      int
	newInCodeBlock bool
}

func handleCodeBlockBoundary(
	runes []rune,
	chunk string,
	chunkSize int,
	inCodeBlock, newInCodeBlock bool,
	maxLength int,
) codeBlockResult {
	if !newInCodeBlock || chunkSize >= len(runes) {
		return codeBlockResult{chunk, chunkSize, newInCodeBlock}
	}

	originalChunk := chunk
	lastCodeBlockStart := findLastCodeBlockStart(originalChunk, inCodeBlock)

	if lastCodeBlockStart > 0 {
		adjustedChunk := originalChunk[:lastCodeBlockStart]
		adjustedSize := lastCodeBlockStart
		codeBlockCount := strings.Count(adjustedChunk, "```")
		adjustedNewInCodeBlock := calculateNewCodeBlockState(inCodeBlock, codeBlockCount)
		return codeBlockResult{adjustedChunk, adjustedSize, adjustedNewInCodeBlock}
	}

	if lastCodeBlockStart == 0 {
		closingPos := findCodeBlockClosing(runes, chunkSize)
		if closingPos != -1 && closingPos+3 <= maxLength {
			fullChunk := string(runes[:closingPos+3])
			return codeBlockResult{fullChunk, closingPos + 3, false}
		}
	}

	return codeBlockResult{chunk, chunkSize, newInCodeBlock}
}

func findNextEmptyLinePosition(remainingText string) int {
	currentPos := 0
	lines := strings.Split(remainingText, "\n")
	for i, line := range lines {
		if i > 0 && strings.TrimSpace(line) == "" {
			return currentPos
		}
		currentPos += len(line) + 1
	}
	return -1
}

func findDoubleNewlineBreakPoint(chunk string, chunkSize, remainingInChunk int) (int, bool) {
	for i := chunkSize - 1; i >= chunkSize-remainingInChunk; i-- {
		if i < len(chunk) && chunk[i] == '\n' && i > 0 && chunk[i-1] == '\n' {
			return i, true
		}
	}
	return chunkSize, false
}

func findSingleNewlineBreakPoint(runes []rune, chunk string, chunkSize, remainingInChunk int) int {
	const singleNewlineThreshold = 100
	if remainingInChunk >= singleNewlineThreshold {
		return chunkSize
	}

	remainingText := string(runes[chunkSize:])
	nextNewlinePos := -1
	for i, char := range remainingText {
		if char == '\n' {
			nextNewlinePos = i
			break
		}
	}

	if nextNewlinePos > 0 && nextNewlinePos <= singleNewlineThreshold {
		for i := chunkSize - 1; i >= chunkSize-remainingInChunk; i-- {
			if i < len(chunk) && chunk[i] == '\n' {
				return i + 1
			}
		}
	}
	return chunkSize
}

func optimizeTextBlockBoundary(
	runes []rune,
	chunk string,
	chunkSize, adjustedMaxLength int,
	newInCodeBlock bool,
) (string, int) {
	if newInCodeBlock || chunkSize >= len(runes) {
		return chunk, chunkSize
	}

	const threshold = 200
	remainingInChunk := adjustedMaxLength - chunkSize

	if remainingInChunk <= 0 || remainingInChunk >= threshold {
		return chunk, chunkSize
	}

	remainingText := string(runes[chunkSize:])
	nextEmptyLinePos := findNextEmptyLinePosition(remainingText)

	if nextEmptyLinePos <= 0 || nextEmptyLinePos > threshold {
		return chunk, chunkSize
	}

	newChunkSize, found := findDoubleNewlineBreakPoint(chunk, chunkSize, remainingInChunk)
	if found {
		return string(runes[:newChunkSize]), newChunkSize
	}

	newChunkSize = findSingleNewlineBreakPoint(runes, chunk, chunkSize, remainingInChunk)
	return string(runes[:newChunkSize]), newChunkSize
}

func smartSplitWithCodeBlocks(text string, maxLength int) []string {
	var chunks []string
	runes := []rune(text)
	inCodeBlock := false

	for len(runes) > 0 {
		adjustedMaxLength := calculateAdjustedMaxLength(maxLength, inCodeBlock)
		chunkSize := calculateInitialChunkSize(runes, adjustedMaxLength)
		chunkSize = findWordBoundary(runes, chunkSize, adjustedMaxLength)

		chunk := string(runes[:chunkSize])
		originalChunk := chunk

		codeBlockCount := strings.Count(originalChunk, "```")
		newInCodeBlock := calculateNewCodeBlockState(inCodeBlock, codeBlockCount)

		result := handleCodeBlockBoundary(runes, chunk, chunkSize, inCodeBlock, newInCodeBlock, maxLength)
		chunk = result.chunk
		chunkSize = result.chunkSize
		newInCodeBlock = result.newInCodeBlock

		chunk, chunkSize = optimizeTextBlockBoundary(runes, chunk, chunkSize, adjustedMaxLength, newInCodeBlock)

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
