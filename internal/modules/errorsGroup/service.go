package errorsgroup

import "context"

type Service interface {
	GetByID(ctx context.Context, id string) (*Entity, error)
	GetAll(ctx context.Context, params GetAllParams) ([]*Entity, int, error)
}

type service struct {
	repo   Repository
	logger Logger
}

func NewService(repo Repository, logger Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) GetByID(ctx context.Context, id string) (*Entity, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toResponse(entity), nil
}

func (s *service) GetAll(ctx context.Context, params GetAllParams) ([]*Entity, int, error) {
	entities, err := s.repo.GetAll(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx, params.FilterParams)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*Entity, 0, len(entities))
	for _, entity := range entities {
		responses = append(responses, toResponse(entity))
	}
	return responses, total, nil
}

func toResponse(g *Group) *Entity {
	return &Entity{
		ID:          g.ID,
		Message:     g.Message,
		File:        g.File,
		Line:        g.Line,
		FirstSeenAt: g.FirstSeenAt,
		LastSeenAt:  g.LastSeenAt,
		Counter:     g.Counter,
	}
}
