package loggroup

import (
	"context"
)

type Repository interface {
	GetAll(ctx context.Context, params GetAllParams) ([]*Group, error)
	Count(ctx context.Context, params FilterParams) (int, error)
	GetByID(ctx context.Context, id string) (*Group, error)
}
