package common

type (
	ConnectionHandler[T any] func() (T, error)
	MessageHandler[T any]    func(message T)
)
