package article

import (
	"context"
	"encoding/json"
	"fmt"
	"mime"

	"github.com/rabbitmq/amqp091-go"

	"github.com/hrvadl/goweekly/crawler/internal/crawler"
)

type Config struct {
	Addr string
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

type Publisher struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
	q    amqp091.Queue
}

func (c *Publisher) Connect(addr string) error {
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

	c.q = q
	return nil
}

func (c *Publisher) Close() error {
	closeChanErr := c.ch.Close()
	closeConnErr := c.conn.Close()
	if closeChanErr != nil {
		return closeChanErr
	}
	return closeConnErr
}

func (c *Publisher) Publish(ctx context.Context, article crawler.Article) error {
	body, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("failed to construct message body: %w", err)
	}

	return c.ch.PublishWithContext(
		ctx,
		"",
		c.q.Name,
		false,
		false,
		amqp091.Publishing{ContentType: mime.TypeByExtension("json"), Body: body},
	)
}
