package producer

type (
	Producer interface {
		Produce(queueName string, message []byte) error
	}
)
