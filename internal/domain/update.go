package domain

type Update struct {
	ExternalID        string
	ExternalUserID    string
	UserLanguage      string
	MessageText       string
	ExternalMessageID int // Telegram message ID
}
