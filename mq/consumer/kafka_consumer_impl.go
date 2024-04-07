package consumer

import (
	"context"
	"fmt"

	"github.com/cetnfurkan/core/config"
	"github.com/cetnfurkan/core/mq/common"
	"github.com/segmentio/kafka-go"
)

type (
	KafkaConsumer struct {
		cfg            *config.MQ
		dialerHandler  common.ConnectionHandler[*kafka.Dialer]
		messageHandler common.MessageHandler[*kafka.Message]
	}
)

func NewKafkaConsumer(
	cfg *config.MQ,
	dialerHandler common.ConnectionHandler[*kafka.Dialer],
	messageHandler common.MessageHandler[*kafka.Message],
) Consumer {

	return &KafkaConsumer{
		cfg:            cfg,
		dialerHandler:  dialerHandler,
		messageHandler: messageHandler,
	}
}

func (consumer *KafkaConsumer) Consume(topic string) error {
	dialer, err := consumer.dialerHandler()
	if err != nil {
		return err
	}

	broker := fmt.Sprintf("%s:%d", consumer.cfg.Host, consumer.cfg.Port)
	readerConfig := kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: fmt.Sprintf("g_%s", topic),
		Dialer:  dialer,
	}

	reader := kafka.NewReader(readerConfig)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		consumer.messageHandler(&msg)
	}
}
