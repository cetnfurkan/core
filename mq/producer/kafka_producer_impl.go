package producer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cetnfurkan/core/config"
	"github.com/cetnfurkan/core/mq/common"
	"github.com/segmentio/kafka-go"
)

type (
	KafkaProducer struct {
		cfg           *config.MQ
		dialerHandler common.ConnectionHandler[*kafka.Dialer]
		mutex         sync.Mutex
	}
)

func NewKafkaProducer(
	cfg *config.MQ,
	dialerHandler common.ConnectionHandler[*kafka.Dialer],
) Producer {

	return &KafkaProducer{
		cfg:           cfg,
		dialerHandler: dialerHandler,
	}
}

func (producer *KafkaProducer) Produce(topic string, message []byte) error {

	producer.mutex.Lock()
	defer producer.mutex.Unlock()

	dialer, err := producer.dialerHandler()
	if err != nil {
		return err
	}

	broker := fmt.Sprintf("%s:%d", producer.cfg.Host, producer.cfg.Port)
	writerConfig := kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   topic,
		Dialer:  dialer,
	}

	writer := kafka.NewWriter(writerConfig)
	writer.AllowAutoTopicCreation = true
	defer writer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return writer.WriteMessages(ctx, kafka.Message{
		Value: message,
	})
}
