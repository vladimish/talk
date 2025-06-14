package queue

import (
	"context"
	"errors"
	"time"

	"github.com/vladimish/talk/internal/domain"
)

var ErrAlreadyProcessing = errors.New("user is already processing a request")

// QueuedItem represents an update with metadata
type QueuedItem struct {
	Update              domain.Update
	QueueNotificationID string // Message ID of the "queued" notification to delete
}

// Queue interface for managing user request queues
type Queue interface {
	// EnqueueWithNotification adds an update to a user's queue with notification message ID
	EnqueueWithNotification(ctx context.Context, userID string, update domain.Update, notificationID string) error

	// DequeueWithMetadata retrieves and removes the next update from a user's queue
	// Returns nil if queue is empty
	DequeueWithMetadata(ctx context.Context, userID string) (*QueuedItem, error)

	// SetProcessing marks a user as currently processing a request
	// ttl defines how long the lock should be held
	SetProcessing(ctx context.Context, userID string, ttl time.Duration) error

	// IsProcessing checks if a user is currently processing a request
	IsProcessing(ctx context.Context, userID string) (bool, error)

	// ClearProcessing removes the processing lock for a user
	ClearProcessing(ctx context.Context, userID string) error

	// GetQueueLength returns the number of queued updates for a user
	GetQueueLength(ctx context.Context, userID string) (int, error)
}
