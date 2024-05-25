package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"online-store/internal/entity"
	"online-store/pkg/constant"
	"online-store/pkg/utils"
	"strconv"
)

func (u *Usecase) GetListCategory(ctx context.Context, arg entity.InGetListCategory) ([]entity.Category, int64, error) {
	var (
		result []entity.Category
		total  int64
		err    error
	)

	cacheKeyCount := fmt.Sprintf(constant.CacheKeyListCategoryCount)
	cacheKeyList := fmt.Sprintf(constant.CacheKeyListCategory, arg.Limit, arg.Offset)

	cacheData, _ := u.cache.Get(ctx, cacheKeyCount)
	if cacheData != "" {
		total, err = strconv.ParseInt(cacheData, 10, 64)
		if err != nil {
			slog.Error(
				"Failed to GetListCategory Parse Cache",
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
					"Failed to GetListCategory Parse Cache",
					slog.Any("err", err),
					slog.Any("cacheKey", cacheKeyList),
					slog.Any("cacheData", cacheData),
				)
				return nil, 0, utils.NewErrInternalServer("Failed to parse cache data")
			}

			return result, total, nil
		}
	}

	result, total, err = u.repo.GetListCategory(ctx, arg)
	if err != nil {
		return nil, 0, err
	}

	jsonList, err := json.Marshal(result)
	if err != nil {
		slog.Error(
			"Failed to GetListCategory Marshal Cache",
			slog.Any("err", err),
			slog.Any("data", result),
		)
		return nil, 0, err
	}

	_ = u.cache.Set(ctx, cacheKeyCount, []byte(fmt.Sprintf("%d", total)), constant.DefaultExpiration)
	_ = u.cache.Set(ctx, cacheKeyList, jsonList, constant.DefaultExpiration)

	return result, total, nil
}
