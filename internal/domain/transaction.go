package domain

import "time"

// TokenType represents the type of token.
type TokenType string

const (
	TokenTypePremium TokenType = "premium"
	TokenTypeRegular TokenType = "regular"
)

// TransactionType represents the type of transaction.
type TransactionType string

const (
	TransactionTypeInitialCredit TransactionType = "initial_credit"
	TransactionTypeMessageCost   TransactionType = "message_cost"
	TransactionTypeAdminCredit   TransactionType = "admin_credit"
	TransactionTypeAdminDebit    TransactionType = "admin_debit"
)

// Transaction represents a token transaction.
type Transaction struct {
	ID              int64           `json:"id"`
	UserID          int64           `json:"user_id"`
	TokenType       TokenType       `json:"token_type"`
	Amount          int32           `json:"amount"` // positive for credits, negative for debits
	TransactionType TransactionType `json:"transaction_type"`
	ModelUsed       *string         `json:"model_used,omitempty"`
	Description     *string         `json:"description,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

// TokenBalance represents user token balances.
type TokenBalance struct {
	PremiumBalance int32 `json:"premium_balance"`
	RegularBalance int32 `json:"regular_balance"`
}

// ModelCost represents the cost configuration for AI models.
type ModelCost struct {
	ModelName string    `json:"model_name"`
	Cost      int32     `json:"cost"`
	TokenType TokenType `json:"token_type"`
}

// Token cost constants.
const (
	TokenCostLow    int32 = 1
	TokenCostMedium int32 = 3
	TokenCostHigh   int32 = 5
)

// ModelCosts defines the token costs for each AI model.
var ModelCosts = map[string]ModelCost{
	"google/gemini-2.5-flash-preview-05-20": {
		ModelName: "google/gemini-2.5-flash-preview-05-20",
		Cost:      TokenCostLow,
		TokenType: TokenTypeRegular,
	},
	"openai/gpt-4o": {
		ModelName: "openai/gpt-4o",
		Cost:      TokenCostHigh,
		TokenType: TokenTypePremium,
	},
	"openai/gpt-4o-mini": {
		ModelName: "openai/gpt-4o-mini",
		Cost:      TokenCostLow,
		TokenType: TokenTypeRegular,
	},
	"anthropic/claude-3.5-sonnet": {
		ModelName: "anthropic/claude-3.5-sonnet",
		Cost:      TokenCostMedium,
		TokenType: TokenTypePremium,
	},
	"anthropic/claude-3.5-haiku": {
		ModelName: "anthropic/claude-3.5-haiku",
		Cost:      TokenCostLow,
		TokenType: TokenTypeRegular,
	},
}
