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
	Text string `json:"text"`
}
