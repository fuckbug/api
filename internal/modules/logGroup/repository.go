package loggroup

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("not found")

type Repository interface {
	GetAll(ctx context.Context, params GetAllParams) ([]*Group, error)
	Count(ctx context.Context, params FilterParams) (int, error)
	GetByID(ctx context.Context, id string) (*Group, error)
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

func (r *repository) GetAll(ctx context.Context, params GetAllParams) ([]*Group, error) {
	query := `
        SELECT id, project_id, level, message, first_seen_at, last_seen_at, counter, status 
        FROM log_groups 
        WHERE 1=1
    `

	args := map[string]interface{}{
		"limit":  params.Limit,
		"offset": params.Offset,
	}

	query, args = applyFilters(query, params.FilterParams, args)

	query += " ORDER BY last_seen_at " + params.SortOrder
	query += " LIMIT :limit OFFSET :offset"

	query, namedArgs, err := sqlx.Named(query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare named query: %w", err)
	}

	query = r.db.Rebind(query)

	r.logger.Debug(query)

	var entities []*Group
	err = r.db.SelectContext(ctx, &entities, query, namedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get log groups: %w", err)
	}
	return entities, nil
}

func (r *repository) Count(ctx context.Context, params FilterParams) (int, error) {
	query := "SELECT COUNT(*) FROM log_groups WHERE 1=1"
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

func (r *repository) GetByID(ctx context.Context, id string) (*Group, error) {
	const query = `SELECT id, project_id, level, message, first_seen_at, last_seen_at, counter, status 
		FROM log_groups WHERE id = $1`

	var entity Group
	err := r.db.GetContext(ctx, &entity, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get log by id: %w", err)
	}
	return &entity, nil
}

func applyFilters(baseQuery string, params FilterParams, args map[string]interface{}) (string, map[string]interface{}) {
	query := baseQuery

	if params.ProjectID != "" {
		query += " AND project_id = :projectId"
		args["projectId"] = params.ProjectID
	}

	if params.TimeFrom != 0 {
		query += " AND last_seen_at >= :timeFrom"
		args["timeFrom"] = params.TimeFrom
	}

	if params.TimeTo != 0 {
		query += " AND last_seen_at <= :timeTo"
		args["timeTo"] = params.TimeTo
	}

	if params.Level != "" {
		query += " AND level = :level"
		args["level"] = params.Level
	}

	if params.Search != "" {
		query += " AND message ILIKE :search"
		args["search"] = "%" + params.Search + "%"
	}

	return query, args
}
