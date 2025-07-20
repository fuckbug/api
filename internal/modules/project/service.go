package project

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	GetByID(ctx context.Context, id string) (*Entity, error)
	GetDSNByID(ctx context.Context, id string) (string, error)
	GetAll(ctx context.Context, params GetAllParams) ([]*Entity, int, error)
	Create(ctx context.Context, req *Create) (*Entity, error)
	Update(ctx context.Context, id string, req *Update) (*Entity, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo   Repository
	logger Logger
	domain string
}

func NewService(repo Repository, logger Logger, domain string) Service {
	return &service{
		repo:   repo,
		logger: logger,
		domain: domain,
	}
}

func (s *service) GetByID(ctx context.Context, id string) (*Entity, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toResponse(project), nil
}

func (s *service) GetDSNByID(ctx context.Context, id string) (string, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}

	dsn := "https://" + s.domain + "/api/ingest/" + project.ID + ":" + project.PublicKey

	return dsn, nil
}

func (s *service) GetAll(ctx context.Context, params GetAllParams) ([]*Entity, int, error) {
	projects, err := s.repo.GetAll(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*Entity, 0, len(projects))
	for _, project := range projects {
		responses = append(responses, toResponse(project))
	}
	return responses, total, nil
}

func (s *service) Create(ctx context.Context, req *Create) (*Entity, error) {
	project := &Project{
		ID:   uuid.New().String(),
		Name: req.Name,
	}

	if err := s.repo.Create(ctx, project); err != nil {
		return nil, err
	}

	return toResponse(project), nil
}

func (s *service) Update(ctx context.Context, id string, req *Update) (*Entity, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		project.Name = req.Name
	}

	if err := s.repo.Update(ctx, id, project); err != nil {
		return nil, err
	}

	return toResponse(project), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func toResponse(p *Project) *Entity {
	return &Entity{
		ID:   p.ID,
		Name: p.Name,
	}
}
