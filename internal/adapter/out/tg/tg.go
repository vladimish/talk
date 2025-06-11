package tg

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Sender struct {
	bot *bot.Bot
}

func NewSender(bot *bot.Bot) *Sender {
	return &Sender{bot: bot}
}

func (u *Sender) SendMessage(ctx context.Context, externalUserID string, text string) (string, error) {
	_, err := u.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    externalUserID,
		Text:      text,
		ParseMode: models.ParseModeMarkdownV1,
	})
	if err != nil {
		return "", fmt.Errorf("can't send message: %w", err)
	}

	return text, nil
}
