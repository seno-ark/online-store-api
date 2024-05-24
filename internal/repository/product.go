package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"online-store/internal/entity"
	"online-store/pkg/utils"
)

func (r *Repository) CreateProduct(ctx context.Context, arg entity.Product) (string, error) {
	var id string

	query := `INSERT INTO products
	(category_id, name, description, price, stock)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

	row := r.dbtx.QueryRowContext(ctx, query,
		arg.CategoryID,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.Stock,
	)

	err := row.Scan(&id)
	if err != nil {
		slog.Error(
			"Failed to CreateProduct Scan",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return id, utils.NewErrInternalServer("Failed to create product")
	}

	return id, nil
}

func (r *Repository) UpdateProductStock(ctx context.Context, id string, stockChange int64) error {
	query := `UPDATE products 
	SET stock = stock + $1
	WHERE id = $2`

	_, err := r.dbtx.ExecContext(ctx, query, stockChange, id)
	if err != nil {
		slog.Error(
			"Failed to UpdateProductStock ExecContext",
			slog.Any("err", err),
			slog.Any("id", id),
			slog.Any("stockChange", stockChange),
		)
		return utils.NewErrInternalServer("Failed to update product stock")
	}

	return nil
}

func (r *Repository) GetProduct(ctx context.Context, id string) (*entity.OutGetProduct, error) {
	var i entity.OutGetProduct

	query := `SELECT 
	products.id, products.category_id, products.name, products.description, 
	products.price, products.stock, products.created_at, products.updated_at,
	categories.name AS category_name
	FROM products
	INNER JOIN categories ON categories.id = products.category_id
	WHERE products.id = $1
	LIMIT 1`

	row := r.dbtx.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&i.ID,
		&i.CategoryID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.Stock,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CategoryName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewErrNotFound("Product not found")
		}

		slog.Error(
			"Failed to GetProduct Scan",
			slog.Any("err", err),
			slog.Any("id", id),
		)
		return nil, utils.NewErrInternalServer("Failed to get product")
	}

	return &i, nil
}

func (r *Repository) GetListProduct(ctx context.Context, arg entity.InGetListProduct) ([]entity.OutGetProduct, int64, error) {
	var count int64
	items := []entity.OutGetProduct{}

	queryCount := `SELECT COUNT(id) FROM products`

	queryData := `SELECT products.id, products.category_id, products.name, 
	products.description, products.price, products.stock, products.created_at, products.updated_at,
	categories.name AS category_name 
	FROM products 
	INNER JOIN categories ON categories.id = products.category_id
	`

	queryArgs := []any{}
	if arg.CategoryID != "" {
		queryArgs = append(queryArgs, arg.CategoryID)
		whereQuery := " WHERE products.category_id = $1"

		queryCount += whereQuery
		queryData += whereQuery
	}

	row := r.dbtx.QueryRowContext(ctx, queryCount, queryArgs...)
	err := row.Scan(&count)
	if err != nil {
		slog.Error(
			"Failed to GetListProduct QueryContext Count",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get product list")
	}

	if count == 0 {
		return items, count, nil
	}

	queryArgs = append(queryArgs, arg.Limit)
	queryData += fmt.Sprintf(" ORDER BY products.updated_at DESC LIMIT $%d", len(queryArgs))

	queryArgs = append(queryArgs, arg.Offset)
	queryData += fmt.Sprintf(" OFFSET $%d", len(queryArgs))

	rows, err := r.dbtx.QueryContext(ctx, queryData, queryArgs...)
	if err != nil {
		slog.Error(
			"Failed to GetListProduct QueryContext Data",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get product list")
	}
	defer rows.Close()

	for rows.Next() {
		var i entity.OutGetProduct

		if err := rows.Scan(
			&i.ID,
			&i.CategoryID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Stock,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CategoryName,
		); err != nil {
			slog.Error(
				"Failed to GetListProduct rows.Next",
				slog.Any("err", err),
				slog.Any("arg", arg),
			)
			return items, count, utils.NewErrInternalServer("Failed to get product list")
		}

		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		slog.Error(
			"Failed to GetListProduct rows.Close",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get product list")
	}

	if err := rows.Err(); err != nil {
		slog.Error(
			"Failed to GetListProduct rows.Err",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get product list")
	}

	return items, count, nil
}
