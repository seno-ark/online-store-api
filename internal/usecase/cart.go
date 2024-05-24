package usecase

import (
	"context"
	"online-store/internal/entity"
	"online-store/internal/repository"
)

func (u *Usecase) CreateCartItem(ctx context.Context, userID string, arg entity.InCreateCartItem) (*entity.CartItem, error) {
	_, err := u.repo.GetProduct(ctx, arg.ProductID)
	if err != nil {
		return nil, err
	}

	var cartID string
	err = u.repo.WithTx(ctx, func(rtx *repository.Repository) error {
		var errTx error

		cartID, errTx = rtx.CreateCartItem(ctx, entity.CartItem{
			UserID:    userID,
			ProductID: arg.ProductID,
			Notes:     arg.Notes,
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return u.repo.GetCartItem(ctx, userID, cartID)
}

func (u *Usecase) DeleteCartItem(ctx context.Context, userID, cartID string) error {
	_, err := u.repo.GetCartItem(ctx, userID, cartID)
	if err != nil {
		return err
	}
	return u.repo.WithTx(ctx, func(rtx *repository.Repository) error {
		return rtx.DeleteCartItem(ctx, cartID)
	})
}

func (u *Usecase) GetListCartItem(ctx context.Context, userID string, arg entity.InGetListCartItem) ([]entity.OutGetCartItem, int64, error) {
	return u.repo.GetListCartItem(ctx, userID, arg)
}
