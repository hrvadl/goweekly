package adapter

import (
	"github.com/hrvadl/go-weekly/internal/crawler"
)

type ArticleFormatter interface {
	FormatArticle(a crawler.Article) string
}

func NewArticle(articles []crawler.Article, fmt ArticleFormatter) *Adapter {
	return &Adapter{
		articles:  articles,
		formatter: fmt,
	}
}

type Adapter struct {
	articles  []crawler.Article
	formatter ArticleFormatter
}

func (a *Adapter) ToMessages() []string {
	msg := make([]string, 0, len(a.articles))
	for _, ar := range a.articles {
		msg = append(msg, a.formatter.FormatArticle(ar))
	}
	return msg
}
