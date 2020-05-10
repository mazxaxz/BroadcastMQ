package rabbitmq

func (c *Client) CreateExchange(exchange, kind string, durable, autoDelete bool) error {
	// TODO channel pool
	ch, _ := c.connection.Channel()
	defer ch.Close()

	err := ch.ExchangeDeclare(exchange, kind, durable, autoDelete, false, false, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateQueue(queue string, durable, autoDelete, exclusive bool) error {
	// TODO channel pool
	ch, _ := c.connection.Channel()
	defer ch.Close()

	_, err := ch.QueueDeclare(queue, durable, autoDelete, exclusive, false, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Bind(queue, exchange, bindingKey string) error {
	// TODO channel pool
	ch, _ := c.connection.Channel()
	defer ch.Close()

	err := ch.QueueBind(queue, bindingKey, exchange, false, nil)
	if err != nil {
		return err
	}

	return nil
}
