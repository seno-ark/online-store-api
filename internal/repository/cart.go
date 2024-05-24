package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"online-store/internal/entity"
	"online-store/pkg/utils"
)

func (r *Repository) CreateCartItem(ctx context.Context, arg entity.CartItem) (string, error) {
	var id string

	query := `INSERT INTO carts
	(user_id, product_id, notes)
	VALUES ($1, $2, $3)
	ON CONFLICT(user_id, product_id) 
	DO UPDATE SET notes = $4, updated_at = NOW()
	RETURNING id`

	row := r.dbtx.QueryRowContext(ctx, query,
		arg.UserID,
		arg.ProductID,
		arg.Notes,
		arg.Notes,
	)
	err := row.Scan(&id)
	if err != nil {
		slog.Error(
			"Failed to CreateCartItem Scan",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return id, utils.NewErrInternalServer("Failed to create cart")
	}

	return id, nil
}

func (r *Repository) DeleteCartItem(ctx context.Context, id string) error {
	query := `DELETE FROM carts
	WHERE id = $1`

	_, err := r.dbtx.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error(
			"Failed to DeleteCartItem",
			slog.Any("err", err),
			slog.Any("id", id),
		)
		return utils.NewErrInternalServer("Failed to delete cart")
	}

	return nil
}

func (r *Repository) GetCartItem(ctx context.Context, userID, id string) (*entity.CartItem, error) {
	var i entity.CartItem

	query := `SELECT 
	id, user_id, product_id, notes, created_at, updated_at
	FROM carts 
	WHERE id = $1 AND user_id = $2`

	row := r.dbtx.QueryRowContext(ctx, query, id, userID)
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProductID,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewErrNotFound("Cart not found")
		}

		slog.Error(
			"Failed to GetCartItem",
			slog.Any("err", err),
			slog.Any("id", id),
			slog.Any("userID", userID),
		)
		return nil, utils.NewErrInternalServer("Failed to get cart")
	}

	return &i, nil
}

func (r *Repository) GetCartItemByProduct(ctx context.Context, userID, productID string) (*entity.CartItem, error) {
	var i entity.CartItem

	query := `SELECT 
	id, user_id, product_id, notes, created_at, updated_at
	FROM carts 
	WHERE user_id = $1 AND product_id = $2`

	row := r.dbtx.QueryRowContext(ctx, query, userID, productID)
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProductID,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewErrNotFound("Cart not found")
		}

		slog.Error(
			"Failed to GetCartItemByProduct",
			slog.Any("err", err),
			slog.Any("productID", productID),
			slog.Any("userID", userID),
		)
		return nil, utils.NewErrInternalServer("Failed to get cart")
	}

	return &i, nil
}

func (r *Repository) GetListCartItem(ctx context.Context, userID string, arg entity.InGetListCartItem) ([]entity.OutGetCartItem, int64, error) {
	var count int64
	items := []entity.OutGetCartItem{}

	queryCount := `SELECT COUNT(id) FROM carts
	WHERE carts.user_id = $1`

	queryData := `SELECT 
	carts.id, carts.user_id, carts.product_id, carts.notes, carts.created_at, carts.updated_at,
	products.category_id, products.name AS product_name, products.description AS product_description,
	products.price AS product_price, products.stock AS product_stock,
	categories.name AS category_name
	FROM carts
	INNER JOIN products ON products.ID = carts.product_id
	INNER JOIN categories ON categories.ID = products.category_id
	WHERE carts.user_id = $1
	ORDER BY carts.updated_at DESC
	LIMIT $2 OFFSET $3
	`

	row := r.dbtx.QueryRowContext(ctx, queryCount, userID)
	err := row.Scan(&count)
	if err != nil {
		slog.Error(
			"Failed to GetListCartItem QueryContext Count",
			slog.Any("err", err),
			slog.Any("userID", userID),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get cart item list")
	}

	if count == 0 {
		return items, count, nil
	}

	rows, err := r.dbtx.QueryContext(ctx, queryData, userID, arg.Limit, arg.Offset)
	if err != nil {
		slog.Error(
			"Failed to GetListCartItem QueryContext Data",
			slog.Any("err", err),
			slog.Any("userID", userID),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get cart item list")
	}
	defer rows.Close()

	for rows.Next() {
		var i entity.OutGetCartItem

		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ProductID,
			&i.Notes,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CategoryID,
			&i.ProductName,
			&i.ProductDescription,
			&i.ProductPrice,
			&i.ProductStock,
			&i.CategoryName,
		); err != nil {
			slog.Error(
				"Failed to GetListCartItem rows.Next",
				slog.Any("err", err),
				slog.Any("userID", userID),
				slog.Any("arg", arg),
			)
			return items, count, utils.NewErrInternalServer("Failed to get cart item list")
		}

		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		slog.Error(
			"Failed to GetListCartItem rows.Close",
			slog.Any("err", err),
			slog.Any("userID", userID),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get cart item list")
	}

	if err := rows.Err(); err != nil {
		slog.Error(
			"Failed to GetListCartItem rows.Err",
			slog.Any("err", err),
			slog.Any("userID", userID),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get cart item list")
	}

	return items, count, nil
}
