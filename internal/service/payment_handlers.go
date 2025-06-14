package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vladimish/talk/internal/domain"
)

// grantSubscriptionTokens grants tokens to the user based on their subscription.
func (s *UpdateService) grantSubscriptionTokens(ctx context.Context, payment *domain.Payment) error {
	now := time.Now()

	// Grant regular tokens
	regularTransaction := &domain.Transaction{
		UserID:          payment.UserID,
		TokenType:       domain.TokenTypeRegular,
		Amount:          domain.MonthlyRegularTokenReward,
		TransactionType: domain.TransactionTypeAdminCredit,
		ModelUsed:       nil,
		Description:     stringPtr("Monthly subscription reward - regular tokens"),
		CreatedAt:       now,
	}

	_, err := s.storage.CreateTransaction(ctx, regularTransaction)
	if err != nil {
		return fmt.Errorf("failed to create regular token transaction: %w", err)
	}

	// Grant premium tokens
	premiumTransaction := &domain.Transaction{
		UserID:          payment.UserID,
		TokenType:       domain.TokenTypePremium,
		Amount:          domain.MonthlyPremiumTokenReward,
		TransactionType: domain.TransactionTypeAdminCredit,
		ModelUsed:       nil,
		Description:     stringPtr("Monthly subscription reward - premium tokens"),
		CreatedAt:       now,
	}

	_, err = s.storage.CreateTransaction(ctx, premiumTransaction)
	if err != nil {
		return fmt.Errorf("failed to create premium token transaction: %w", err)
	}

	return nil
}

// stringPtr returns a pointer to the given string.
func stringPtr(s string) *string {
	return &s
}
