package errors

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/fuckbug/api/pkg/pointers"
	"github.com/google/uuid"
)

type Service interface {
	GetByID(ctx context.Context, id string) (*Entity, error)
	GetAll(ctx context.Context, params GetAllParams) ([]*Entity, int, error)
	Create(ctx context.Context, req *Create) (*Entity, error)
	Update(ctx context.Context, id string, req *Update) (*Entity, error)
	Delete(ctx context.Context, id string) error
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

func (s *service) Create(ctx context.Context, req *Create) (*Entity, error) {
	entity := &Error{
		ID:          uuid.New().String(),
		ProjectID:   req.ProjectID,
		Message:     req.Message,
		Stacktrace:  req.Stacktrace,
		File:        req.File,
		Line:        req.Line,
		Context:     req.Context,
		Ip:          req.Ip,
		Url:         req.Url,
		Method:      req.Method,
		Headers:     req.Headers,
		QueryParams: req.QueryParams,
		BodyParams:  req.BodyParams,
		Cookies:     req.Cookies,
		Session:     req.Session,
		Files:       req.Files,
		Env:         req.Env,
		Time:        req.Time,
	}

	entity.Fingerprint = generateFingerprint(entity)

	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}

	return toResponse(entity), nil
}

func (s *service) Update(ctx context.Context, id string, req *Update) (*Entity, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Message != "" {
		entity.Message = req.Message
	}
	if req.Stacktrace != "" {
		entity.Stacktrace = req.Stacktrace
	}
	if req.File != "" {
		entity.File = req.File
	}
	if req.Line != 0 {
		entity.Line = req.Line
	}
	if pointers.DerefString(req.Context) != "" {
		entity.Context = req.Context
	}

	entity.Fingerprint = generateFingerprint(entity)

	if err := s.repo.Update(ctx, id, entity); err != nil {
		return nil, err
	}

	return toResponse(entity), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func generateFingerprint(e *Error) string {
	cleanMsg := regexp.MustCompile(`\d+|0x[0-9a-f]+`).ReplaceAllString(e.Message, "*")

	data := fmt.Sprintf(
		"%s:%s:%d",
		cleanMsg,
		e.File,
		e.Line,
	)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func toResponse(e *Error) *Entity {
	return &Entity{
		ID:          e.ID,
		Message:     e.Message,
		Stacktrace:  e.Stacktrace,
		File:        e.File,
		Line:        e.Line,
		Context:     e.Context,
		Ip:          e.Ip,
		Url:         e.Url,
		Method:      e.Method,
		Headers:     e.Headers,
		QueryParams: e.QueryParams,
		BodyParams:  e.BodyParams,
		Cookies:     e.Cookies,
		Session:     e.Session,
		Files:       e.Files,
		Env:         e.Env,
		Time:        e.Time,
	}
}
