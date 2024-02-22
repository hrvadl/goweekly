package app

import "github.com/hrvadl/go-weekly/internal/crawler"

type Writer interface {
	GetArticles() ([]crawler.Article, error)
}
