package tg

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/service"
	"github.com/vladimish/talk/pkg/slogctx"

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

	// Handle successful payments first
	if update.Message.SuccessfulPayment != nil && update.Message.From != nil {
		b.l.DebugContext(ctx, "user already has active subscription, declining pre-checkout")
		// b.HandleSuccessfulPayment(ctx, nil, update)
		return
	}

	if update.Message == nil || update.Message.From == nil {
		b.l.WarnContext(ctx, "unsupported update type")
		return
	}

	err := b.s.HandleUpdate(ctx, domain.Update{
		ExternalID:        strconv.FormatInt(update.ID, 10),
		ExternalUserID:    strconv.FormatInt(update.Message.From.ID, 10),
		UserLanguage:      update.Message.From.LanguageCode,
		MessageText:       update.Message.Text,
		ExternalMessageID: update.Message.ID,
	})
	if err != nil {
		b.l.ErrorContext(ctx, fmt.Errorf("error while handling update: %w", err).Error())
	}
}

func (b *Bot) HandleCallback(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil || update.CallbackQuery == nil || update.CallbackQuery.From.ID == 0 {
		return
	}

	ctx = slogctx.WithField(ctx, "callback_query_id", update.CallbackQuery.ID)
	b.l.DebugContext(ctx, "handling callback query")

	err := b.s.HandleCallbackQuery(ctx, domain.CallbackQuery{
		ID:             update.CallbackQuery.ID,
		ExternalUserID: strconv.FormatInt(update.CallbackQuery.From.ID, 10),
		UserLanguage:   update.CallbackQuery.From.LanguageCode,
		Data:           update.CallbackQuery.Data,
	})
	if err != nil {
		b.l.ErrorContext(ctx, fmt.Errorf("error while handling callback query: %w", err).Error())
	}
}

func (b *Bot) HandlePreCheckoutQuery(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil || update.PreCheckoutQuery == nil || update.PreCheckoutQuery.From.ID == 0 {
		return
	}

	ctx = slogctx.WithField(ctx, "pre_checkout_query_id", update.PreCheckoutQuery.ID)
	b.l.DebugContext(ctx, "handling pre-checkout query")

	err := b.s.HandlePreCheckoutQuery(ctx, domain.PreCheckoutQuery{
		ID:               update.PreCheckoutQuery.ID,
		ExternalUserID:   strconv.FormatInt(update.PreCheckoutQuery.From.ID, 10),
		Currency:         update.PreCheckoutQuery.Currency,
		TotalAmount:      int64(update.PreCheckoutQuery.TotalAmount),
		InvoicePayload:   update.PreCheckoutQuery.InvoicePayload,
		TelegramChargeID: "", // Not available in pre-checkout
		ProviderChargeID: "", // Not available in pre-checkout
	})
	if err != nil {
		b.l.ErrorContext(ctx, fmt.Errorf("error while handling pre-checkout query: %w", err).Error())
	}
}

func (b *Bot) HandleSuccessfulPayment(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil || update.Message.SuccessfulPayment == nil ||
		update.Message.From.ID == 0 {
		return
	}

	ctx = slogctx.WithField(ctx, "successful_payment", true)
	b.l.DebugContext(ctx, "handling successful payment")

	payment := update.Message.SuccessfulPayment
	err := b.s.HandleSuccessfulPayment(ctx, domain.SuccessfulPayment{
		ExternalUserID:   strconv.FormatInt(update.Message.From.ID, 10),
		Currency:         payment.Currency,
		TotalAmount:      int64(payment.TotalAmount),
		InvoicePayload:   payment.InvoicePayload,
		TelegramChargeID: payment.TelegramPaymentChargeID,
		ProviderChargeID: payment.ProviderPaymentChargeID,
	})
	if err != nil {
		b.l.ErrorContext(ctx, fmt.Errorf("error while handling successful payment: %w", err).Error())
	}
}
