package errors

import "errors"

var (
	ErrMQConnectionIsClosed = errors.New("mq connection is closed")
)
