package errors

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("not found")

type Repository interface {
	GetAll(ctx context.Context, params GetAllParams) ([]*Error, error)
	Count(ctx context.Context, params FilterParams) (int, error)
	GetByID(ctx context.Context, id string) (*Error, error)
	Create(ctx context.Context, entity *Error) error
	Update(ctx context.Context, id string, entity *Error) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db     *sqlx.DB
	logger Logger
}

func NewRepository(db *sqlx.DB, logger Logger) Repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

func (r *repository) GetAll(ctx context.Context, params GetAllParams) ([]*Error, error) {
	query := `
        SELECT id, project_id, fingerprint, message, stacktrace, file, line, context, time, created_at, updated_at 
        FROM errors 
        WHERE 1=1
    `

	args := map[string]interface{}{
		"limit":  params.Limit,
		"offset": params.Offset,
	}

	query, args = applyFilters(query, params.FilterParams, args)

	query += " ORDER BY time " + params.SortOrder
	query += " LIMIT :limit OFFSET :offset"

	query, namedArgs, err := sqlx.Named(query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare named query: %w", err)
	}

	query = r.db.Rebind(query)

	r.logger.Debug(query)

	var entities []*Error
	err = r.db.SelectContext(ctx, &entities, query, namedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get errors: %w", err)
	}
	return entities, nil
}

func (r *repository) Count(ctx context.Context, params FilterParams) (int, error) {
	query := "SELECT COUNT(*) FROM errors WHERE 1=1"
	query, args := applyFilters(query, params, make(map[string]interface{}))

	query, namedArgs, err := sqlx.Named(query, args)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare named query: %w", err)
	}

	query = r.db.Rebind(query)

	r.logger.Debug(query)

	var count int
	err = r.db.GetContext(ctx, &count, query, namedArgs...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*Error, error) {
	const query = `SELECT id, project_id, fingerprint, message, stacktrace, file, line, context, time, created_at, updated_at 
		FROM errors WHERE id = $1`

	var entity Error
	err := r.db.GetContext(ctx, &entity, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get error by id: %w", err)
	}
	return &entity, nil
}

func (r *repository) Create(ctx context.Context, e *Error) error {
	const query = `INSERT INTO errors 
		(id, project_id, fingerprint, message, stacktrace, file, line, context, time, created_at, updated_at)
		VALUES (:id, :project_id, :fingerprint, :message, :stacktrace, :file, :line, :context, :time, :created_at, :updated_at)`

	if e.ID == "" {
		e.ID = uuid.New().String()
	}

	now := time.Now().UnixMilli()
	e.CreatedAt = now
	e.UpdatedAt = now

	_, err := r.db.NamedExecContext(ctx, query, e)
	if err != nil {
		return fmt.Errorf("failed to create error: %w", err)
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string, updated *Error) error {
	const query = `UPDATE errors 
		SET 
		    fingerprint = :fingerprint,
		    message = :message,
		    stacktrace = :stacktrace,
		    file = :file,
		    line = :line,
		    context = :context,
		    time = :time,
		    updated_at = :updated_at
		WHERE id = :id`

	updated.ID = id
	updated.UpdatedAt = time.Now().UnixMilli()

	result, err := r.db.NamedExecContext(ctx, query, updated)
	if err != nil {
		return fmt.Errorf("failed to update error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM errors WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func applyFilters(baseQuery string, params FilterParams, args map[string]interface{}) (string, map[string]interface{}) {
	query := baseQuery

	if params.ProjectID != "" {
		query += " AND project_id = :projectId"
		args["projectId"] = params.ProjectID
	}

	if params.TimeFrom != 0 {
		query += " AND time >= :timeFrom"
		args["timeFrom"] = params.TimeFrom
	}

	if params.TimeTo != 0 {
		query += " AND time <= :timeTo"
		args["timeTo"] = params.TimeTo
	}

	if params.SearchQuery != "" {
		query += " AND message LIKE :search"
		args["search"] = "%" + params.SearchQuery + "%"
	}

	return query, args
}
