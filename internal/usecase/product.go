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

func (u *Usecase) GetListProduct(ctx context.Context, arg entity.InGetListProduct) ([]entity.OutGetProduct, int64, error) {
	var (
		result []entity.OutGetProduct
		total  int64
		err    error
	)

	cacheKeyCount := fmt.Sprintf(constant.CacheKeyListProductByCategoryCount, arg.CategoryID)
	cacheKeyList := fmt.Sprintf(constant.CacheKeyListProductByCategory, arg.CategoryID, arg.Limit, arg.Offset)

	cacheData, _ := u.cache.Get(ctx, cacheKeyCount)
	if cacheData != "" {
		total, err = strconv.ParseInt(cacheData, 10, 64)
		if err != nil {
			slog.Error(
				"Failed to GetListProduct Parse Cache",
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
					"Failed to GetListProduct Parse Cache",
					slog.Any("err", err),
					slog.Any("cacheKey", cacheKeyList),
					slog.Any("cacheData", cacheData),
				)
				return nil, 0, utils.NewErrInternalServer("Failed to parse cache data")
			}

			return result, total, nil
		}
	}

	result, total, err = u.repo.GetListProduct(ctx, arg)
	if err != nil {
		return nil, 0, err
	}

	jsonList, err := json.Marshal(result)
	if err != nil {
		slog.Error(
			"Failed to GetListProduct Marshal Cache",
			slog.Any("err", err),
			slog.Any("data", result),
		)
		return nil, 0, err
	}

	_ = u.cache.Set(ctx, cacheKeyCount, []byte(fmt.Sprintf("%d", total)), constant.DefaultExpiration)
	_ = u.cache.Set(ctx, cacheKeyList, jsonList, constant.DefaultExpiration)

	return result, total, nil
}

func (u *Usecase) DummyProductAndCategory(ctx context.Context) error {
	_, total, err := u.repo.GetListCategory(ctx, entity.InGetListCategory{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		return err
	}

	if total > 0 {
		return nil
	}

	err = u.repo.WithTx(ctx, func(rtx *repository.Repository) error {

		id, err := rtx.CreateCategory(ctx, entity.Category{
			Name:        "Pakaian",
			Description: "Baju, Celana, dan aksesoris untuk semua umur",
		})
		if err != nil {
			return err
		}

		_, err = rtx.CreateProduct(ctx, entity.Product{
			CategoryID:  id,
			Name:        "Celana panjang anak",
			Description: "Untuk umur 6 tahun",
			Price:       100000,
			Stock:       10,
		})
		if err != nil {
			return err
		}

		_, err = rtx.CreateProduct(ctx, entity.Product{
			CategoryID:  id,
			Name:        "Celana pendek anak laki-laki",
			Description: "Untuk umur 8 tahun",
			Price:       80000,
			Stock:       10,
		})
		if err != nil {
			return err
		}

		id, err = rtx.CreateCategory(ctx, entity.Category{
			Name:        "Elektronik",
			Description: "Peralatan elektronik berkualitas",
		})
		if err != nil {
			return err
		}

		_, err = rtx.CreateProduct(ctx, entity.Product{
			CategoryID:  id,
			Name:        "TV Samsung 32 Inch",
			Description: "Bergaransi seumur jagung",
			Price:       3000000,
			Stock:       5,
		})
		if err != nil {
			return err
		}

		_, err = rtx.CreateProduct(ctx, entity.Product{
			CategoryID:  id,
			Name:        "Kipas angin kesedot sampah",
			Description: "Wuss",
			Price:       180000,
			Stock:       8,
		})
		if err != nil {
			return err
		}

		_, err = rtx.CreateCategory(ctx, entity.Category{
			Name:        "Makanan",
			Description: "Snack dan minuman",
		})
		return err
	})
	if err != nil {
		panic(err)
	}

	return nil
}
