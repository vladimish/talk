package tg

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/service"
	"github.com/vladimish/talk/pkg/slogctx"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot struct {
	l     *slog.Logger
	s     *service.UpdateService
	bot   *bot.Bot
	token string
}

func NewBot(l *slog.Logger, s *service.UpdateService, telegramBot *bot.Bot, token string) *Bot {
	return &Bot{
		l:     l,
		s:     s,
		bot:   telegramBot,
		token: token,
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

	// Extract image data if present
	var imageData []byte
	var imageMimeType string
	var pdfData []byte
	var pdfMimeType string
	var pdfFileName string
	messageText := update.Message.Text

	// Handle photo messages
	if len(update.Message.Photo) > 0 {
		largestPhoto := update.Message.Photo[len(update.Message.Photo)-1]
		imageData, imageMimeType = b.downloadPhoto(ctx, largestPhoto)

		// Use caption as message text if no text is provided
		if messageText == "" && update.Message.Caption != "" {
			messageText = update.Message.Caption
		}

		b.l.InfoContext(ctx, "Photo received",
			slog.String("file_id", largestPhoto.FileID),
			slog.Int("file_size", largestPhoto.FileSize),
			slog.Int("downloaded_size", len(imageData)))
	}

	// Handle document messages (for PDFs)
	if update.Message.Document != nil && update.Message.Document.MimeType == "application/pdf" {
		pdfData, pdfMimeType = b.downloadDocument(ctx, update.Message.Document)
		pdfFileName = update.Message.Document.FileName

		// Use caption as message text if no text is provided
		if messageText == "" && update.Message.Caption != "" {
			messageText = update.Message.Caption
		}

		b.l.InfoContext(ctx, "PDF document received",
			slog.String("file_id", update.Message.Document.FileID),
			slog.Int64("file_size", update.Message.Document.FileSize),
			slog.String("file_name", pdfFileName),
			slog.Int("downloaded_size", len(pdfData)))
	}

	err := b.s.HandleUpdate(ctx, domain.Update{
		ExternalID:        strconv.FormatInt(update.ID, 10),
		ExternalUserID:    strconv.FormatInt(update.Message.From.ID, 10),
		UserLanguage:      update.Message.From.LanguageCode,
		MessageText:       messageText,
		ImageData:         imageData,
		ImageMimeType:     imageMimeType,
		PDFData:           pdfData,
		PDFMimeType:       pdfMimeType,
		PDFFileName:       pdfFileName,
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

func (b *Bot) downloadPhoto(ctx context.Context, photo models.PhotoSize) ([]byte, string) {
	// Download the photo using Telegram Bot API
	file, err := b.bot.GetFile(ctx, &bot.GetFileParams{
		FileID: photo.FileID,
	})
	if err != nil {
		b.l.ErrorContext(ctx, "failed to get file info", slog.String("error", err.Error()))
		return nil, ""
	}

	if file.FilePath == "" {
		return nil, ""
	}

	// Download the actual file content
	downloadURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.token, file.FilePath)

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if reqErr != nil {
		b.l.ErrorContext(ctx, "failed to create request", slog.String("error", reqErr.Error()))
		return nil, ""
	}

	httpResp, httpErr := http.DefaultClient.Do(req)
	if httpErr != nil {
		b.l.ErrorContext(ctx, "failed to download file", slog.String("error", httpErr.Error()))
		return nil, ""
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		b.l.ErrorContext(ctx, "failed to download file",
			slog.Int("status_code", httpResp.StatusCode))
		return nil, ""
	}

	// Read the file content
	fileData, readErr := io.ReadAll(httpResp.Body)
	if readErr != nil {
		b.l.ErrorContext(ctx, "failed to read file content", slog.String("error", readErr.Error()))
		return nil, ""
	}

	b.l.InfoContext(ctx, "Successfully downloaded image",
		slog.Int("size", len(fileData)))

	return fileData, "image/jpeg" // Telegram photos are always JPEG
}

func (b *Bot) downloadDocument(ctx context.Context, document *models.Document) ([]byte, string) {
	// Download the document using Telegram Bot API
	file, err := b.bot.GetFile(ctx, &bot.GetFileParams{
		FileID: document.FileID,
	})
	if err != nil {
		b.l.ErrorContext(ctx, "failed to get file info", slog.String("error", err.Error()))
		return nil, ""
	}

	if file.FilePath == "" {
		return nil, ""
	}

	// Download the actual file content
	downloadURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.token, file.FilePath)

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if reqErr != nil {
		b.l.ErrorContext(ctx, "failed to create request", slog.String("error", reqErr.Error()))
		return nil, ""
	}

	httpResp, httpErr := http.DefaultClient.Do(req)
	if httpErr != nil {
		b.l.ErrorContext(ctx, "failed to download file", slog.String("error", httpErr.Error()))
		return nil, ""
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		b.l.ErrorContext(ctx, "failed to download file",
			slog.Int("status_code", httpResp.StatusCode))
		return nil, ""
	}

	// Read the file content
	fileData, readErr := io.ReadAll(httpResp.Body)
	if readErr != nil {
		b.l.ErrorContext(ctx, "failed to read file content", slog.String("error", readErr.Error()))
		return nil, ""
	}

	b.l.InfoContext(ctx, "Successfully downloaded document",
		slog.Int("size", len(fileData)),
		slog.String("mime_type", document.MimeType))

	return fileData, document.MimeType
}
