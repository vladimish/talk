package queue

import (
	"context"
	"errors"
	"time"

	"github.com/vladimish/talk/internal/domain"
)

//go:generate go tool mockgen -source=queue.go -destination=../../../mocks/mock_queue.go -package=mocks

var (
	ErrAlreadyProcessing = errors.New("user is already processing a request")
	ErrEmptyQueue        = errors.New("queue is empty")
)

// QueuedItem represents an update with metadata.
type QueuedItem struct {
	Update              domain.Update `json:"update"`
	QueueNotificationID string        `json:"notification_id"` // Message ID of the "queued" notification to delete
}

// PendingMessages represents messages waiting for concatenation.
type PendingMessages struct {
	Messages  []domain.Update `json:"messages"`
	Timestamp time.Time       `json:"timestamp"`
}

// Queue interface for managing user request queues.
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

	// Message concatenation methods
	// SetPendingMessages stores messages that are waiting for potential concatenation
	SetPendingMessages(ctx context.Context, userID string, messages *PendingMessages, ttl time.Duration) error

	// GetPendingMessages retrieves pending messages for concatenation
	GetPendingMessages(ctx context.Context, userID string) (*PendingMessages, error)

	// ClearPendingMessages removes pending messages for a user
	ClearPendingMessages(ctx context.Context, userID string) error

	// SetGenerationLock marks a user as having an active AI generation (for cancellation)
	SetGenerationLock(ctx context.Context, userID string, ttl time.Duration) error

	// IsGenerating checks if a user has an active AI generation
	IsGenerating(ctx context.Context, userID string) (bool, error)

	// ClearGenerationLock removes the generation lock for a user
	ClearGenerationLock(ctx context.Context, userID string) error
}
