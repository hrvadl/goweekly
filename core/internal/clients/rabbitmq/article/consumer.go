package article

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

type Article struct {
	URL         string `json:"url,omitempty"`
	Header      string `json:"header,omitempty"`
	Content     string `json:"content,omitempty"`
	Author      string `json:"author,omitempty"`
	IsSponsored bool   `json:"isSponsored,omitempty"`
}

type Config struct {
	Addr string
}

func NewConsumer(log *slog.Logger, processor ArticleProccessor) *Consumer {
	return &Consumer{
		log:       log,
		processor: processor,
	}
}

type ArticleProccessor interface {
	Process(msg Article) error
}

type Consumer struct {
	conn      *amqp091.Connection
	ch        *amqp091.Channel
	log       *slog.Logger
	processor ArticleProccessor
}

func (c *Consumer) Connect(addr string) error {
	conn, err := amqp091.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to rabbitMQ: %w", err)
	}

	c.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel from the conn: %w", err)
	}

	c.ch = ch
	return nil
}

func (c *Consumer) Close() error {
	closeChanErr := c.ch.Close()
	closeConnErr := c.conn.Close()
	if closeChanErr != nil {
		return closeChanErr
	}
	return closeConnErr
}

func (c *Consumer) Consume() error {
	q, err := c.ch.QueueDeclare(
		"articles",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := c.ch.Consume(
		q.Name,
		"",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	g := new(errgroup.Group)
	go func() {
		for d := range msgs {
			g.Go(func() error {
				var article Article
				if err := json.Unmarshal(d.Body, &article); err != nil {
					c.log.Error("Failed to unmarshall message's body", "err", err)
					return err
				}
				c.log.Info("Got article from q", "article", article)
				return c.processor.Process(article)
			})
		}
	}()

	return nil
}
