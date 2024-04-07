package database

// Database is an interface for databases like postgres, mysql, etc.
type Database interface {
	// Get returns the database client instance.
	Get() any

	// UnmarshalExtra unmarshals extra config data.
	// It will panic if it fails to unmarshal.
	UnmarshalExtra()
}

type option[T any] func(*T) error

func WithCallback[T any](callback func(*T) error) option[T] {
	return func(client *T) error {
		return callback(client)
	}
}
