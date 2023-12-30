package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/hrvadl/go-weekly/pkg/logger"
)

const URL = "https://api.telegram.org/bot"
const ContentTypeJSON = "application/json"
const parseMode = "Markdown"
const daysInWeek = 7

var daysTillTuesday = map[time.Weekday]float64{
	time.Monday:    1.0,
	time.Tuesday:   1.0,
	time.Wednesday: 6.0,
	time.Thursday:  5.0,
	time.Friday:    4.0,
	time.Saturday:  3.0,
	time.Sunday:    2.0,
}

func NewBot(token, chatID string) Bot {
	return Bot{
		url:    URL + token,
		chatID: chatID,
	}
}

type Bot struct {
	url    string
	chatID string
}

type MessagePayload struct {
	Message   string `json:"text"`
	ChatID    string `json:"chat_id"`
	ParseMode string `json:"parse_mode"`
}

func (b Bot) SendMessage(msg string) error {
	body, err := json.Marshal(MessagePayload{
		Message:   msg,
		ChatID:    b.chatID,
		ParseMode: parseMode,
	})

	if err != nil {
		return err
	}

	res, err := http.Post(b.url+"/sendMessage", ContentTypeJSON, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("sending message failed with status %v: %v", res.StatusCode, resp)
	}

	return nil
}

func (b Bot) SendMessagesThroughoutWeek(messages []string) {
	today := time.Now().Weekday()
	perDay := int(math.Ceil(float64(len(messages)) / daysTillTuesday[today]))
	for idx, msg := range messages {
		if b.dayLimitExceeded(idx, perDay) {
			time.Sleep(time.Hour * 24)
		}

		go func(msg string) {
			if err := b.SendMessage(msg); err != nil {
				logger.Errorf("Failed to send message: %v", err)
			}
		}(msg)
	}
}

func (b Bot) dayLimitExceeded(idx, perDay int) bool {
	return idx != 0 && idx%perDay == 0
}
