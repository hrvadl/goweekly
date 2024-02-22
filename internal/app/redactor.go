package app

import "github.com/hrvadl/go-weekly/internal/crawler"

type Redactor interface {
	Review([]crawler.Article) ([]crawler.Article, error)
}
