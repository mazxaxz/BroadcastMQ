package rabbitmq

import (
	"context"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Client struct {
	connection *amqp.Connection
	logger     *logrus.Logger
	pool       *channelPool
}

// NewClient creates new RabbitMQ client
func NewClient(ctx context.Context, uri string, log *logrus.Logger) (Client, error) {
	client := Client{
		connection: nil,
		logger:     log,
	}

	var err error
	client.connection, err = amqp.Dial(uri)
	if err != nil {
		return client, errors.Wrap(err, "Could not establish connection with RabbitMQ")
	}

	client.pool = newChannelPool(ctx, client.connection, 10, 10)

	go disposeClient(ctx, client)
	return client, nil
}

// Consume subscribes to given queue
func (c *Client) Consume(queueName string, callback func(delivery amqp.Delivery)) error {
	ch, err := c.pool.get(c.connection)
	if err != nil {
		return errors.Wrap(err, "Could not establish channel connection")
	}

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
		return errors.Wrap(err, "Could not subscribe to queue")
	}

	go func() {
		for d := range messages {
			go callback(d)
		}
	}()

	return nil
}

// Publish sends payload to given exchange
func (c *Client) Publish(ex, key string, body []byte, headers map[string]interface{}) error {
	ch, err := c.pool.get(c.connection)
	if err != nil {
		return errors.Wrap(err, "Could not establish channel connection")
	}
	defer c.pool.release(ch)

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
		return errors.Wrap(err, "Could not publish message")
	}

	return nil
}

func disposeClient(ctx context.Context, c Client) {
	for {
		select {
		case <-ctx.Done():
			if e := c.connection.Close(); e != nil {
				c.logger.Error(e)
			}
			return
		case <-c.connection.NotifyClose(make(chan *amqp.Error)):
			c.logger.Infof("Closed because of AMQP.")
			return
		}
	}
}
