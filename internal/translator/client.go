package translator

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/hrvadl/go-weekly/internal/crawler"
	"github.com/hrvadl/go-weekly/pkg/logger"
)

const LingvaAPIURL = "https://lingva.ml/api/v1/en/uk/"

type Config struct {
	BatchRequests   int
	Retries         int
	RetriesInterval time.Duration
	BatchInterval   time.Duration
	Timeout         time.Duration
}

func NewLingvaClient(cfg *Config) *LingvaClient {
	return &LingvaClient{
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

type LingvaClient struct {
	BatchRequests   int
	Retries         int
	RetriesInterval time.Duration
	BatchInterval   time.Duration

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
		if err != nil {
			time.Sleep(c.RetriesInterval)
			continue
		}

		if res.StatusCode != http.StatusOK {
			logger.Errorf("failed to translate, status: %v, err: %v", res.StatusCode, err)
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

func (c *LingvaClient) TranslateArticles(articles []crawler.Article) error {
	var (
		wg     sync.WaitGroup
		errCh  = make(chan error, len(articles))
		doneCh = make(chan struct{}, c.BatchRequests)
	)

	for i := 0; i < len(articles); i++ {
		if c.isStartOfTheChunk(i) {
			c.waitForPreviousBatch(doneCh)
		}

		wg.Add(1)
		go func(article *crawler.Article) {
			defer wg.Done()
			translated, err := c.Translate(article.Content)
			doneCh <- struct{}{}

			if err != nil {
				errCh <- err
			}

			article.Content = translated
		}(&articles[i])
	}

	wg.Wait()

	var err error
	for i := 0; i < len(errCh); i++ {
		err = errors.Join(err, <-errCh)
	}

	return err
}

func (c *LingvaClient) isStartOfTheChunk(i int) bool {
	return i != 0 && i%c.BatchRequests == 0
}

func (c *LingvaClient) waitForPreviousBatch(doneCh <-chan struct{}) {
	for j := 0; j < c.BatchRequests; j++ {
		<-doneCh
	}
	time.Sleep(c.BatchInterval)
}
