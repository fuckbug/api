package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
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

func (r *repository) Create(ctx context.Context, user *User) error {
	const query = `
		INSERT INTO users (id, email, password, role, created_at, updated_at)
		VALUES (:id, :email, :password, :role, :created_at, :updated_at)
	`

	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	const query = `
		SELECT id, email, password, role
		FROM users
		WHERE email = $1
	`

	var user User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}
