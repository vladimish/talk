package tg

import (
	"context"
	"fmt"
	"strconv"
	"talk/internal/domain"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Sender struct {
	bot *bot.Bot
}

func NewSender(bot *bot.Bot) *Sender {
	return &Sender{bot: bot}
}

func (u *Sender) SendMessage(ctx context.Context, message domain.Message) (*domain.Message, error) {
	m, err := u.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    message.ExternalUserID,
		Text:      message.Content.Text,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		return nil, fmt.Errorf("can't send message: %w", err)
	}

	message.ExternalID = strconv.Itoa(m.ID)
	return &message, nil
}
