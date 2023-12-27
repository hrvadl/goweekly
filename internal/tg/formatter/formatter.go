package formatter

import (
	"fmt"
	"strings"

	"github.com/hrvadl/go-weekly/internal/crawler"
)

func NewMarkdownV2() *MarkdownV2 {
	return &MarkdownV2{}
}

type MarkdownV2 struct {
}

func (v MarkdownV2) FormatArticles(articles []crawler.Article) []string {
	formatted := make([]string, 0, len(articles))
	for _, a := range articles {
		formatted = append(formatted, v.FormatArticle(a))
	}
	return formatted
}

func (v MarkdownV2) FormatArticle(a crawler.Article) string {
	var builder strings.Builder
	builder.WriteString(v.FormatTitle(a.Header))
	builder.WriteString(v.FormatContent(a.Content))
	builder.WriteString(v.FormatAuthor(a.Author))
	builder.WriteString(v.FormatURL(a.URL))

	return builder.String()
}

func (v MarkdownV2) FormatTitle(title string) string {
	return fmt.Sprintf("*%v*\n", title)
}

func (v MarkdownV2) FormatContent(content string) string {
	return fmt.Sprintf("\n%v\n", content)
}

func (v MarkdownV2) FormatAuthor(author string) string {
	return fmt.Sprintf("\nАвтор: *%v*\n", author)
}

func (v MarkdownV2) FormatURL(url string) string {
	return fmt.Sprintf("Читати повну статтю: %v\n", url)
}
