package producer

import (
	"sync"

	"github.com/cetnfurkan/core/config"
	"github.com/cetnfurkan/core/mq/common"

	"github.com/streadway/amqp"
)

type (
	RabbitMQProducer struct {
		args              amqp.Table
		cfg               *config.MQ
		connectionHandler common.ConnectionHandler[*amqp.Connection]
		expiration        string
		mutex             sync.Mutex
	}

	RabbitMQProducerOption func(*RabbitMQProducer)
)

func WithRabbitMQArgs(args amqp.Table) RabbitMQProducerOption {
	return func(producer *RabbitMQProducer) {
		producer.args = args
	}
}

func WithRabbitMQExpiration(expiration string) RabbitMQProducerOption {
	return func(producer *RabbitMQProducer) {
		producer.expiration = expiration
	}
}

func NewRabbitMQProducer(
	cfg *config.MQ,
	connectionHandler common.ConnectionHandler[*amqp.Connection],
	opts ...RabbitMQProducerOption,
) Producer {

	producer := &RabbitMQProducer{
		cfg:               cfg,
		connectionHandler: connectionHandler,
	}

	for _, opt := range opts {
		opt(producer)
	}

	return producer
}

func (producer *RabbitMQProducer) Produce(queueName string, message []byte) error {

	producer.mutex.Lock()
	defer producer.mutex.Unlock()

	connection, err := producer.connectionHandler()
	if err != nil {
		return err
	}

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare(
		queueName,     // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		producer.args, // arguments
	)
	if err != nil {
		return err
	}

	err = channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
			Expiration:  producer.expiration,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
