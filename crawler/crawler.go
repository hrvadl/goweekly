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

func getSiteHTML(url string) ([]byte, error) {
	log.SetPrefix("network GetSiteHtml: ")
	log.SetFlags(0)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	return body, nil
}

func CrawlSite(url string) ([]Article, error) {
	siteHTML, err := getSiteHTML(url)
	if err != nil {
		return nil, err
	}
	// create new reader to close the html body,even if copying is done
	tokenizer := html.NewTokenizer(bytes.NewReader(siteHTML))

	var articles = make([]Article, 15)
	currentArticle := 0
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

		fmt.Printf("article %v\n", currentArticle+1)
		article, err := processArticle(tokenizer)
		if err != nil {
			return nil, err
		}

		if currentArticle < len(articles) {
			articles[currentArticle] = article
		} else {
			articles = append(articles, article)
		}
		currentArticle++
	}
	return articles, nil
}

func findTokenByAttribute(tokenizer *html.Tokenizer, attrName string, attrValue string) (bool, error) {
	attrNameBytes := []byte(attrName)
	attrValueBytes := []byte(attrValue)
	attrNameLength := len(attrNameBytes)
	attrValueLength := len(attrValueBytes)
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
				if len(attr) != attrNameLength || len(value) != attrValueLength {
					if !more {
						break
					}
					continue
				}

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
	currentContent := 0
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
			} else {
				if article.Content == nil {
					currentContent = 0
					article.Content = make([]string, 5)
				}
				if currentContent < len(article.Content) {
					article.Content[currentContent] = string(text)
				} else {
					article.Content = append(article.Content, string(text))
				}
				currentContent++
			}
		}
	}
}
