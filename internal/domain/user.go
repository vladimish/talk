package domain

import "time"

const (
	UserStateMenu         = "menu"
	UserStateConversation = "conversation"
)

const (
	ButtonStartConversation = "🗣️ Start Conversation"
	ButtonBackToMenu        = "🔙 Back to Menu"
)

type User struct {
	ID                  int64
	ExternalID          string
	Language            string
	CurrentStep         string
	SelectedModel       string
	CurrentConversation *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
