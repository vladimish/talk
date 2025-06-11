package tg

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"talk/internal/domain"
	"talk/internal/service"
	"talk/pkg/slogctx"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot struct {
	l *slog.Logger
	s *service.UpdateService
}

func NewBot(l *slog.Logger, s *service.UpdateService) *Bot {
	return &Bot{
		l: l,
		s: s,
	}
}

func (b *Bot) Handle(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil {
		return
	}

	ctx = slogctx.WithField(ctx, "update_id", update.ID)
	b.l.DebugContext(ctx, "handling update")

	if update.Message == nil || update.Message.From == nil {
		b.l.WarnContext(ctx, "unsupported update type")
		return
	}

	err := b.s.HandleUpdate(ctx, domain.Update{
		ExternalID: strconv.FormatInt(update.ID, 10),
		Message: domain.Message{
			ExternalID:     strconv.Itoa(update.Message.ID),
			ExternalUserID: strconv.FormatInt(update.Message.From.ID, 10),
			Content: domain.MessageContent{
				Text: update.Message.Text,
			},
		},
	})
	if err != nil {
		b.l.ErrorContext(ctx, fmt.Errorf("error while handling update: %w", err).Error())
	}
}
