package broadcast

import (
	"github.com/mazxaxz/BroadcastMQ/cmd/config"
	"github.com/mazxaxz/BroadcastMQ/pkg/rabbitmq"
	"github.com/pkg/errors"
)

func ensureSource(c rabbitmq.Client, src config.Source) error {
	var err error
	if err = c.CreateQueue(src.BMQQueueName, true, false, true); err != nil {
		return errors.Wrap(err, "An error has occured while BMQ queue creation")
	}

	if err = c.Bind(src.BMQQueueName, src.Exchange, src.RoutingKey); err != nil {
		return errors.Wrap(err, "An error has occured while binding BMQ queue to the source Exchange")
	}

	return nil
}

func ensureDestination(c rabbitmq.Client, dest config.Destination) error {
	if err := c.CreateExchange(dest.BMQExchange, "topic", false, true); err != nil {
		return errors.Wrap(err, "An error has occured while BMQ destination Exchange creation")
	}

	for _, queue := range dest.Queues {
		if queue.EnsureExists {
			if err := c.CreateQueue(queue.Name, true, false, false); err != nil {
				return errors.Wrapf(err, "Could not create destination queue: '%s'", queue.Name)
			}
		}

		if err := c.Bind(queue.Name, dest.BMQExchange, queue.BMQBindingKey); err != nil {
			return errors.Wrap(err, "An error has occured while binding destination Queue to BMQ exchange")
		}
	}

	return nil
}