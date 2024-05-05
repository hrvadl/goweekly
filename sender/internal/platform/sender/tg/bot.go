package tg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hrvadl/goweekly/sender/internal/platform/sender"
)

const (
	URL             = "https://api.telegram.org/bot"
	ContentTypeJSON = "application/json"
	daysInWeek      = 7
)

func NewBot(token, chatID, parseMode string) Bot {
	return Bot{
		url:       URL + token,
		chatID:    chatID,
		parseMode: parseMode,
	}
}

type MessagePayload struct {
	Message   string `json:"text"`
	ChatID    string `json:"chat_id"`
	ParseMode string `json:"parse_mode"`
}

type Bot struct {
	url       string
	chatID    string
	parseMode string
}

func (b Bot) Send(ctx context.Context, msg sender.Message) error {
	body, err := json.Marshal(MessagePayload{
		Message:   msg.Message,
		ChatID:    b.chatID,
		ParseMode: b.parseMode,
	})
	if err != nil {
		return fmt.Errorf("failed to construct request body: %w", err)
	}

	res, err := http.Post(b.url+"/sendMessage", ContentTypeJSON, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer res.Body.Close()
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"sending message failed with status %v: %v",
			res.StatusCode,
			string(resp),
		)
	}

	return nil
}
