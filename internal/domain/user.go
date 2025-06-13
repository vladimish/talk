package domain

import "time"

const (
	UserStateMenu         = "menu"
	UserStateConversation = "conversation"
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
