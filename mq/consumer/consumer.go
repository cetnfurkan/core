package consumer

type (
	Consumer interface {
		Consume(queueName string) error
	}
)
