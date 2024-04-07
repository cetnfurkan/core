package mq

import (
	"github.com/cetnfurkan/core/mq/consumer"
	"github.com/cetnfurkan/core/mq/producer"
)

type (
	MQ interface {
		Connection() (any, error)
		Consumer() consumer.Consumer
		Producer() producer.Producer
	}
)
