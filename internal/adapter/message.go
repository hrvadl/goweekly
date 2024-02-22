package adapter

import (
	"github.com/hrvadl/go-weekly/internal/crawler"
)

type ArticleFormatter interface {
	FormatArticle(a crawler.Article) string
}

type WeeklySender interface {
	SendWeekly(messages []string)
}

func NewArticleSender(snd WeeklySender, fmt ArticleFormatter) *Adapter {
	return &Adapter{
		formatter: fmt,
		sender:    snd,
	}
}

type Adapter struct {
	articles  []crawler.Article
	formatter ArticleFormatter
	sender    WeeklySender
}

func (a *Adapter) SendWeekly(articles []crawler.Article) {
	msg := make([]string, 0, len(a.articles))
	for _, article := range articles {
		msg = append(msg, a.formatter.FormatArticle(article))
	}

	a.sender.SendWeekly(msg)
}
