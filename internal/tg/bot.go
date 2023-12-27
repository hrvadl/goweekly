package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/hrvadl/go-weekly/pkg/logger"
)

const URL = "https://api.telegram.org/bot"
const ContentTypeJSON = "application/json"
const parseMode = "Markdown"
const daysInWeek = 7

func NewBot(url, token, chatID string) Bot {
	return Bot{
		url:    url + token,
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
	perDay := int(math.Ceil(float64(len(messages)) / daysInWeek))

	var wg sync.WaitGroup
	for idx, msg := range messages {
		if b.dayLimitExceeded(idx, perDay) {
			time.Sleep(time.Hour * 24)
		}

		wg.Add(1)
		go func(msg string) {
			defer wg.Done()
			if err := b.SendMessage(msg); err != nil {
				logger.Errorf("Failed to send message: %v", err)
			}
		}(msg)
	}

	wg.Wait()
}

func (b Bot) dayLimitExceeded(idx, perDay int) bool {
	return idx != 0 && idx%perDay == 0
}
