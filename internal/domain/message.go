package domain

type Message struct {
	ID             string
	ExternalID     string
	ExternalUserID string
	Content        MessageContent
}
