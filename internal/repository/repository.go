package repository

import (
	"context"
	"database/sql"
	"fmt"
	"online-store/pkg/database"
)

type Repository struct {
	db   *sql.DB
	dbtx database.DBTX
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db:   db,
		dbtx: db,
	}
}

func (r *Repository) WithTx(ctx context.Context, fn func(*Repository) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := &Repository{dbtx: tx}
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
