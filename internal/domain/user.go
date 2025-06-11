package domain

import "time"

type User struct {
	ID         string
	ExternalID string
	Language   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
