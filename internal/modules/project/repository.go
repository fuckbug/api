package project

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound = errors.New("not found")
	KeyLength   = 64
)

type Repository interface {
	GetAll(ctx context.Context, params GetAllParams) ([]*Project, error)
	Count(ctx context.Context) (int, error)
	GetByID(ctx context.Context, id string) (*Project, error)
	Create(ctx context.Context, project *Project) error
	Update(ctx context.Context, id string, project *Project) error
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

func (r *repository) GetAll(ctx context.Context, params GetAllParams) ([]*Project, error) {
	query := `
        SELECT id, name, public_key, created_at, updated_at, deleted_at
        FROM projects 
        WHERE 1=1
    `

	args := map[string]interface{}{
		"limit":  params.Limit,
		"offset": params.Offset,
	}

	query += " ORDER BY id " + params.SortOrder
	query += " LIMIT :limit OFFSET :offset"

	query, namedArgs, err := sqlx.Named(query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare named query: %w", err)
	}

	query = r.db.Rebind(query)

	r.logger.Debug(query)

	var projects []*Project
	err = r.db.SelectContext(ctx, &projects, query, namedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	return projects, nil
}

func (r *repository) Count(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM projects WHERE deleted_at IS NULL"

	var count int
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*Project, error) {
	const query = `SELECT id, name, public_key, created_at, updated_at, deleted_at 
		FROM projects WHERE id = $1`

	var entity Project
	err := r.db.GetContext(ctx, &entity, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get project by id: %w", err)
	}
	return &entity, nil
}

func (r *repository) Create(ctx context.Context, l *Project) error {
	const query = `INSERT INTO projects 
		(id, name, public_key, created_at, updated_at, deleted_at)
		VALUES (:id, :name, :public_key, :created_at, :updated_at, null)`

	if l.ID == "" {
		l.ID = uuid.New().String()
	}

	now := time.Now().UnixMilli()
	l.PublicKey = generateRandomKey()
	l.CreatedAt = now
	l.UpdatedAt = now

	_, err := r.db.NamedExecContext(ctx, query, l)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string, updated *Project) error {
	const query = `UPDATE projects 
		SET name = :name,
		    updated_at = :updated_at
		WHERE id = :id`

	updated.ID = id
	updated.UpdatedAt = time.Now().UnixMilli()

	result, err := r.db.NamedExecContext(ctx, query, updated)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
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
	const query = `DELETE FROM projects WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
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

func generateRandomKey() string {
	b := make([]byte, KeyLength)
	_, err := rand.Read(b)
	if err != nil {
		return uuid.New().String()
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
