package domain

import "time"

const (
	UserStateMenu             = "menu"
	UserStateConversation     = "conversation"
	UserStateModelSelect      = "model_select"
	UserStateConversationList = "conversation_list"
	UserStateSettings         = "settings"
	UserStateLanguageSelect   = "language_select"
)

const (
	ButtonStartConversation = "💬️ Start Conversation"
	ButtonModelSelect       = "🤖 Select Model"
	ButtonBackToMenu        = "🔙 Back to Menu"
	ButtonNewConversation   = "➕ New Conversation"
	ButtonSettings          = "⚙️ Settings"
	ButtonLanguage          = "🌐 Language"
)

var AvailableModels = []string{
	"google/gemini-2.5-flash-preview-05-20",
	"deepseek/deepseek-r1-0528-qwen3-8b",
	"thedrummer/valkyrie-49b-v1",
}

var SupportedLanguages = map[string]string{
	"🇺🇸 English":   "en",
	"🇪🇸 Español":   "es",
	"🇫🇷 Français":  "fr",
	"🇩🇪 Deutsch":   "de",
	"🇮🇹 Italiano":  "it",
	"🇷🇺 Русский":   "ru",
	"🇨🇳 中文":        "zh",
	"🇯🇵 日本語":       "ja",
	"🇰🇷 한국어":       "ko",
	"🇵🇹 Português": "pt",
}

type User struct {
	ID                    int64
	ExternalID            string
	Language              string
	CurrentStep           string
	SelectedModel         string
	CurrentConversationID *int64
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
