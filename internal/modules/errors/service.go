package errors

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

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

func (s *service) GetStats(ctx context.Context, projectID string, fingerprint string) (*Stats, error) {
	stats, err := s.repo.GetStats(ctx, projectID, fingerprint)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *service) Create(ctx context.Context, req *Create) (*Entity, error) {
	stacktrace, err := stacktraceToString(req.Stacktrace)
	if err != nil {
		return nil, err
	}

	contextStr, err := contextToStringPtr(req.Context)
	if err != nil {
		return nil, err
	}

	headers, err := mapToStringPtr(req.Headers)
	if err != nil {
		return nil, err
	}

	queryParams, err := mapToStringPtr(req.QueryParams)
	if err != nil {
		return nil, err
	}

	bodyParams, err := mapToStringPtr(req.BodyParams)
	if err != nil {
		return nil, err
	}

	cookies, err := mapToStringPtr(req.Cookies)
	if err != nil {
		return nil, err
	}

	session, err := mapToStringPtr(req.Session)
	if err != nil {
		return nil, err
	}

	files, err := mapToStringPtr(req.Files)
	if err != nil {
		return nil, err
	}

	env, err := mapToStringPtr(req.Env)
	if err != nil {
		return nil, err
	}

	entity := &Error{
		ID:          uuid.New().String(),
		ProjectID:   req.ProjectID,
		Message:     req.Message,
		Stacktrace:  stacktrace,
		File:        req.File,
		Line:        req.Line,
		Context:     contextStr,
		IP:          req.IP,
		URL:         req.URL,
		Method:      req.Method,
		Headers:     headers,
		QueryParams: queryParams,
		BodyParams:  bodyParams,
		Cookies:     cookies,
		Session:     session,
		Files:       files,
		Env:         env,
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
	stacktrace, err := stacktraceToString(req.Stacktrace)
	if err != nil {
		return nil, err
	}
	entity.Stacktrace = stacktrace
	if req.File != "" {
		entity.File = req.File
	}
	if req.Line != 0 {
		entity.Line = req.Line
	}
	contextStr, err := contextToStringPtr(req.Context)
	if err != nil {
		return nil, err
	}
	entity.Context = contextStr

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
		"%s:%s:%s:%d",
		cleanMsg,
		e.ProjectID,
		e.File,
		e.Line,
	)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func toResponse(e *Error) *Entity {
	response := &Entity{
		ID:      e.ID,
		Message: e.Message,
		File:    e.File,
		Line:    e.Line,
		IP:      e.IP,
		URL:     e.URL,
		Method:  e.Method,
		Time:    e.Time,
	}

	if err := parseJSONField(&e.Stacktrace, &response.Stacktrace); err != nil {
		*response.Stacktrace = e.Stacktrace
	}

	if err := parseJSONField(e.Context, &response.Context); err != nil {
		*response.Context = e.Context
	}

	if err := parseJSONField(e.Headers, &response.Headers); err != nil {
		*response.Headers = map[string]interface{}{
			"Headers": e.Headers,
		}
	}

	if err := parseJSONField(e.QueryParams, &response.QueryParams); err != nil {
		*response.QueryParams = map[string]interface{}{
			"QueryParams": e.QueryParams,
		}
	}

	if err := parseJSONField(e.BodyParams, &response.BodyParams); err != nil {
		*response.BodyParams = map[string]interface{}{
			"BodyParams": e.BodyParams,
		}
	}

	if err := parseJSONField(e.Cookies, &response.Cookies); err != nil {
		*response.Cookies = map[string]interface{}{
			"Cookies": e.Cookies,
		}
	}

	if err := parseJSONField(e.Session, &response.Session); err != nil {
		*response.Session = map[string]interface{}{
			"Session": e.Session,
		}
	}

	if err := parseJSONField(e.Files, &response.Files); err != nil {
		*response.Files = map[string]interface{}{
			"Files": e.Files,
		}
	}

	if err := parseJSONField(e.Env, &response.Env); err != nil {
		*response.Env = map[string]interface{}{
			"Env": e.Env,
		}
	}

	return response
}

func parseJSONField(src *string, dest interface{}) error {
	if src == nil || *src == "" {
		return nil
	}
	return json.Unmarshal([]byte(*src), dest)
}

func contextToStringPtr(context *interface{}) (*string, error) {
	if context != nil {
		jsonData, err := json.Marshal(*context)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal context: %w", err)
		}

		str := string(jsonData)
		return &str, nil
	}
	return nil, nil
}

func stacktraceToString(stacktrace interface{}) (string, error) {
	jsonData, err := json.Marshal(stacktrace)
	if err != nil {
		return "", fmt.Errorf("failed to marshal context: %w", err)
	}
	return string(jsonData), nil
}

func mapToStringPtr(m *map[string]interface{}) (*string, error) {
	if m == nil {
		return nil, nil
	}

	jsonBytes, err := json.Marshal(*m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal map to JSON: %w", err)
	}

	result := string(jsonBytes)
	return &result, nil
}
