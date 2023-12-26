package translator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hrvadl/go-weekly/internal/crawler"
)

const LingvaAPIURL = "https://lingva.ml/api/v1/en/uk/"

type Config struct {
	Timeout         time.Duration
	Retries         int
	RetriesInterval time.Duration
	URL             string
}

func NewLingvaClient(cfg *Config) *LingvaClient {
	return &LingvaClient{
		Retries:         cfg.Retries,
		RetriesInterval: cfg.RetriesInterval,
		url:             cfg.URL,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

type LingvaResponse struct {
	Translation string `json:"translation"`
}

type LingvaClient struct {
	Retries         int
	RetriesInterval time.Duration

	client *http.Client
	url    string
}

func (c *LingvaClient) Translate(msg string) (string, error) {
	var (
		err  error
		res  *http.Response
		body []byte
	)

	req, err := http.NewRequest(http.MethodGet, c.url+url.QueryEscape(msg), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	for i := 0; i <= c.Retries; i++ {
		res, err = c.client.Do(req)
		if err != nil || (res != nil && res.StatusCode != http.StatusOK) {
			err = fmt.Errorf("failed to translate, status: %v, err: %w", res.StatusCode, err)
			time.Sleep(c.RetriesInterval)
			continue
		}

		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			time.Sleep(c.RetriesInterval)
			continue
		}

		var result LingvaResponse
		if err = json.Unmarshal(body, &result); err != nil {
			time.Sleep(c.RetriesInterval)
			continue
		}

		s, _ := url.QueryUnescape(result.Translation)
		return s, nil
	}

	return "", err
}

// NOTE: This function MUST be sync, without any goroutines
// because we're getting 429 otherwise
func (c *LingvaClient) TranslateArticles(articles []crawler.Article) error {
	for i := 0; i < len(articles); i++ {
		translated, err := c.Translate(articles[i].Content)
		if err != nil {
			return err
		}
		articles[i].Content = translated
	}

	return nil
}
