package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"online-store/internal/entity"
	"online-store/pkg/utils"
)

func (r *Repository) CreateOrderPayment(ctx context.Context, arg entity.Payment) (string, error) {
	var id string

	query := `INSERT INTO payments
	(order_id, payment_method, payment_provider, bill_amount, paid_amount, status, transaction_id, paid_at, log)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id`

	row := r.dbtx.QueryRowContext(ctx, query,
		arg.OrderID,
		arg.PaymentMethod,
		arg.PaymentProvider,
		arg.BillAmount,
		arg.PaidAmount,
		arg.Status,
		arg.TransactionID,
		arg.PaidAt,
		arg.Log,
	)

	err := row.Scan(&id)
	if err != nil {
		slog.Error(
			"Failed to CreateOrderPayment Scan",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return id, utils.NewErrInternalServer("Failed to create order payment")
	}

	return id, nil
}

func (r *Repository) UpdateOrderPayment(ctx context.Context, arg entity.Payment) error {
	query := `UPDATE payments
	SET 
	paid_amount = $1,
	status = $2, 
	transaction_id = $3,
	paid_at = $4,
	log = $5
	WHERE order_id = $6`

	_, err := r.dbtx.ExecContext(ctx, query,
		arg.PaidAmount,
		arg.Status,
		arg.TransactionID,
		arg.PaidAt,
		arg.Log,
		arg.OrderID,
	)
	if err != nil {
		slog.Error(
			"Failed to UpdateOrderPayment ExecContext",
			slog.Any("err", err),
			slog.Any("arg", arg),
		)
		return utils.NewErrInternalServer("Failed to update order payment")
	}

	return nil
}

func (r *Repository) GetOrderPayment(ctx context.Context, id string) (*entity.Payment, error) {
	var i entity.Payment

	query := `SELECT id, order_id, payment_method, payment_provider, bill_amount, paid_amount, status, transaction_id, paid_at, created_at, updated_at 
	FROM payments
	WHERE order_id = $1
	LIMIT 1`

	row := r.dbtx.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.PaymentMethod,
		&i.PaymentProvider,
		&i.BillAmount,
		&i.PaidAmount,
		&i.Status,
		&i.TransactionID,
		&i.PaidAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewErrNotFound("Payment not found")
		}

		slog.Error(
			"Failed to GetOrderPayment Scan",
			slog.Any("err", err),
			slog.Any("id", id),
		)
		return nil, utils.NewErrInternalServer("Failed to get payment")
	}

	return &i, nil
}
