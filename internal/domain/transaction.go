package domain

import (
	"time"

	"github.com/vladimish/talk/pkg/i18n"
)

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
	ID             string    `json:"id"`                    // Internal model identifier
	DisplayName    string    `json:"display_name"`          // User-friendly name for buttons (deprecated - use I18nKey)
	I18nKey        string    `json:"i18n_key"`              // Internationalization key for display name
	Cost           int64     `json:"cost"`                  // Token cost per message
	TokenType      TokenType `json:"token_type"`            // Type of tokens required
	ImageSupport   bool      `json:"image_support"`         // Whether the model supports image inputs
	Reasoning      bool      `json:"reasoning"`             // Whether the model has reasoning capabilities
	WebSearch      bool      `json:"web_search"`            // Whether the model has web search capabilities
	NoSubscription bool      `json:"no_subscription"`       // If true, requires active subscription to use
	SearchCost     *int64    `json:"search_cost,omitempty"` // Additional cost when using web search (optional)
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
		ID:             "google/gemini-2.5-flash-preview-05-20",
		I18nKey:        i18n.ModelGeminiFlash,
		Cost:           TokenCostLow,
		TokenType:      TokenTypeRegular,
		ImageSupport:   true,
		Reasoning:      false,
		WebSearch:      true,
		NoSubscription: false,
		SearchCost:     &[]int64{2}[0], // 2 additional tokens when using web search
	},
	{
		ID:             "openai/gpt-4o",
		I18nKey:        i18n.ModelGPT4o,
		Cost:           TokenCostLow,
		TokenType:      TokenTypePremium,
		ImageSupport:   true,
		Reasoning:      false,
		WebSearch:      false,
		NoSubscription: false,
	},
	{
		ID:             "openai/gpt-4o-mini",
		I18nKey:        i18n.ModelGPT4oMini,
		Cost:           TokenCostLow,
		TokenType:      TokenTypeRegular,
		ImageSupport:   true,
		Reasoning:      false,
		WebSearch:      false,
		NoSubscription: false,
	},
	{
		ID:             "anthropic/claude-4-sonnet-20250522",
		I18nKey:        i18n.ModelClaude4Sonnet,
		Cost:           TokenCostLow,
		TokenType:      TokenTypePremium,
		ImageSupport:   true,
		Reasoning:      false,
		WebSearch:      false,
		NoSubscription: false,
	},
	{
		ID:             "google/gemini-2.5-pro-preview",
		I18nKey:        i18n.ModelGeminiPro,
		Cost:           TokenCostLow,
		TokenType:      TokenTypeRegular,
		ImageSupport:   true,
		Reasoning:      false,
		WebSearch:      true,
		NoSubscription: false,
	},
	{
		ID:             "openai/o3-mini",
		I18nKey:        i18n.ModelO3Mini,
		Cost:           TokenCostMedium,
		TokenType:      TokenTypePremium,
		ImageSupport:   false,
		Reasoning:      true,
		WebSearch:      false,
		NoSubscription: true,
	},
	{
		ID:             "deepseek/deepseek-chat-v3-0324:free",
		I18nKey:        i18n.ModelDeepSeekChat,
		Cost:           TokenCostLow,
		TokenType:      TokenTypeRegular,
		ImageSupport:   false,
		Reasoning:      false,
		WebSearch:      false,
		NoSubscription: false,
	},
	{
		ID:             "deepseek/deepseek-r1:free",
		I18nKey:        i18n.ModelDeepSeekR1,
		Cost:           TokenCostMedium,
		TokenType:      TokenTypeRegular,
		ImageSupport:   false,
		Reasoning:      true,
		WebSearch:      false,
		NoSubscription: false,
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

// GetDisplayNameWithEmojis returns the model display name with appropriate capability emojis.
func (m *ModelInfo) GetDisplayNameWithEmojis(language string) string {
	// Get localized name or fallback to DisplayName
	name := i18n.GetString(language, m.I18nKey)
	if name == m.I18nKey && m.DisplayName != "" {
		// Fallback to DisplayName if i18n key not found
		name = m.DisplayName
	}

	// Add capability emojis
	if m.ImageSupport {
		name += " 👁️"
	}
	if m.Reasoning {
		name += " 🧠"
	}
	if m.WebSearch {
		name += " 🌐"
	}

	return name
}

// CanUserUseModel checks if a user can use a specific model based on subscription and token balance.
func CanUserUseModel(model *ModelInfo, hasActiveSubscription bool, tokenBalance TokenBalance) bool {
	// Check subscription requirement
	if model.NoSubscription && !hasActiveSubscription {
		return false
	}

	// Check token balance requirement
	switch model.TokenType {
	case TokenTypePremium:
		return tokenBalance.PremiumBalance >= model.Cost
	case TokenTypeRegular:
		return tokenBalance.RegularBalance >= model.Cost
	}

	return false
}

// GetDisplayNameForUser returns the model display name with availability indicator for a specific user.
func (m *ModelInfo) GetDisplayNameForUser(
	language string,
	hasActiveSubscription bool,
	tokenBalance TokenBalance,
) string {
	displayName := m.GetDisplayNameWithEmojis(language)

	// Add red cross if user can't use the model
	if !CanUserUseModel(m, hasActiveSubscription, tokenBalance) {
		displayName = "❌ " + displayName
	}

	return displayName
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
