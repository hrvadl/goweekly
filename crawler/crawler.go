package crawler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Article struct {
	Url     string
	Header  string
	Content []string
}

func getSiteHTML(url string) (io.ReadCloser, error) {
	log.SetPrefix("network GetSiteHtml: ")
	log.SetFlags(0)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	return resp.Body, nil
}

func GetArticles(url string) ([]Article, error) {
	siteHTML, err := getSiteHTML(url)
	defer siteHTML.Close()
	if err != nil {
		return nil, err
	}
	// create new reader to close the html body,even if copying is done
	tokenizer := html.NewTokenizer(siteHTML)

	articles := make([]Article, 15)
	for {
		// articles is contained in tables with class 'el-item item  '
		// they have same structure regardless of content table > tbody > tr > td > (content)
		found, err := findTokenByAttribute(tokenizer, "class", "el-item item  ")
		if err != nil {
			return nil, err
		}

		if !found {
			break
		}

		article, err := processArticle(tokenizer)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}
	return articles, nil
}

func findTokenByAttribute(tokenizer *html.Tokenizer, attrName string, attrValue string) (bool, error) {
	attrNameBytes := []byte(attrName)
	attrValueBytes := []byte(attrValue)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return false, nil
			}
			return false, tokenizer.Err()
		case html.StartTagToken:
			for {
				attr, value, more := tokenizer.TagAttr()

				if bytes.Equal(attr, attrNameBytes) && bytes.Equal(value, attrValueBytes) {
					return true, nil
				}

				if !more {
					break
				}
			}
		}
	}
}

func processArticle(tokenizer *html.Tokenizer) (Article, error) {
	article := Article{
		Content: make([]string, 5),
	}
	title := true
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return article, nil
			}
			return article, tokenizer.Err()
		case html.EndTagToken:
			tagName, _ := tokenizer.TagName()

			fmt.Printf("end %s\n", tagName)
			if string(tagName) == "table" {
				return article, nil
			}
		case html.TextToken:
			text := strings.TrimSpace(string(tokenizer.Text()))
			if text == "" {
				break
			}

			if title {
				article.Header = text
				fmt.Printf("header: %v\n", text)
				title = false
				continue
			}

			article.Content = append(article.Content, string(text))
		}
	}
}
