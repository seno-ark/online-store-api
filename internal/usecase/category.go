package usecase

import (
	"context"
	"online-store/internal/entity"
)

func (u *Usecase) GetListCategory(ctx context.Context, arg entity.InGetListCategory) ([]entity.Category, int64, error) {
	return u.repo.GetListCategory(ctx, arg)
}
