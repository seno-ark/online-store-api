package repository

import (
	"context"
	"log/slog"
	"online-store/internal/entity"
	"online-store/pkg/utils"
)

func (r *Repository) CreateOrderItem(ctx context.Context, arg entity.OrderItem) error {
	var id string

	query := `INSERT INTO order_items
	(order_id, product_id, notes, qty, product_price)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

	row := r.dbtx.QueryRowContext(ctx, query,
		arg.OrderID,
		arg.ProductID,
		arg.Notes,
		arg.Qty,
		arg.ProductPrice,
	)

	err := row.Scan(&id)
	if err != nil {
		slog.Error(
			"Failed to CreateOrderItem Scan",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return utils.NewErrInternalServer("Failed to create order item")
	}

	return nil
}

func (r *Repository) GetListOrderItem(ctx context.Context, orderID string, arg entity.InGetListOrderItem) ([]entity.OutGetOrderItem, int64, error) {
	var count int64
	items := []entity.OutGetOrderItem{}

	queryCount := `SELECT COUNT(id) FROM order_items WHERE order_id = $1`

	queryData := `SELECT 
	oi.id, oi.order_id, oi.product_id, oi.notes, oi.qty, oi.product_price, oi.created_at,
	p.name AS product_name, p.description AS product_description
	FROM order_items oi
	INNER JOIN products p ON p.ID = oi.product_id
	WHERE oi.order_id = $1
	ORDER BY oi.created_at
	LIMIT $2 OFFSET $3
	`

	row := r.dbtx.QueryRowContext(ctx, queryCount, orderID)
	err := row.Scan(&count)
	if err != nil {
		slog.Error(
			"Failed to GetListOrderItem QueryContext Count",
			slog.Any("err", err),
			slog.Any("orderID", orderID),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get order item list")
	}

	if count == 0 {
		return items, count, nil
	}

	rows, err := r.dbtx.QueryContext(ctx, queryData, orderID, arg.Limit, arg.Offset)
	if err != nil {
		slog.Error(
			"Failed to GetListOrderItem QueryContext Data",
			slog.Any("err", err),
			slog.Any("orderID", orderID),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get order item list")
	}
	defer rows.Close()

	for rows.Next() {
		var i entity.OutGetOrderItem

		if err := rows.Scan(
			&i.ID,
			&i.OrderID,
			&i.ProductID,
			&i.Notes,
			&i.Qty,
			&i.ProductPrice,
			&i.CreatedAt,
			&i.ProductName,
			&i.ProductDescription,
		); err != nil {
			slog.Error(
				"Failed to GetListOrderItem rows.Next",
				slog.Any("err", err),
				slog.Any("orderID", orderID),
				slog.Any("arg", arg),
			)
			return items, count, utils.NewErrInternalServer("Failed to get order item list")
		}

		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		slog.Error(
			"Failed to GetListOrderItem rows.Close",
			slog.Any("err", err),
			slog.Any("orderID", orderID),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get order item list")
	}

	if err := rows.Err(); err != nil {
		slog.Error(
			"Failed to GetListOrderItem rows.Err",
			slog.Any("err", err),
			slog.Any("orderID", orderID),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get order item list")
	}

	return items, count, nil
}
