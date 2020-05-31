package broadcast

import (
	"context"
	"github.com/mazxaxz/BroadcastMQ/cmd/config"
	"github.com/mazxaxz/BroadcastMQ/pkg/rabbitmq"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
)

type Broadcaster struct {
	Config  []config.Broadcast
	Clients map[string]rabbitmq.Client
	Logger  *logrus.Logger
}

// Initialize makes sure MQ clients are connected and source/destination resources are present
func (b *Broadcaster) Initialize(ctx context.Context, log *logrus.Logger) error {
	b.Clients = make(map[string]rabbitmq.Client, 0)
	b.Logger = log

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

func (b *Broadcaster) addMqClient(ctx context.Context, cs string) (rabbitmq.Client, error) {
	var client rabbitmq.Client
	var err error

	if _, exists := b.Clients[cs]; !exists {
		client, err = rabbitmq.NewClient(ctx, cs, b.Logger)
		if err != nil {
			return client, err
		}

		b.Clients[cs] = client
	}

	return client, nil
}

// Start begins broadcasting
func (b *Broadcaster) Start() {
	var wg sync.WaitGroup
	for _, broadcast := range b.Config {
		wg.Add(1)
		go func(bc config.Broadcast) {
			defer wg.Done()

			if src, ok := b.Clients[bc.Source.ConnectionString]; ok {
				b.Logger.Infof("Started broadcasting, '%s' === '%s' ===> *", bc.Source.Exchange, bc.Destination.BMQRoutingKey)

				if err := src.Consume(bc.Source.BMQQueueName, b.forward(bc.Destination)); err != nil {
					b.Logger.Error(err)
				}
			}
		}(broadcast)
	}
	wg.Wait()
}

func (b *Broadcaster) forward(cfg config.Destination) func(msg amqp.Delivery) {
	dest, ok := b.Clients[cfg.ConnectionString]

	return func(msg amqp.Delivery) {
		if !ok {
			return
		}

		var headers map[string]interface{}
		if cfg.PersistHeaders {
			headers = msg.Headers
		}

		if err := dest.Publish(cfg.BMQExchange, cfg.BMQRoutingKey, msg.Body, headers); err != nil {
			return
		}
	}
}
