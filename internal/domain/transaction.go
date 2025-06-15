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
	Amount          int64           `json:"amount"` // positive for credits, negative for debits
	TransactionType TransactionType `json:"transaction_type"`
	ModelUsed       *string         `json:"model_used,omitempty"`
	Description     *string         `json:"description,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

// TokenBalance represents user token balances.
type TokenBalance struct {
	PremiumBalance int64 `json:"premium_balance"`
	RegularBalance int64 `json:"regular_balance"`
}

// ModelInfo represents the complete information for an AI model.
type ModelInfo struct {
	ID           string    `json:"id"`            // Internal model identifier
	DisplayName  string    `json:"display_name"`  // User-friendly name for buttons
	Cost         int64     `json:"cost"`          // Token cost per message
	TokenType    TokenType `json:"token_type"`    // Type of tokens required
	ImageSupport bool      `json:"image_support"` // Whether the model supports image inputs
}

// Token cost constants.
const (
	TokenCostLow    int64 = 1
	TokenCostMedium int64 = 3
	TokenCostHigh   int64 = 5
)

// AvailableModels defines all available AI models with their complete information.
var AvailableModels = []ModelInfo{
	{
		ID:           "google/gemini-2.5-flash-preview-05-20",
		DisplayName:  "🚀 Gemini 2.5 Flash (Fast) 👁️",
		Cost:         TokenCostLow,
		TokenType:    TokenTypeRegular,
		ImageSupport: true,
	},
	{
		ID:           "openai/gpt-4o",
		DisplayName:  "🧠 GPT-4o (Most Capable) 👁️",
		Cost:         TokenCostHigh,
		TokenType:    TokenTypePremium,
		ImageSupport: true,
	},
	{
		ID:           "openai/gpt-4o-mini",
		DisplayName:  "⚡ GPT-4o Mini (Balanced) 👁️",
		Cost:         TokenCostLow,
		TokenType:    TokenTypeRegular,
		ImageSupport: true,
	},
	{
		ID:           "anthropic/claude-4-sonnet-20250522",
		DisplayName:  "🎭 Claude 4 Sonnet (Creative) 👁️",
		Cost:         TokenCostMedium,
		TokenType:    TokenTypePremium,
		ImageSupport: true,
	},
	{
		ID:           "google/gemini-2.5-pro-preview",
		DisplayName:  "🌸 Gemini 2.5 Pro (Long context) 👁️",
		Cost:         TokenCostLow,
		TokenType:    TokenTypeRegular,
		ImageSupport: true,
	},
}

// GetModelByID returns the model info for a given model ID.
func GetModelByID(modelID string) *ModelInfo {
	for _, model := range AvailableModels {
		if model.ID == modelID {
			return &model
		}
	}
	return nil
}

// ModelCosts defines the token costs for each AI model (for backward compatibility).
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

// ModelCost represents the cost configuration for AI models (for backward compatibility).
type ModelCost struct {
	ModelName string    `json:"model_name"`
	Cost      int64     `json:"cost"`
	TokenType TokenType `json:"token_type"`
}
