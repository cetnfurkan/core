package config

import "time"

type Service struct {
	Address        string
	RequestTimeout time.Duration
}
