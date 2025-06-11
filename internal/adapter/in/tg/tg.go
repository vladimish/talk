package tg

import (
	"context"
	"fmt"
	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/service"
	"github.com/vladimish/talk/pkg/slogctx"
	"log/slog"
	"strconv"

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
		ExternalID:     strconv.FormatInt(update.ID, 10),
		ExternalUserID: strconv.FormatInt(update.Message.From.ID, 10),
		UserLanguage:   update.Message.From.LanguageCode,
		MessageText:    update.Message.Text,
	})
	if err != nil {
		b.l.ErrorContext(ctx, fmt.Errorf("error while handling update: %w", err).Error())
	}
}
