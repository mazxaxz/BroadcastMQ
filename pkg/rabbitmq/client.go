package rabbitmq

import (
	"context"
	"fmt"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Client struct {
	connection *amqp.Connection
	done       chan error
	logger     *logrus.Logger
}

func NewClient(ctx context.Context, uri string, log *logrus.Logger) (*Client, error) {
	c := Client{
		connection: nil,
		done:       make(chan error),
		logger:     log,
	}

	var err error
	c.connection, err = amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("Could not establish connection with RabbitMQ: %v", err)
	}

	go func(client *Client) {
		for {
			select {
			case _ <- ctx.Done():
				if e := client.shutdown(); e != nil {
					client.logger.Error(e)
				}
				return
			case _ <- c.connection.NotifyClose(make(chan *amqp.Error)):
				client.logger.Infof("Closed because of AMQP.")
				return
			}
		}
	}(&c)

	return &c, nil
}

func (c *Client) shutdown() error {
	c.logger.Infof("Shutting down connection")

	if err := c.connection.Close(); err != nil {
		return fmt.Errorf("Could not close connection: %v", err)
	}

	return <-c.done
}

func (c *Client) Consume(ctx context.Context, queueName string, callback func(delivery amqp.Delivery)) error {
	// TODO: create channel pool

	ch, err := c.connection.Channel()
	if err != nil {
		return fmt.Errorf("Could not establish channel connection: %v", err)
	}
	defer ch.Close()

	consumerName := ksuid.New().String() + ".BroadcastMQ"
	messages, err := ch.Consume(
		queueName,    // queue
		consumerName, // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("Could not subscribe to queue: %v", err)
	}

	go handle(ctx, messages, c.done, callback)

	return nil
}

func (c *Client) Publish(ex, key string, body []byte, headers map[string]interface{}) error {
	// TODO channel pool

	ch, err := c.connection.Channel()
	if err != nil {
		return fmt.Errorf("Could not establish channel connection: %v", err)
	}
	defer ch.Close()

	err = ch.Publish(
		ex,    // exchange
		key,   // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Body:    body,
			Headers: headers,
		})
	if err != nil {
		return fmt.Errorf("Could not publish message: %v", err)
	}

	return nil
}

func handle(ctx context.Context, messages <-chan amqp.Delivery, done chan error, cb func(delivery amqp.Delivery)) {
	for d := range messages {
		select {
		case _ <- ctx.Done():
			done <- nil
			return
		default:
			cb(d)
		}
	}
	done <- nil
}
