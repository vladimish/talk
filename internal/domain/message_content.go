package domain

type KeyboardButton struct {
	Text string
}

type ReplyKeyboard struct {
	Buttons     [][]KeyboardButton
	Resize      bool
	OneTime     bool
	Placeholder string
}

type InlineKeyboardButton struct {
	Text         string
	URL          string
	CallbackData string
}

type InlineKeyboard struct {
	Buttons [][]InlineKeyboardButton
}

type MessageContent struct {
	Text             string
	ReplyKeyboard    *ReplyKeyboard
	InlineKeyboard   *InlineKeyboard
	ReplyToMessageID *int64
	IsPersistent     bool
}
