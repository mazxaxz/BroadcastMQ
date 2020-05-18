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
		if err = ensureSource(client, bc.Source); err != nil {
			return err
		}

		client, err = b.addMqClient(ctx, bc.Destination.ConnectionString)
		if err != nil {
			return err
		}
		if err = ensureDestination(client, bc.Destination); err != nil {
			return err
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

func ensureSource(c *rabbitmq.Client, src config.Source) error {
	var err error
	if err = c.CreateQueue(src.BmqQueueName, true, false, true); err != nil {
		return fmt.Errorf("An error has occured while BMQ queue creation: %v", err)
	}

	if err = c.Bind(src.BmqQueueName, src.Exchange, src.RoutingKey); err != nil {
		return fmt.Errorf("An error has occured while binding BMQ queue to the source Exchange: %v", err)
	}

	return nil
}

func ensureDestination(c *rabbitmq.Client, dest config.Destination) error {
	if err := c.CreateExchange(dest.BmqExchange, "topic", false, true); err != nil {
		return fmt.Errorf("An error has occured while BMQ destination Exchange creation: %v", err)
	}

	for _, queue := range dest.Queues {
		if queue.EnsureExists {
			if err := c.CreateQueue(queue.Name, true, false, false); err != nil {
				return fmt.Errorf("Could not create destination queue: '%s'. Error: %v", queue.Name, err)
			}
		}

		if err := c.Bind(queue.Name, dest.BmqExchange, queue.BmqBindingKey); err != nil {
			return fmt.Errorf("An error has occured while binding destination Queue to BMQ exchange: %v", err)
		}
	}

	return nil
}

func (b *Broadcast) Start() {
	var wg sync.WaitGroup
	for _, broadcast := range b.Config {
		wg.Add(1)
		go func(bc config.Broadcast) {
			defer wg.Done()

			if src, ok := b.clients[bc.Source.ConnectionString]; ok {
				b.logger.Infof("Started broadcasting, src: '%s', RK: '%s'", bc.Source.Exchange, bc.Destination.BmqRoutingKey)

				err := src.Consume(bc.Source.BmqQueueName, b.forward(bc.Destination))
				if err != nil {
					b.logger.Error(err)
				}
			}
		}(broadcast)
	}
	wg.Wait()
}

func (b *Broadcast) forward(cfg config.Destination) func(msg amqp.Delivery) {
	dest, ok := b.clients[cfg.ConnectionString]

	return func(msg amqp.Delivery) {
		if !ok {
			return
		}

		var headers map[string]interface{}
		if cfg.PersistHeaders {
			headers = msg.Headers
		}

		if err := dest.Publish(cfg.BmqExchange, cfg.BmqRoutingKey, msg.Body, headers); err != nil {
			return
		}
	}
}
