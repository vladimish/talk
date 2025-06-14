package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/queue"
)

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(redisURL string) (*RedisQueue, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisQueue{
		client: client,
	}, nil
}

func (r *RedisQueue) EnqueueWithNotification(ctx context.Context, userID string, update domain.Update, notificationID string) error {
	item := queue.QueuedItem{
		Update:              update,
		QueueNotificationID: notificationID,
	}

	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal queued item: %w", err)
	}

	queueKey := r.queueKey(userID)
	if err := r.client.RPush(ctx, queueKey, data).Err(); err != nil {
		return fmt.Errorf("failed to enqueue item: %w", err)
	}

	// Set expiration for the queue to prevent memory leaks (24 hours)
	r.client.Expire(ctx, queueKey, 24*time.Hour)

	return nil
}

func (r *RedisQueue) DequeueWithMetadata(ctx context.Context, userID string) (*queue.QueuedItem, error) {
	queueKey := r.queueKey(userID)

	data, err := r.client.LPop(ctx, queueKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Empty queue
		}
		return nil, fmt.Errorf("failed to dequeue item: %w", err)
	}

	var item queue.QueuedItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal queued item: %w", err)
	}

	return &item, nil
}

func (r *RedisQueue) SetProcessing(ctx context.Context, userID string, ttl time.Duration) error {
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

func (r *RedisQueue) IsProcessing(ctx context.Context, userID string) (bool, error) {
	lockKey := r.lockKey(userID)

	exists, err := r.client.Exists(ctx, lockKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check processing lock: %w", err)
	}

	return exists > 0, nil
}

func (r *RedisQueue) ClearProcessing(ctx context.Context, userID string) error {
	lockKey := r.lockKey(userID)

	if err := r.client.Del(ctx, lockKey).Err(); err != nil {
		return fmt.Errorf("failed to clear processing lock: %w", err)
	}

	return nil
}

func (r *RedisQueue) GetQueueLength(ctx context.Context, userID string) (int, error) {
	queueKey := r.queueKey(userID)

	length, err := r.client.LLen(ctx, queueKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}

	return int(length), nil
}

func (r *RedisQueue) Close() error {
	return r.client.Close()
}

// Helper methods for key generation
func (r *RedisQueue) queueKey(userID string) string {
	return fmt.Sprintf("queue:user:%s", userID)
}

func (r *RedisQueue) lockKey(userID string) string {
	return fmt.Sprintf("lock:user:%s", userID)
}
