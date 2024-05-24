package repository

import (
	"context"
	"log/slog"
	"online-store/internal/entity"
	"online-store/pkg/utils"
)

func (r *Repository) CreateCategory(ctx context.Context, arg entity.Category) (string, error) {
	var id string

	query := `INSERT INTO categories
	(name, description)
	VALUES ($1, $2)
	RETURNING id`

	row := r.dbtx.QueryRowContext(ctx, query, arg.Name, arg.Description)

	err := row.Scan(&id)
	if err != nil {
		slog.Error(
			"Failed to CreateCategory Scan",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return id, utils.NewErrInternalServer("Failed to create category")
	}

	return id, nil
}

func (r *Repository) GetListCategory(ctx context.Context, arg entity.InGetListCategory) ([]entity.Category, int64, error) {
	var count int64
	items := []entity.Category{}

	queryCount := `SELECT COUNT(id) FROM categories`

	queryData := `SELECT id, name, description, created_at, updated_at FROM categories
	ORDER BY updated_at DESC
	LIMIT $1 OFFSET $2
	`

	row := r.dbtx.QueryRowContext(ctx, queryCount)
	err := row.Scan(&count)
	if err != nil {
		slog.Error(
			"Failed to GetListCategory QueryContext Count",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get category list")
	}

	if count == 0 {
		return items, count, nil
	}

	rows, err := r.dbtx.QueryContext(ctx, queryData, arg.Limit, arg.Offset)
	if err != nil {
		slog.Error(
			"Failed to GetListCategory QueryContext Data",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get category list")
	}
	defer rows.Close()

	for rows.Next() {
		var i entity.Category

		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			slog.Error(
				"Failed to GetListCategory rows.Next",
				slog.Any("err", err),
				slog.Any("arg", arg),
			)
			return items, count, utils.NewErrInternalServer("Failed to get category list")
		}

		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		slog.Error(
			"Failed to GetListCategory rows.Close",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get category list")
	}

	if err := rows.Err(); err != nil {
		slog.Error(
			"Failed to GetListCategory rows.Err",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get category list")
	}

	return items, count, nil
}
