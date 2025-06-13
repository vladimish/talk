package domain

import "time"

type Conversation struct {
	ID        int64
	Name      string
	UserID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
