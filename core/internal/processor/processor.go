package processor

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/hrvadl/goweekly/core/internal/clients/rabbitmq/article"
)

//go:generate mockgen -destination=./mocks/mock_translator.go -package=mocks -source ./processor.go Translator
type Translator interface {
	Translate(ctx context.Context, msg string) (string, error)
}

//go:generate mockgen -destination=./mocks/mock_sender.go -package=mocks -source ./processor.go Sender
type Sender interface {
	Send(ctx context.Context, msg string) error
}

//go:generate mockgen -destination=./mocks/mock_formatter.go -package=mocks -source ./processor.go ArticleFormatter
type ArticleFormatter interface {
	FormatArticle(a article.Article) string
}

func New(fmter ArticleFormatter, sender Sender, translator Translator) *Processor {
	return &Processor{
		formatter:  fmter,
		sender:     sender,
		translator: translator,
	}
}

type Processor struct {
	sender     Sender
	translator Translator
	formatter  ArticleFormatter
}

func (p *Processor) Process(a article.Article) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	translated, err := p.translator.Translate(ctx, a.Content)
	if err != nil {
		return fmt.Errorf("failed to translate article: %w", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	a.Content = translated
	if err := p.sender.Send(ctx, p.formatter.FormatArticle(a)); err != nil {
		return fmt.Errorf("failed to translate article: %w", err)
	}

	return nil
}
