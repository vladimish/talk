package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/queue"
)

const (
	connectionTimeout = 5 * time.Second
	queueExpiration   = 24 * time.Hour
)

type Queue struct {
	client *redis.Client
}

func NewRedisQueue(redisURL string) (*Queue, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	if pingErr := client.Ping(ctx).Err(); pingErr != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", pingErr)
	}

	return &Queue{
		client: client,
	}, nil
}

func (r *Queue) EnqueueWithNotification(
	ctx context.Context,
	userID string,
	update domain.Update,
	notificationID string,
) error {
	item := queue.QueuedItem{
		Update:              update,
		QueueNotificationID: notificationID,
	}

	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal queued item: %w", err)
	}

	queueKey := r.queueKey(userID)
	if pushErr := r.client.RPush(ctx, queueKey, data).Err(); pushErr != nil {
		return fmt.Errorf("failed to enqueue item: %w", pushErr)
	}

	// Set expiration for the queue to prevent memory leaks (24 hours)
	r.client.Expire(ctx, queueKey, queueExpiration)

	return nil
}

func (r *Queue) DequeueWithMetadata(ctx context.Context, userID string) (*queue.QueuedItem, error) {
	queueKey := r.queueKey(userID)

	data, err := r.client.LPop(ctx, queueKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, queue.ErrEmptyQueue
		}
		return nil, fmt.Errorf("failed to dequeue item: %w", err)
	}

	var item queue.QueuedItem
	if unmarshalErr := json.Unmarshal(data, &item); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal queued item: %w", unmarshalErr)
	}

	return &item, nil
}

func (r *Queue) SetProcessing(ctx context.Context, userID string, ttl time.Duration) error {
	lockKey := r.lockKey(userID)

	// Use SET with NX (only set if not exists) and EX (expiration)
	ok, err := r.client.SetNX(ctx, lockKey, "1", ttl).Result()
	if err != nil {
		return fmt.Errorf("failed to set processing lock: %w", err)
	}

	if !ok {
		return queue.ErrAlreadyProcessing
	}

	return nil
}

func (r *Queue) IsProcessing(ctx context.Context, userID string) (bool, error) {
	lockKey := r.lockKey(userID)

	exists, err := r.client.Exists(ctx, lockKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check processing lock: %w", err)
	}

	return exists > 0, nil
}

func (r *Queue) ClearProcessing(ctx context.Context, userID string) error {
	lockKey := r.lockKey(userID)

	if err := r.client.Del(ctx, lockKey).Err(); err != nil {
		return fmt.Errorf("failed to clear processing lock: %w", err)
	}

	return nil
}

func (r *Queue) GetQueueLength(ctx context.Context, userID string) (int, error) {
	queueKey := r.queueKey(userID)

	length, err := r.client.LLen(ctx, queueKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}

	return int(length), nil
}

// SetPendingMessages stores messages that are waiting for potential concatenation.
func (r *Queue) SetPendingMessages(
	ctx context.Context,
	userID string,
	messages *queue.PendingMessages,
	ttl time.Duration,
) error {
	pendingKey := r.pendingKey(userID)

	data, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to marshal pending messages: %w", err)
	}

	if setErr := r.client.Set(ctx, pendingKey, data, ttl).Err(); setErr != nil {
		return fmt.Errorf("failed to set pending messages: %w", setErr)
	}

	return nil
}

// GetPendingMessages retrieves pending messages for concatenation.
func (r *Queue) GetPendingMessages(ctx context.Context, userID string) (*queue.PendingMessages, error) {
	pendingKey := r.pendingKey(userID)

	data, err := r.client.Get(ctx, pendingKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, queue.ErrEmptyQueue // No pending messages
		}
		return nil, fmt.Errorf("failed to get pending messages: %w", err)
	}

	var messages queue.PendingMessages
	if unmarshalErr := json.Unmarshal(data, &messages); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal pending messages: %w", unmarshalErr)
	}

	return &messages, nil
}

// ClearPendingMessages removes pending messages for a user.
func (r *Queue) ClearPendingMessages(ctx context.Context, userID string) error {
	pendingKey := r.pendingKey(userID)

	if err := r.client.Del(ctx, pendingKey).Err(); err != nil {
		return fmt.Errorf("failed to clear pending messages: %w", err)
	}

	return nil
}

// SetGenerationLock marks a user as having an active AI generation (for cancellation).
func (r *Queue) SetGenerationLock(ctx context.Context, userID string, ttl time.Duration) error {
	genKey := r.generationKey(userID)

	// Use SET with EX (expiration)
	if err := r.client.Set(ctx, genKey, "1", ttl).Err(); err != nil {
		return fmt.Errorf("failed to set generation lock: %w", err)
	}

	return nil
}

// IsGenerating checks if a user has an active AI generation.
func (r *Queue) IsGenerating(ctx context.Context, userID string) (bool, error) {
	genKey := r.generationKey(userID)

	exists, err := r.client.Exists(ctx, genKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check generation lock: %w", err)
	}

	return exists > 0, nil
}

// ClearGenerationLock removes the generation lock for a user.
func (r *Queue) ClearGenerationLock(ctx context.Context, userID string) error {
	genKey := r.generationKey(userID)

	if err := r.client.Del(ctx, genKey).Err(); err != nil {
		return fmt.Errorf("failed to clear generation lock: %w", err)
	}

	return nil
}

func (r *Queue) Close() error {
	return r.client.Close()
}

// Helper methods for key generation.
func (r *Queue) queueKey(userID string) string {
	return fmt.Sprintf("queue:user:%s", userID)
}

func (r *Queue) lockKey(userID string) string {
	return fmt.Sprintf("lock:user:%s", userID)
}

func (r *Queue) pendingKey(userID string) string {
	return fmt.Sprintf("pending:user:%s", userID)
}

func (r *Queue) generationKey(userID string) string {
	return fmt.Sprintf("generation:user:%s", userID)
}
