package log

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/fuckbug/api/pkg/utils"
	"github.com/google/uuid"
)

var ErrInvalidLogLevel = errors.New("invalid log level")

type Service interface {
	GetByID(ctx context.Context, id string) (*Entity, error)
	GetAll(ctx context.Context, params GetAllParams) ([]*Entity, int, error)
	GetStats(ctx context.Context, projectID string, fingerprint string) (*Stats, error)
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
	log, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toResponse(log), nil
}

func (s *service) GetAll(ctx context.Context, params GetAllParams) ([]*Entity, int, error) {
	logs, err := s.repo.GetAll(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx, params.FilterParams)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*Entity, 0, len(logs))
	for _, log := range logs {
		responses = append(responses, toResponse(log))
	}

	return responses, total, nil
}

func (s *service) GetStats(ctx context.Context, projectID string, fingerprint string) (*Stats, error) {
	stats, err := s.repo.GetStats(ctx, projectID, fingerprint)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *service) Create(ctx context.Context, req *Create) (*Entity, error) {
	if !isValidLogLevel(req.Level) {
		return nil, ErrInvalidLogLevel
	}

	log := &Log{
		ID:        uuid.New().String(),
		ProjectID: req.ProjectID,
		Level:     Level(req.Level),
		Message:   req.Message,
		Context:   req.Context,
		Time:      req.Time,
	}

	log.Fingerprint = generateFingerprint(log)

	if err := s.repo.Create(ctx, log); err != nil {
		return nil, err
	}

	return toResponse(log), nil
}

func (s *service) Update(ctx context.Context, id string, req *Update) (*Entity, error) {
	log, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Level != "" {
		if !isValidLogLevel(req.Level) {
			return nil, ErrInvalidLogLevel
		}
		log.Level = Level(req.Level)
	}
	if req.Message != "" {
		log.Message = req.Message
	}
	if utils.DerefString(req.Context) != "" {
		log.Context = req.Context
	}

	log.Fingerprint = generateFingerprint(log)

	if err := s.repo.Update(ctx, id, log); err != nil {
		return nil, err
	}

	return toResponse(log), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func isValidLogLevel(level string) bool {
	switch Level(level) {
	case LevelInfo, LevelWarn, LevelError, LevelDebug:
		return true
	default:
		return false
	}
}

func generateFingerprint(e *Log) string {
	data := fmt.Sprintf(
		"%s:%s:%s",
		e.ProjectID,
		e.Level,
		e.Message,
	)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func toResponse(l *Log) *Entity {
	return &Entity{
		ID:      l.ID,
		Level:   string(l.Level),
		Message: l.Message,
		Context: l.Context,
		Time:    l.Time,
	}
}
