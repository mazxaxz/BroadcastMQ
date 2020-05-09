package rabbitmq

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"github.com/streadway/amqp"
)

type Client struct {
	ConnectionString string
	connection       *amqp.Connection
}

func (c *Client) Connect() error {
	if c.connection != nil {
		return nil
	}

	conn, err := amqp.Dial(c.ConnectionString)
	if err != nil {
		return fmt.Errorf("Could not establish connection with RabbitMQ: %v", err)
	}

	c.connection = conn
	// TODO handle dispose
	return nil
}

func (c *Client) Consume(queueName string, callback func(delivery amqp.Delivery)) error {
	// TODO: create channel pool

	ch, err := c.connection.Channel()
	if err != nil {
		return fmt.Errorf("Could not establish channel connection: %v", err)
	}
	defer ch.Close()

	cid := ksuid.New()
	messages, err := ch.Consume(
		queueName,      // queue
		cid.String(),   // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,      // args
	)
	if err != nil {
		return fmt.Errorf("Could not subscribe to queue: %v", err)
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			callback(d)
		}
	}()
	<-forever

	return nil
}

func (c *Client) Publish(ex string, key string, body []byte, headers map[string]interface{}) error {
	// TODO channel pool

	ch, err := c.connection.Channel()
	if err != nil {
		return fmt.Errorf("Could not establish channel connection: %v", err)
	}
	defer ch.Close()

	err = ch.Publish(
		ex,     // exchange
		key,    // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			Body:        body,
			Headers:     headers,
		})
	if err != nil {
		return fmt.Errorf("Could not publish message: %v", err)
	}

	return nil
}