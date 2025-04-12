package app

import "context"

type Service interface {
	Health(ctx context.Context) []byte
}
