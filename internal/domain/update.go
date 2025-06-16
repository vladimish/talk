package domain

type CallbackQuery struct {
	ID             string `json:"id"`
	ExternalUserID string `json:"external_user_id"`
	UserLanguage   string `json:"user_language"`
	Data           string `json:"data"`
}

type PreCheckoutQuery struct {
	ID               string `json:"id"`
	ExternalUserID   string `json:"external_user_id"`
	Currency         string `json:"currency"`
	TotalAmount      int64  `json:"total_amount"`
	InvoicePayload   string `json:"invoice_payload"`
	TelegramChargeID string `json:"telegram_charge_id"`
	ProviderChargeID string `json:"provider_charge_id"`
}

type SuccessfulPayment struct {
	ExternalUserID   string `json:"external_user_id"`
	Currency         string `json:"currency"`
	TotalAmount      int64  `json:"total_amount"`
	InvoicePayload   string `json:"invoice_payload"`
	TelegramChargeID string `json:"telegram_charge_id"`
	ProviderChargeID string `json:"provider_charge_id"`
}

type Update struct {
	ExternalID        string             `json:"external_id"`
	ExternalUserID    string             `json:"external_user_id"`
	UserLanguage      string             `json:"user_language"`
	MessageText       string             `json:"message_text"`
	ImageData         []byte             `json:"image_data,omitempty"`      // Base64 encoded image data
	ImageMimeType     string             `json:"image_mime_type,omitempty"` // MIME type of the image
	PDFData           []byte             `json:"pdf_data,omitempty"`        // Base64 encoded PDF data
	PDFMimeType       string             `json:"pdf_mime_type,omitempty"`   // MIME type of the PDF
	PDFFileName       string             `json:"pdf_filename,omitempty"`    // Original file name
	ExternalMessageID int                `json:"external_message_id"`       // Telegram message ID
	CallbackQuery     *CallbackQuery     `json:"callback_query,omitempty"`
	PreCheckoutQuery  *PreCheckoutQuery  `json:"pre_checkout_query,omitempty"`
	SuccessfulPayment *SuccessfulPayment `json:"successful_payment,omitempty"`
}
