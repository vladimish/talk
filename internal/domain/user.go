package domain

import "time"

const (
	UserStateMenu         = "menu"
	UserStateConversation = "conversation"
	UserStateModelSelect  = "model_select"
)

const (
	ButtonStartConversation = "💬️ Start Conversation"
	ButtonModelSelect       = "🤖 Select Model"
	ButtonBackToMenu        = "🔙 Back to Menu"
)

var AvailableModels = []string{
	"google/gemini-2.5-flash-preview-05-20",
	"deepseek/deepseek-r1-0528-qwen3-8b",
	"thedrummer/valkyrie-49b-v1",
}

type User struct {
	ID                  int64
	ExternalID          string
	Language            string
	CurrentStep         string
	SelectedModel       string
	CurrentConversation *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
