package consumer

import (
	"github.com/cetnfurkan/core/config"
	"github.com/cetnfurkan/core/errors"
	"github.com/cetnfurkan/core/mq/common"

	"github.com/streadway/amqp"
)

type (
	RabbitMQConsumer struct {
		cfg               *config.MQ
		connectionHandler common.ConnectionHandler[*amqp.Connection]
		messageHandler    common.MessageHandler[*amqp.Delivery]
	}
)

func NewRabbitMQConsumer(
	cfg *config.MQ,
	connectionHandler common.ConnectionHandler[*amqp.Connection],
	messageHandler common.MessageHandler[*amqp.Delivery],
) Consumer {

	return &RabbitMQConsumer{
		cfg:               cfg,
		connectionHandler: connectionHandler,
		messageHandler:    messageHandler,
	}
}

func (consumer *RabbitMQConsumer) Consume(queueName string) error {
	connection, err := consumer.connectionHandler()
	if err != nil {
		return err
	}

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	queue, err := channel.QueueDeclare(
		queueName,    // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		amqp.Table{}, // arguments
	)
	if err != nil {
		return err
	}

	messages, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-channel.NotifyClose(make(chan *amqp.Error)):
			return errors.ErrMQConnectionIsClosed

		case message := <-messages:
			consumer.messageHandler(&message)
		}
	}
}
