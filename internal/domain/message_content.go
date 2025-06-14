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

type MessageContent struct {
	Text             string
	ReplyKeyboard    *ReplyKeyboard
	ReplyToMessageID *int64
}
