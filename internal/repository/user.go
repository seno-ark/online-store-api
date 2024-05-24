package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"online-store/internal/entity"
	"online-store/pkg/utils"
)

func (r *Repository) CreateUser(ctx context.Context, arg entity.User) (string, error) {
	var id string

	query := `INSERT INTO users
	(email, password, full_name)
	VALUES ($1, $2, $3)
	RETURNING id`

	row := r.dbtx.QueryRowContext(ctx, query, arg.Email, arg.Password, arg.FullName)
	err := row.Scan(&id)
	if err != nil {
		slog.Error(
			"Failed to CreateUser",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return id, utils.NewErrInternalServer("Failed to create user")
	}

	return id, err
}

func (r *Repository) GetUser(ctx context.Context, id string) (*entity.User, error) {
	var i entity.User

	query := `SELECT id, email, password, full_name, created_at, updated_at 
	FROM users
	WHERE id = $1`

	row := r.dbtx.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.FullName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewErrNotFound("User not found")
		}
		slog.Error(
			"Failed to CreateUser",
			slog.Any("err", err),
			slog.Any("id", id),
		)
		return nil, utils.NewErrInternalServer("Failed to get user")
	}

	return &i, err
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var i entity.User

	query := `SELECT id, email, password, full_name, created_at, updated_at 
	FROM users
	WHERE email = $1`

	row := r.dbtx.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.FullName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewErrNotFound("User not found")
		}
		slog.Error(
			"Failed to CreateUser",
			slog.Any("err", err),
			slog.Any("email", email),
		)
		return nil, utils.NewErrInternalServer("Failed to get user")
	}

	return &i, err
}
