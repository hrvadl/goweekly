package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

const tableAttrName = "table"

var (
	attrName  = []byte("class")
	attrValue = []byte("el-item item  ")
)

type Article struct {
	URL     string
	Header  string
	Content []string
}

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
	siteHTML, err := c.getSiteHTML()
	if err != nil {
		return nil, err
	}

	defer siteHTML.Close()
	c.tokenizer = html.NewTokenizer(siteHTML)
	articles := make([]Article, 0, 15)

	// articles is contained in tables with class 'el-item item  '
	// they have same structure regardless of content table > tbody > tr > td > (content)
	for {
		ok, err := c.findTableToken(attrName, attrValue)
		if err != nil {
			return nil, err
		}

		if !ok {
			break
		}

		article, err := c.processArticle()
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (c *Crawler) findTableToken(attrName []byte, attrValue []byte) (bool, error) {
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

func (c *Crawler) processArticle() (Article, error) {
	isTitle := true
	article := Article{
		Content: make([]string, 0, 5),
	}

	for c.tokenizer.Err() == nil {
		switch c.tokenizer.Next() {
		case html.EndTagToken:
			tagName, _ := c.tokenizer.TagName()
			if string(tagName) == tableAttrName {
				return article, nil
			}

		case html.TextToken:
			text := bytes.TrimSpace(c.tokenizer.Text())
			if len(text) == 0 {
				break
			}

			if isTitle {
				article.Header = string(text)
				isTitle = false
				continue
			}

			article.Content = append(article.Content, string(text))
		}
	}

	if errors.Is(c.tokenizer.Err(), io.EOF) {
		return article, nil
	}

	return article, c.tokenizer.Err()
}

func (c *Crawler) getSiteHTML() (io.ReadCloser, error) {
	var (
		err error
		res *http.Response
	)

	for i := 0; i < int(c.Retries); i++ {
		res, err = c.client.Do(&http.Request{
			Method: http.MethodGet,
			URL:    c.URL,
		})

		if err != nil || res.StatusCode != http.StatusOK {
			err = fmt.Errorf("failed to get site HTML, status: %v, err: %w", res.StatusCode, err)
			time.Sleep(c.RetriesInterval)
			continue
		}

		return res.Body, nil
	}

	return nil, err
}
