package formatter

import (
	"fmt"
	"strings"

	"github.com/hrvadl/go-weekly/internal/crawler"
)

const MarkdownType = "Markdown"

func NewMarkdown() *Markdown {
	return &Markdown{}
}

type Markdown struct{}

func (v Markdown) FormatArticles(articles []crawler.Article) []string {
	formatted := make([]string, 0, len(articles))
	for _, a := range articles {
		formatted = append(formatted, v.FormatArticle(a))
	}
	return formatted
}

func (v Markdown) FormatType() string {
	return MarkdownType
}

func (v Markdown) FormatArticle(a crawler.Article) string {
	var builder strings.Builder
	builder.WriteString(v.FormatTitle(a.Header))
	builder.WriteString(v.FormatContent(a.Content))
	builder.WriteString(v.FormatAuthor(a.Author))
	builder.WriteString(v.FormatURL(a.URL))

	return builder.String()
}

func (v Markdown) FormatTitle(title string) string {
	return fmt.Sprintf("*%v*\n", title)
}

func (v Markdown) FormatContent(content string) string {
	return fmt.Sprintf("\n%v\n", content)
}

func (v Markdown) FormatAuthor(author string) string {
	return fmt.Sprintf("\nАвтор: *%v*\n", author)
}

func (v Markdown) FormatURL(url string) string {
	return fmt.Sprintf("Читати повну статтю: %v\n", url)
}
