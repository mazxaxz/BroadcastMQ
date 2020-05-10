package broadcast

import (
	"context"
	"fmt"
	"github.com/mazxaxz/BroadcastMQ/cmd/config"
	"github.com/mazxaxz/BroadcastMQ/pkg/rabbitmq"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
)

type Broadcast struct {
	Config  []config.Broadcast
	clients map[string]*rabbitmq.Client
	logger  *logrus.Logger
}

func (b *Broadcast) Initialize(ctx context.Context, log *logrus.Logger) error {
	b.clients = make(map[string]*rabbitmq.Client, 0)
	b.logger = log

	for _, bc := range b.Config {
		client, err := b.addMqClient(ctx, bc.Source.ConnectionString)
		if err != nil {
			return err
		}

		// TODO
		// bc.Source.BmqQueueName ?? DefaultBmqQueueName
		// bc.Destination.BmqExchange ?? DefaultBmqExchange
		// queue.BmqBindingKey ?? DefaultBmqBindingKey

		client.CreateQueue(bc.Source.BmqQueueName, false, true, true)                 // TODO err
		client.Bind(bc.Source.BmqQueueName, bc.Source.Exchange, bc.Source.RoutingKey) // TODO err
		// TODO what if exchange does not exist

		client, err = b.addMqClient(ctx, bc.Destination.ConnectionString)
		if err != nil {
			return err
		}

		client.CreateExchange(bc.Destination.BmqExchange, "topic", false, true) // TODO err
		for _, queue := range bc.Destination.Queues {
			// TODO what if queue does not exist
			client.Bind(queue.Name, bc.Destination.BmqExchange, queue.BmqBindingKey) // TODO err
		}
	}

	return nil
}

func (b *Broadcast) addMqClient(ctx context.Context, cs string) (*rabbitmq.Client, error) {
	var client *rabbitmq.Client
	var err error

	if _, exists := b.clients[cs]; !exists {
		client, err = rabbitmq.NewClient(ctx, cs, b.logger)
		if err != nil {
			return nil, fmt.Errorf("Could not connect: %v", err)
		}

		b.clients[cs] = client
	}

	return client, nil
}

func (b *Broadcast) Start(ctx context.Context) {
	var wg sync.WaitGroup
	for _, broadcast := range b.Config {
		wg.Add(1)
		go func(bc *config.Broadcast) {
			defer wg.Done()

			queueName := bc.Source.BmqQueueName
			if queueName == "" {
				queueName = DefaultBqQueue
			}

			if src, ok := b.clients[bc.Source.ConnectionString]; ok {
				err := src.Consume(ctx, queueName, b.forward(bc.Destination))
				if err != nil {
					b.logger.Error(err)
				}
			}
		}(&broadcast)
	}
	wg.Wait()
}

func (b *Broadcast) forward(cfg config.Destination) func(msg amqp.Delivery) {
	dest, ok := b.clients[cfg.ConnectionString]

	exchange := cfg.BmqExchange
	if exchange == "" {
		exchange = DefaultBqExchange
	}

	routingKey := cfg.BmqRoutingKey
	if routingKey == "" {
		routingKey = DefaultBqRoutingKey
	}

	return func(msg amqp.Delivery) {
		if !ok {
			return
		}

		var headers map[string]interface{}
		if cfg.PersistHeaders {
			headers = msg.Headers
		}

		if err := dest.Publish(exchange, routingKey, msg.Body, headers); err != nil {
			// TODO handle
		}
	}
}
