package main

import (
	"fmt"
	"time"

	"github.com/hrvadl/go-weekly/internal/crawler"
)

const (
	articlesURL = "https://golangweekly.com/issues/latest"
	retries     = 3
	timeout     = 30 * time.Second
)

func main() {
	crawler := crawler.Must(crawler.New(articlesURL, timeout, retries))

	articles, err := crawler.ParseArticles()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(articles)
}
