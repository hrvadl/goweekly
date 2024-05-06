package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/html"
)

const (
	tableTag = "table"
	linkTag  = "a"
)

const (
	maxArticlesPerWeek = 15
	articlesURL        = "https://golangweekly.com/issues/latest"
)

var (
	hrefAttrName  = []byte("href")
	classAttrName = []byte("class")
	tableClasses  = []byte("el-item item  ")
)

func New(timeout time.Duration, retries int) *Crawler {
	return &Crawler{
		URL:             articlesURL,
		Retries:         retries,
		RetriesInterval: time.Second * 15,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

type Crawler struct {
	Retries         int
	RetriesInterval time.Duration
	URL             string

	tokenizer *html.Tokenizer
	client    *http.Client
}

/*
Fetches and parses the articles

Articles is contained in tables with class 'el-item item  '.
They have same structure regardless of content table > tbody > tr > td > (content)
*/
func (c *Crawler) ParseArticles() ([]Article, error) {
	siteHTML, err := c.getHTMLStream()
	if err != nil {
		return nil, err
	}

	defer siteHTML.Close()
	c.tokenizer = html.NewTokenizer(siteHTML)
	articles := make([]Article, 0, maxArticlesPerWeek)

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

/*
Moves the tokenizer cursor to the element

which matches the provided attribute name
and attribute class
*/
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

/*
Parses tokens from the tokenizer

Tokenizer's cursor should already be on <table> element
Exits on the EOF exception or on the closing </table> element
*/
// nolint
func (c *Crawler) getArticleFromStream() (*Article, error) {
	once := sync.Once{}
	tokens := make([]string, 2, 5)

	var isParsingTitle bool
	for c.tokenizer.Err() == nil {
		switch c.tokenizer.Next() {
		case html.StartTagToken:
			tagRaw, _ := c.tokenizer.TagName()
			tagName := string(tagRaw)
			if tagName == linkTag {
				isParsingTitle = true
				once.Do(func() {
					tokens[urlTokenIdx] = string(c.getTokensAttr(hrefAttrName))
				})
			}

		case html.EndTagToken:
			tagRaw, _ := c.tokenizer.TagName()
			tagName := string(tagRaw)
			if tagName == linkTag {
				isParsingTitle = false
			}

			if tagName == tableTag {
				return newArticleFromTextTokens(tokens)
			}

		case html.TextToken:
			rawText := bytes.TrimSpace(c.tokenizer.Text())
			if len(rawText) <= 0 {
				continue
			}

			if isParsingTitle {
				tokens[headerTokenIdx] += " " + string(rawText)
				continue
			}

			tokens = append(tokens, string(rawText))
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

/*
Fetches the given URL and return body stream

In case request is failed retries a couple times specified
in the config with the provided timeout.

User of this function is ought to close the outgoing stream.
*/
func (c *Crawler) getHTMLStream() (io.ReadCloser, error) {
	var (
		err error
		res *http.Response
	)

	req, err := http.NewRequest(http.MethodGet, c.URL, nil)
	if err != nil {
		return nil, err
	}

	for i := 0; i < c.Retries; i++ {
		res, err = c.client.Do(req)
		if err != nil {
			time.Sleep(c.RetriesInterval)
			continue
		}

		if res.StatusCode != http.StatusOK {
			err = fmt.Errorf("failed to get site HTML, status: %v, err: %w", res.StatusCode, err)
			time.Sleep(c.RetriesInterval)
			continue
		}

		return res.Body, nil
	}

	return nil, err
}
