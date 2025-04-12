package app

import (
	"context"
)

type service struct {
	logger Logger
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

func New(logger Logger) Service {
	a := &service{
		logger: logger,
	}

	return a
}

func (s *service) Health(_ context.Context) []byte {
	return []byte("OK")
}
