package log

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	loggroup "github.com/fuckbug/api/internal/modules/logGroup"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("not found")

type Repository interface {
	GetAll(ctx context.Context, params GetAllParams) ([]*Log, error)
	Count(ctx context.Context, params FilterParams) (int, error)
	GetByID(ctx context.Context, id string) (*Log, error)
	Create(ctx context.Context, log *Log) error
	Update(ctx context.Context, id string, log *Log) error
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

func (r *repository) GetAll(ctx context.Context, params GetAllParams) ([]*Log, error) {
	query := `
        SELECT id, project_id, level, message, context, time, created_at, updated_at 
        FROM logs 
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

	var logs []*Log
	err = r.db.SelectContext(ctx, &logs, query, namedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	return logs, nil
}

func (r *repository) Count(ctx context.Context, params FilterParams) (int, error) {
	query := "SELECT COUNT(*) FROM logs WHERE 1=1"
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

func (r *repository) GetByID(ctx context.Context, id string) (*Log, error) {
	const query = `SELECT id, level, message, context, time, created_at, updated_at 
		FROM logs WHERE id = $1`

	var entity Log
	err := r.db.GetContext(ctx, &entity, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get log by id: %w", err)
	}
	return &entity, nil
}

func (r *repository) Create(ctx context.Context, l *Log) error {

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			r.logger.Warn(fmt.Sprintf("failed to rollback transaction: %v", err))
		}
	}()

	const logGroupQuery = `
        INSERT INTO error_groups (id, project_id, file, line, message, first_seen_at, last_seen_at, counter)
        VALUES (:id, :project_id, :file, :line, :message, :first_seen_at, :last_seen_at, 1)
        ON CONFLICT (id) DO UPDATE 
        SET counter = error_groups.counter + 1, last_seen_at = EXCLUDED.last_seen_at
    `

	now := time.Now().Unix()
	logGroup := loggroup.Group{
		ID:          l.Fingerprint,
		ProjectID:   l.ProjectID,
		Level:       loggroup.Level(l.Level),
		Message:     l.Message,
		FirstSeenAt: int(now),
		LastSeenAt:  int(now),
		Counter:     0,
		Status:      loggroup.StatusUnresolved,
	}

	_, err = tx.NamedExecContext(ctx, logGroupQuery, logGroup)
	if err != nil {
		return fmt.Errorf("failed to upsert log group: %w", err)
	}

	const query = `
		INSERT INTO logs (
	  		id, project_id, fingerprint, level, message, context, time, created_at, updated_at
		) VALUES (
	  		:id, :project_id, :fingerprint, :level, :message, :context, :time, :created_at, :updated_at
		)
	`

	if l.ID == "" {
		l.ID = uuid.New().String()
	}

	l.CreatedAt = now
	l.UpdatedAt = now

	_, err = r.db.NamedExecContext(ctx, query, l)
	if err != nil {
		return fmt.Errorf("failed to create log: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string, updated *Log) error {
	const query = `
		UPDATE
		    logs 
		SET
		    fingerprint = :fingerprint,
		    level = :level,
		    message = :message,
		    context = :context,
		    time = :time,
		    updated_at = :updated_at
		WHERE
		    id = :id
	`

	updated.ID = id
	updated.UpdatedAt = time.Now().UnixMilli()

	result, err := r.db.NamedExecContext(ctx, query, updated)
	if err != nil {
		return fmt.Errorf("failed to update log: %w", err)
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
	const query = `DELETE FROM logs WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete log: %w", err)
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

	if params.Fingerprint != "" {
		query += " AND fingerprint = :fingerprint"
		args["fingerprint"] = params.Fingerprint
	}

	if params.TimeFrom != 0 {
		query += " AND time >= :timeFrom"
		args["timeFrom"] = params.TimeFrom
	}

	if params.TimeTo != 0 {
		query += " AND time <= :timeTo"
		args["timeTo"] = params.TimeTo
	}

	if params.LevelFilter != "" {
		query += " AND level = :level"
		args["level"] = params.LevelFilter
	}

	if params.SearchQuery != "" {
		query += " AND message LIKE :search"
		args["search"] = "%" + params.SearchQuery + "%"
	}

	return query, args
}
