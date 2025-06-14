package domain

import "time"

const (
	UserStateMenu             = "menu"
	UserStateConversation     = "conversation"
	UserStateModelSelect      = "model_select"
	UserStateConversationList = "conversation_list"
	UserStateSettings         = "settings"
	UserStateLanguageSelect   = "language_select"
	UserStateProfile          = "profile"
)

type User struct {
	ID                     int64
	ExternalID             string
	Language               string
	CurrentStep            string
	SelectedModel          string
	CurrentConversationID  *int64
	ConversationListOffset int
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
