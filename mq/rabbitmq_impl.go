package mq

import (
	"fmt"
	"sync"

	"github.com/cetnfurkan/core/config"
	"github.com/cetnfurkan/core/mq/common"
	"github.com/cetnfurkan/core/mq/consumer"
	"github.com/cetnfurkan/core/mq/producer"

	"github.com/streadway/amqp"
)

type (
	RabbitMQ struct {
		cfg        *config.MQ
		connection *amqp.Connection
		mutex      sync.Mutex

		// callbacks
		consumerMessageHandler common.MessageHandler[*amqp.Delivery]
	}

	rabbitOption func(*RabbitMQ)
)

func WithRabbitConsumerMessageHandler(handler common.MessageHandler[*amqp.Delivery]) rabbitOption {
	return func(rabbitmq *RabbitMQ) {
		rabbitmq.consumerMessageHandler = handler
	}
}

func NewRabbitMQ(cfg *config.MQ, opts ...rabbitOption) MQ {
	rabbitmq := &RabbitMQ{
		cfg: cfg,
	}

	rabbitmq.consumerMessageHandler = rabbitmq.defaultMessageHandler

	for _, opt := range opts {
		opt(rabbitmq)
	}

	return rabbitmq
}

func (rabbitMQ *RabbitMQ) Connection() (any, error) {
	var (
		err error
	)

	rabbitMQ.mutex.Lock()
	defer rabbitMQ.mutex.Unlock()

	if rabbitMQ.connection != nil && !rabbitMQ.connection.IsClosed() {
		return rabbitMQ.connection, nil
	}

	rabbitMQ.connection, err = amqp.Dial(rabbitMQ.url())
	if err != nil {
		return nil, err
	}

	return rabbitMQ.connection, nil
}

func (rabbitMQ *RabbitMQ) defaultAmqpConnection() (*amqp.Connection, error) {
	connection, err := rabbitMQ.Connection()
	if err != nil {
		return nil, err
	}

	return connection.(*amqp.Connection), nil
}

func (rabbitMQ *RabbitMQ) url() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d",
		rabbitMQ.cfg.User,
		rabbitMQ.cfg.Password,
		rabbitMQ.cfg.Host,
		rabbitMQ.cfg.Port,
	)
}

func (rabbitMQ *RabbitMQ) Consumer() consumer.Consumer {
	return consumer.NewRabbitMQConsumer(
		rabbitMQ.cfg,
		rabbitMQ.defaultAmqpConnection,
		rabbitMQ.consumerMessageHandler,
	)
}

func (rabbitMQ *RabbitMQ) Producer() producer.Producer {
	return producer.NewRabbitMQProducer(
		rabbitMQ.cfg,
		rabbitMQ.defaultAmqpConnection,
	)
}

func (rabbitMQ *RabbitMQ) ProducerWith(opts ...producer.RabbitMQProducerOption) producer.Producer {
	return producer.NewRabbitMQProducer(
		rabbitMQ.cfg,
		rabbitMQ.defaultAmqpConnection,
		opts...,
	)
}

func (rabbitMQ *RabbitMQ) defaultMessageHandler(message *amqp.Delivery) {
	fmt.Println(string(message.Body))
}
