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

type Bot struct {
	url       string
	chatID    string
	parseMode string
}

func (b Bot) Send(ctx context.Context, msg sender.Message) error {
	errCh := make(chan error)
	doneCh := make(chan struct{})

	go func() {
		body, err := json.Marshal(msg)
		if err != nil {
			errCh <- err
			return
		}

		res, err := http.Post(b.url+"/sendMessage", ContentTypeJSON, bytes.NewBuffer(body))
		if err != nil {
			errCh <- err
			return
		}

		defer res.Body.Close()
		resp, err := io.ReadAll(res.Body)
		if err != nil {
			errCh <- err
			return
		}

		if res.StatusCode != http.StatusOK {
			err := fmt.Errorf(
				"sending message failed with status %v: %v",
				res.StatusCode,
				string(resp),
			)
			errCh <- err
			return
		}

		doneCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("timeout exceeded")
	case err := <-errCh:
		return fmt.Errorf("failed to send tg message: %w", err)
	case <-doneCh:
		return nil
	}
}
