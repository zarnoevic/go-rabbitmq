package rabbitClient

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   string
}

func NewRabbitClient(amqpURL, queueName string) (*RabbitClient, error) {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		if err = conn.Close(); err != nil {
			return nil, err
		}
		return nil, err
	}

	_, err = channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		if err := channel.Close(); err != nil {
			return nil, err
		}
		if err = conn.Close(); err != nil {
			return nil, err
		}
		return nil, err
	}

	return &RabbitClient{
		conn:    conn,
		channel: channel,
		queue:   queueName,
	}, nil
}

func (c *RabbitClient) Publish(message string) error {
	if c.conn == nil || c.channel == nil {
		return amqp091.ErrClosed
	}

	return c.channel.PublishWithContext(context.Background(),
		"",
		c.queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (c *RabbitClient) Consume() (<-chan amqp091.Delivery, error) {
	if c.conn == nil || c.channel == nil {
		return nil, amqp091.ErrClosed
	}

	return c.channel.Consume(
		c.queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (c *RabbitClient) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
