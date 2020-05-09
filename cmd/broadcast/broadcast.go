package broadcast

import (
	"fmt"
	"github.com/mazxaxz/BroadcastMQ/cmd/config"
	"github.com/mazxaxz/BroadcastMQ/pkg/rabbitmq"
	"github.com/streadway/amqp"
	"sync"
)

type Broadcast struct {
	Config  []config.Broadcast
	clients map[string]*rabbitmq.Client
}

func (b *Broadcast) Initialize() error {
	b.clients = make(map[string]*rabbitmq.Client, 0)

	for _, broadcast := range b.Config {
		cs := broadcast.Source.ConnectionString
		if _, exists := b.clients[cs]; !exists {
			client := rabbitmq.Client{ConnectionString: cs}
			err := client.Connect()
			if err != nil {
				return fmt.Errorf("Could not connect: %v", err)
			}

			b.clients[cs] = &client
		}

		cs = broadcast.Destination.ConnectionString
		if _, exists := b.clients[cs]; !exists {
			client := rabbitmq.Client{ConnectionString: cs}
			err := client.Connect()
			if err != nil {
				return fmt.Errorf("Could not connect: %v", err)
			}

			b.clients[cs] = &client
		}
	}

	// TODO: bind queues

	return nil
}

func (b *Broadcast) Dispose() {
	// TODO
	// disconnect consumers
	// delete BMQ queues and exchanges
	// close channels
	// close connections
}

func (b *Broadcast) Start() {
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
				src.Consume(
					queueName,
					b.forward(bc.Destination),
					)
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

		err := dest.Publish(exchange, routingKey, msg.Body, headers)
		if err != nil {
			// TODO handle
		}
	}
}