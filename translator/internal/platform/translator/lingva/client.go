package lingva

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const LingvaAPIURL = "https://lingva.ml/api/v1/en/uk/"

type Config struct {
	BatchRequests   int
	Retries         int
	RetriesInterval time.Duration
	BatchInterval   time.Duration
	Timeout         time.Duration
}

func NewClient(cfg *Config) *Client {
	return &Client{
		BatchRequests:   cfg.BatchRequests,
		BatchInterval:   cfg.BatchInterval,
		Retries:         cfg.Retries,
		RetriesInterval: cfg.RetriesInterval,
		url:             LingvaAPIURL,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

type LingvaResponse struct {
	Translation string `json:"translation"`
}

type Client struct {
	BatchRequests   int
	Retries         int
	RetriesInterval time.Duration
	BatchInterval   time.Duration

	client *http.Client
	url    string
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

	for i := 0; i <= c.Retries; i++ {
		res, err = c.translate(req)
		if err == nil {
			return res, nil
		}
		time.Sleep(c.RetriesInterval)
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
