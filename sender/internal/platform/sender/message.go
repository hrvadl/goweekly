package sender

type Message struct {
	Message   string `json:"text"`
	ChatID    string `json:"chat_id"`
	ParseMode string `json:"parse_mode"`
}
