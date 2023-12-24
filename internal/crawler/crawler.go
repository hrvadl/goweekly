package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/html"
)

const (
	tableToken = "table"
	linkToken  = "a"
)

const maxArticlesPerWeek = 15

var (
	hrefAttrName  = []byte("href")
	classAttrName = []byte("class")
	tableClasses  = []byte("el-item item  ")
)

func Must(c *Crawler, err error) *Crawler {
	if err != nil {
		panic(err)
	}

	return c
}

func New(rawURL string, timeout time.Duration, retries uint) (*Crawler, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		URL:             url,
		Retries:         retries,
		RetriesInterval: time.Second * 15,
		client: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

type Crawler struct {
	Retries         uint
	RetriesInterval time.Duration
	URL             *url.URL

	tokenizer *html.Tokenizer
	client    *http.Client
}

func (c *Crawler) ParseArticles() ([]Article, error) {
	siteHTML, err := c.getHTMLStream()
	if err != nil {
		return nil, err
	}

	defer siteHTML.Close()
	c.tokenizer = html.NewTokenizer(siteHTML)
	articles := make([]Article, 0, maxArticlesPerWeek)

	// articles is contained in tables with class 'el-item item  '
	// they have same structure regardless of content table > tbody > tr > td > (content)
	for {
		ok, err := c.findTokenByAttr(classAttrName, tableClasses)
		if err != nil {
			return nil, err
		}

		if !ok {
			break
		}

		article, err := c.getArticleFromStream()
		if err != nil {
			return nil, err
		}

		if !article.IsSponsored {
			articles = append(articles, *article)
		}
	}

	return articles, nil
}

func (c *Crawler) findTokenByAttr(attrName []byte, attrValue []byte) (bool, error) {
	for c.tokenizer.Err() == nil {
		if c.tokenizer.Next() != html.StartTagToken {
			continue
		}

		for {
			attr, value, hasMore := c.tokenizer.TagAttr()
			if bytes.Equal(attr, attrName) && bytes.Equal(value, attrValue) {
				return true, nil
			}

			if !hasMore {
				break
			}
		}
	}

	if errors.Is(c.tokenizer.Err(), io.EOF) {
		return false, nil
	}

	return false, c.tokenizer.Err()
}

func (c *Crawler) getArticleFromStream() (*Article, error) {
	once := sync.Once{}
	tokens := make([]string, 0, 5)

	for c.tokenizer.Err() == nil {
		switch c.tokenizer.Next() {
		case html.StartTagToken:
			tagName, _ := c.tokenizer.TagName()
			if string(tagName) == linkToken {
				once.Do(func() {
					tokens = append(tokens, string(c.getTokensAttr(hrefAttrName)))
				})
			}

		case html.EndTagToken:
			tagName, _ := c.tokenizer.TagName()
			if string(tagName) == tableToken {
				return newArticleFromTextTokens(tokens)
			}

		case html.TextToken:
			if text := bytes.TrimSpace(c.tokenizer.Text()); len(text) > 0 {
				tokens = append(tokens, string(text))
			}
		}
	}

	if errors.Is(c.tokenizer.Err(), io.EOF) {
		return newArticleFromTextTokens(tokens)
	}

	return nil, c.tokenizer.Err()
}

func (c *Crawler) getTokensAttr(attrName []byte) []byte {
	for {
		attr, value, hasMore := c.tokenizer.TagAttr()
		if bytes.Equal(attr, attrName) {
			return value
		}

		if !hasMore {
			return nil
		}
	}
}

func (c *Crawler) getHTMLStream() (io.ReadCloser, error) {
	var (
		err error
		res *http.Response
	)

	req := &http.Request{
		Method: http.MethodGet,
		URL:    c.URL,
	}

	for i := 0; i < int(c.Retries); i++ {
		res, err = c.client.Do(req)
		if err != nil || res.StatusCode != http.StatusOK {
			err = fmt.Errorf("failed to get site HTML, status: %v, err: %w", res.StatusCode, err)
			time.Sleep(c.RetriesInterval)
			continue
		}

		return res.Body, nil
	}

	return nil, err
}
