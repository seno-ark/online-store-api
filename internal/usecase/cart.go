package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"online-store/internal/entity"
	"online-store/internal/repository"
	"online-store/pkg/constant"
	"online-store/pkg/utils"
	"strconv"
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

	cacheKey := fmt.Sprintf(constant.CacheKeyListUserCartAll, userID)
	_ = u.cache.Del(ctx, cacheKey)

	return u.repo.GetCartItem(ctx, userID, cartID)
}

func (u *Usecase) DeleteCartItem(ctx context.Context, userID, cartID string) error {
	_, err := u.repo.GetCartItem(ctx, userID, cartID)
	if err != nil {
		return err
	}

	err = u.repo.WithTx(ctx, func(rtx *repository.Repository) error {
		return rtx.DeleteCartItem(ctx, cartID)
	})
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf(constant.CacheKeyListUserCartAll, userID)
	_ = u.cache.Del(ctx, cacheKey)

	return nil
}

func (u *Usecase) GetListCartItem(ctx context.Context, userID string, arg entity.InGetListCartItem) ([]entity.OutGetCartItem, int64, error) {
	var (
		result []entity.OutGetCartItem
		total  int64
		err    error
	)

	cacheKeyCount := fmt.Sprintf(constant.CacheKeyListUserCartCount, userID)
	cacheKeyList := fmt.Sprintf(constant.CacheKeyListUserCart, userID, arg.Limit, arg.Offset)

	cacheData, _ := u.cache.Get(ctx, cacheKeyCount)
	if cacheData != "" {
		total, err = strconv.ParseInt(cacheData, 10, 64)
		if err != nil {
			slog.Error(
				"Failed to GetListCartItem Parse Cache",
				slog.Any("err", err),
				slog.Any("cacheKey", cacheKeyCount),
				slog.Any("cacheData", cacheData),
			)
			return nil, 0, utils.NewErrInternalServer("Failed to parse cache data")
		}

		cacheData, _ = u.cache.Get(ctx, cacheKeyList)
		if cacheData != "" {
			err = json.Unmarshal([]byte(cacheData), &result)
			if err != nil {
				slog.Error(
					"Failed to GetListCartItem Parse Cache",
					slog.Any("err", err),
					slog.Any("cacheKey", cacheKeyList),
					slog.Any("cacheData", cacheData),
				)
				return nil, 0, utils.NewErrInternalServer("Failed to parse cache data")
			}

			return result, total, nil
		}
	}

	result, total, err = u.repo.GetListCartItem(ctx, userID, arg)
	if err != nil {
		return nil, 0, err
	}

	jsonList, err := json.Marshal(result)
	if err != nil {
		slog.Error(
			"Failed to GetListCartItem Marshal Cache",
			slog.Any("err", err),
			slog.Any("data", result),
		)
		return nil, 0, err
	}

	_ = u.cache.Set(ctx, cacheKeyCount, []byte(fmt.Sprintf("%d", total)), constant.DefaultExpiration)
	_ = u.cache.Set(ctx, cacheKeyList, jsonList, constant.DefaultExpiration)

	return result, total, nil
}
