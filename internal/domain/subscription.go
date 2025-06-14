package domain

import "time"

// SubscriptionStatus represents the status of a subscription.
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
)

// Subscription represents a user subscription record.
type Subscription struct {
	ID               int64              `json:"id"`
	UserID           int64              `json:"user_id"`
	PaymentID        int64              `json:"payment_id"`
	SubscriptionType SubscriptionType   `json:"subscription_type"`
	ValidFrom        time.Time          `json:"valid_from"`
	ValidTo          time.Time          `json:"valid_to"`
	Status           SubscriptionStatus `json:"status"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

// IsActive returns true if the subscription is active and not expired.
func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive && s.ValidTo.After(time.Now())
}

// DaysRemaining returns the number of days remaining in the subscription.
func (s *Subscription) DaysRemaining() int {
	if !s.IsActive() {
		return 0
	}

	remaining := time.Until(s.ValidTo)
	const hoursPerDay = 24
	days := int(remaining.Hours() / hoursPerDay)
	if days < 0 {
		return 0
	}
	return days
}
