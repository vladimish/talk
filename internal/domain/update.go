package domain

type Update struct {
	ExternalID        string `json:"external_id"`
	ExternalUserID    string `json:"external_user_id"`
	UserLanguage      string `json:"user_language"`
	MessageText       string `json:"message_text"`
	ExternalMessageID int    `json:"external_message_id"` // Telegram message ID
}
