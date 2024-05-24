package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"online-store/internal/entity"
	"online-store/pkg/constant"
	"online-store/pkg/utils"
)

func (r *Repository) CreateOrder(ctx context.Context, arg entity.Order) (string, error) {
	var id string

	query := `INSERT INTO orders
	(user_id, status, other_cost, total_cost, shipment_address)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

	row := r.dbtx.QueryRowContext(ctx, query,
		arg.UserID,
		constant.OrderTypeCreated,
		arg.OtherCost,
		arg.TotalCost,
		arg.ShipmentAddress,
	)
	err := row.Scan(&id)
	if err != nil {
		slog.Error(
			"Failed to CreateOrder Scan",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return id, utils.NewErrInternalServer("Failed to create order")
	}

	return id, nil
}

func (r *Repository) GetOrder(ctx context.Context, userID, id string) (*entity.Order, error) {
	query := `SELECT id, user_id, status, other_cost, total_cost, shipment_address, created_at, updated_at 
	FROM orders
	WHERE id = $1
	LIMIT 1`

	row := r.dbtx.QueryRowContext(ctx, query, id)
	var i entity.Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Status,
		&i.OtherCost,
		&i.TotalCost,
		&i.ShipmentAddress,
		&i.CreatedAt,
		&i.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewErrNotFound("Order not found")
		}

		slog.Error(
			"Failed to GetOrder",
			slog.Any("err", err),
			slog.Any("id", id),
			slog.Any("userID", userID),
		)
		return nil, utils.NewErrInternalServer("Failed to get order")
	}

	return &i, nil
}

func (r *Repository) GetListOrder(ctx context.Context, userID string, arg entity.InGetListOrder) ([]entity.Order, int64, error) {
	var count int64
	items := []entity.Order{}

	queryCount := `SELECT COUNT(id) FROM orders WHERE user_id = $1`

	queryData := `SELECT id, user_id, status, other_cost, total_cost, shipment_address, created_at, updated_at 
	FROM orders
	WHERE user_id = $1
	ORDER BY updated_at DESC
	LIMIT $2 OFFSET $3
	`

	row := r.dbtx.QueryRowContext(ctx, queryCount, userID)
	err := row.Scan(&count)
	if err != nil {
		slog.Error(
			"Failed to GetListOrder QueryContext Count",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get order list")
	}

	if count == 0 {
		return items, count, nil
	}

	rows, err := r.dbtx.QueryContext(ctx, queryData, userID, arg.Limit, arg.Offset)
	if err != nil {
		slog.Error(
			"Failed to GetListOrder QueryContext Data",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get order list")
	}
	defer rows.Close()

	for rows.Next() {
		var i entity.Order

		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Status,
			&i.OtherCost,
			&i.TotalCost,
			&i.ShipmentAddress,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			slog.Error(
				"Failed to GetListOrder rows.Next",
				slog.Any("err", err),
				slog.Any("arg", arg),
			)
			return items, count, utils.NewErrInternalServer("Failed to get order list")
		}

		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		slog.Error(
			"Failed to GetListOrder rows.Close",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get order list")
	}

	if err := rows.Err(); err != nil {
		slog.Error(
			"Failed to GetListOrder rows.Err",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return items, count, utils.NewErrInternalServer("Failed to get order list")
	}

	return items, count, nil
}
