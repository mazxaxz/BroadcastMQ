package rabbitmq

import (
	"context"
	"github.com/streadway/amqp"
	"sync"
)

type channelPool struct {
	sync.RWMutex
	channels []*amqp.Channel
	poolSize int
	prefetch int
}

func newChannelPool(ctx context.Context, connection *amqp.Connection, poolSize int, prefetch int) *channelPool {
	pool := channelPool{
		channels: make([]*amqp.Channel, 0, poolSize+10),
		poolSize: poolSize,
		prefetch: prefetch,
	}

	for i := 0; i < pool.poolSize; i++ {
		ch, err := connection.Channel()
		if err != nil {
			continue
		}
		ch.Qos(pool.prefetch, 0, false)

		pool.channels = append(pool.channels, ch)
	}

	go disposePool(ctx, &pool)
	return &pool
}

func (p *channelPool) get(connection *amqp.Connection) (*amqp.Channel, error) {
	p.Lock()
	defer p.Unlock()

	if len(p.channels) == 0 {
		ch, err := connection.Channel()
		if err != nil {
			return nil, err
		}

		ch.Qos(p.prefetch, 0, false)
		return ch, nil
	}

	ch, rest := p.channels[len(p.channels)-1], p.channels[:len(p.channels)-1]
	p.channels = rest

	return ch, nil
}

func (p *channelPool) release(ch *amqp.Channel) {
	p.RLock()
	defer p.RUnlock()

	p.channels = append(p.channels, ch)
}

func disposePool(c context.Context, p *channelPool) {
	for {
		select {
		case <-c.Done():
			p.Lock()
			for _, ch := range p.channels {
				ch.Close()
			}
			p.Unlock()

			return
		}
	}
}
