package mq

import (
	"fmt"
	"time"

	"github.com/cetnfurkan/core/config"
	"github.com/cetnfurkan/core/mq/common"
	"github.com/cetnfurkan/core/mq/consumer"
	"github.com/cetnfurkan/core/mq/producer"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type (
	Kafka struct {
		cfg      *config.MQ
		consumer consumer.Consumer
		producer producer.Producer

		// callbacks
		consumerMessageHandler common.MessageHandler[*kafka.Message]
		kafkaDialerHandler     common.ConnectionHandler[*kafka.Dialer]
	}

	kafkaOption func(*Kafka)
)

func WithKafkaConsumerMessageHandler(handler common.MessageHandler[*kafka.Message]) kafkaOption {
	return func(kafka *Kafka) {
		kafka.consumerMessageHandler = handler
	}
}

func WithKafkaDialerHandler(handler common.ConnectionHandler[*kafka.Dialer]) kafkaOption {
	return func(kafka *Kafka) {
		kafka.kafkaDialerHandler = handler
	}
}

func NewKafka(cfg *config.MQ, opts ...kafkaOption) MQ {
	kafkamq := &Kafka{
		cfg: cfg,
	}

	kafkamq.consumerMessageHandler = kafkamq.defaultMessageHandler
	kafkamq.kafkaDialerHandler = kafkamq.defaultKafkaDialer

	for _, opt := range opts {
		opt(kafkamq)
	}

	kafkamq.consumer = consumer.NewKafkaConsumer(
		cfg,
		kafkamq.kafkaDialerHandler,
		kafkamq.consumerMessageHandler,
	)

	kafkamq.producer = producer.NewKafkaProducer(
		cfg,
		kafkamq.kafkaDialerHandler,
	)

	return kafkamq
}

func (kafkamq *Kafka) Connection() (any, error) {
	return kafkamq.defaultKafkaDialer()
}

func (kafkamq *Kafka) defaultKafkaDialer() (*kafka.Dialer, error) {
	return &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		SASLMechanism: plain.Mechanism{
			Username: kafkamq.cfg.User,
			Password: kafkamq.cfg.Password,
		},
	}, nil
}

func (kafkamq *Kafka) Consumer() consumer.Consumer {
	return kafkamq.consumer
}

func (kafkamq *Kafka) defaultMessageHandler(message *kafka.Message) {
	fmt.Println(string(message.Value))
}

func (kafkamq *Kafka) Producer() producer.Producer {
	return kafkamq.producer
}
