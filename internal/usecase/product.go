package usecase

import (
	"context"
	"online-store/internal/entity"
)

func (u *Usecase) GetListProduct(ctx context.Context, arg entity.InGetListProduct) ([]entity.OutGetProduct, int64, error) {
	return u.repo.GetListProduct(ctx, arg)
}
