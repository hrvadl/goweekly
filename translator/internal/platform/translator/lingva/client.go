package lingva

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

const LingvaAPIURL = "https://lingva.lunar.icu/api/v1/en/uk/"

type Config struct {
	Retries         int
	RetriesInterval time.Duration
	Logger          *slog.Logger
}

func NewClient(cfg *Config) *Client {
	return &Client{
		retries:         cfg.Retries,
		retriesInterval: cfg.RetriesInterval,
		client:          http.DefaultClient,
		url:             LingvaAPIURL,
		log:             cfg.Logger,
	}
}

type LingvaResponse struct {
	Translation string `json:"translation"`
}

type Client struct {
	retries         int
	retriesInterval time.Duration
	log             *slog.Logger
	client          *http.Client
	url             string
}

func (c *Client) Translate(ctx context.Context, msg string) (string, error) {
	var (
		err error
		res string
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+url.QueryEscape(msg), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	for i := 0; i <= c.retries; i++ {
		res, err = c.translate(req)
		if err == nil {
			return res, nil
		}
		c.log.Error("Request failed, waiting for retry...", slog.Any("err", err))
		time.Sleep(c.retriesInterval)
	}

	return "", err
}

func (c *Client) translate(req *http.Request) (string, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to translate, status: %d", res.StatusCode)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read translate response body: %w", err)
	}

	var result LingvaResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse translate response body: %w", err)
	}

	s, _ := url.QueryUnescape(result.Translation)
	return s, nil
}
