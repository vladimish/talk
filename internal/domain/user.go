package domain

import "time"

type User struct {
	ID          int64
	ExternalID  string
	Language    string
	CurrentStep *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
