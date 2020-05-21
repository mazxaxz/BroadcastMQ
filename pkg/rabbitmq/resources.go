package rabbitmq

// CreateExchange makes sure that given exchange is created
func (c *Client) CreateExchange(exchange, kind string, durable, autoDelete bool) error {
	ch, err := c.pool.get(c.connection)
	if err != nil {
		return err
	}
	defer c.pool.release(ch)

	err = ch.ExchangeDeclare(exchange, kind, durable, autoDelete, false, false, nil)
	if err != nil {
		return err
	}

	return nil
}

// CreateQueue makes sure that given queue is created
func (c *Client) CreateQueue(queue string, durable, autoDelete, exclusive bool) error {
	ch, err := c.pool.get(c.connection)
	if err != nil {
		return err
	}
	defer c.pool.release(ch)

	if _, err = ch.QueueDeclare(queue, durable, autoDelete, exclusive, false, nil); err != nil {
		return err
	}

	return nil
}

// Bind makes sure that queue and exchange are linked
func (c *Client) Bind(queue, exchange, bindingKey string) error {
	ch, err := c.pool.get(c.connection)
	if err != nil {
		return err
	}
	defer c.pool.release(ch)

	if err = ch.QueueBind(queue, bindingKey, exchange, false, nil); err != nil {
		return err
	}

	return nil
}
