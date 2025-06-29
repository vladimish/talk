package domain

import "time"

type MessageSender string

const (
	MessageSenderUser MessageSender = "user"
	MessageSenderBot  MessageSender = "bot"
)

type Message struct {
	ID             int64
	UserID         int64
	MessageType    MessageType
	SentBy         MessageSender
	ConversationID *int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type MessageType struct {
	Text          string `json:"text"`
	ImageData     []byte `json:"image_data"`
	ImageMimeType string `json:"image_mime_type"`
	PDFData       []byte `json:"pdf_data"`
	PDFMimeType   string `json:"pdf_mime_type"`
	PDFFileName   string `json:"pdf_filename"`
}
