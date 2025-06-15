package domain

import "time"

type Attachment struct {
	ID          int64
	MessageID   int64
	S3Name      string
	ContentType string
	Size        int64
	CreatedAt   time.Time
}
