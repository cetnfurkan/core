package config

import "time"

type Server struct {
	Port           int
	RequestTimeout time.Duration
}
