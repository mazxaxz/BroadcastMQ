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
	logger     *logrus.Logger
}

func NewClient(ctx context.Context, uri string, log *logrus.Logger) (*Client, error) {
	client := Client{
		connection: nil,
		logger:     log,
	}

	var err error
	client.connection, err = amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("Could not establish connection with RabbitMQ: %v", err)
	}

	go func(c *Client) {
		for {
			select {
			case <- ctx.Done():
				if e := c.connection.Close(); e != nil {
					c.logger.Error(e)
				}
				return
			case <- c.connection.NotifyClose(make(chan *amqp.Error)):
				c.logger.Infof("Closed because of AMQP.")
				return
			}
		}
	}(&client)

	return &client, nil
}

func (c *Client) Consume(queueName string, callback func(delivery amqp.Delivery)) error {
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

	go func() {
		for d := range messages {
			go callback(d)
		}
	}()

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
