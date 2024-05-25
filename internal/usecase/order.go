package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"online-store/internal/entity"
	"online-store/internal/repository"
	"online-store/pkg/constant"
	"online-store/pkg/utils"
	"strconv"
)

func (u *Usecase) CreateOrder(ctx context.Context, userID string, arg entity.InCreateOrder) (*entity.Order, error) {
	type productChartMapper struct {
		Price  int64
		CartID string
	}
	mapProductChart := map[string]productChartMapper{}
	var totalCost int64

	for _, item := range arg.Items {
		product, err := u.repo.GetProduct(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		if product.Stock < item.Qty {
			return nil, utils.NewErrUnprocessable(fmt.Sprintf("Product is out of stock: %s", product.Name))
		}

		cartItem, err := u.repo.GetCartItemByProduct(ctx, userID, item.ProductID)
		if err != nil && !errors.Is(err, utils.ErrNotFound) {
			return nil, err
		}

		productChart := productChartMapper{
			Price: product.Price,
		}
		if cartItem != nil {
			productChart.CartID = cartItem.ID
		}
		mapProductChart[item.ProductID] = productChart

		totalCost += item.Qty * product.Price
	}

	var orderID string

	err := u.repo.WithTx(ctx, func(rtx *repository.Repository) error {
		var errTx error

		// create order
		orderID, errTx = rtx.CreateOrder(ctx, entity.Order{
			UserID:          userID,
			ShipmentAddress: arg.ShipmentAddress,
			TotalCost:       totalCost,
		})
		if errTx != nil {
			return errTx
		}

		// create order payment
		_, errTx = rtx.CreateOrderPayment(ctx, entity.Payment{
			OrderID:         orderID,
			PaymentMethod:   arg.Payment.PaymentMethod,
			PaymentProvider: arg.Payment.PaymentProvider,
			BillAmount:      totalCost,
			Status:          string(constant.PaymentStatusPending),
		})
		if errTx != nil {
			return errTx
		}

		for _, item := range arg.Items {

			// create order items
			errTx = rtx.CreateOrderItem(ctx, entity.OrderItem{
				OrderID:      orderID,
				ProductID:    item.ProductID,
				Qty:          item.Qty,
				ProductPrice: mapProductChart[item.ProductID].Price,
				Notes:        item.Notes,
			})
			if errTx != nil {
				return errTx
			}

			// update product stock
			errTx = rtx.UpdateProductStock(ctx, item.ProductID, item.Qty*-1)
			if errTx != nil {
				return errTx
			}

			// delete related cart
			if mapProductChart[item.ProductID].CartID != "" {
				errTx = rtx.DeleteCartItem(ctx, mapProductChart[item.ProductID].CartID)
				if errTx != nil {
					return errTx
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	_ = u.cache.Del(ctx, fmt.Sprintf(constant.CacheKeyListUserCartAll, userID))
	_ = u.cache.Del(ctx, fmt.Sprintf(constant.CacheKeyListProductByCategoryAll))
	_ = u.cache.Del(ctx, fmt.Sprintf(constant.CacheKeyListUserOrderAll, userID))

	return u.repo.GetOrder(ctx, orderID)
}

func (u *Usecase) GetOrder(ctx context.Context, userID, orderID string) (*entity.Order, error) {
	return u.repo.GetUserOrder(ctx, userID, orderID)
}

func (u *Usecase) GetListOrder(ctx context.Context, userID string, arg entity.InGetListOrder) ([]entity.Order, int64, error) {
	var (
		result []entity.Order
		total  int64
		err    error
	)

	cacheKeyCount := fmt.Sprintf(constant.CacheKeyListUserOrderCount, userID)
	cacheKeyList := fmt.Sprintf(constant.CacheKeyListUserOrder, userID, arg.Limit, arg.Offset)

	cacheData, _ := u.cache.Get(ctx, cacheKeyCount)
	if cacheData != "" {
		total, err = strconv.ParseInt(cacheData, 10, 64)
		if err != nil {
			slog.Error(
				"Failed to GetListOrder Parse Cache",
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
					"Failed to GetListOrder Parse Cache",
					slog.Any("err", err),
					slog.Any("cacheKey", cacheKeyList),
					slog.Any("cacheData", cacheData),
				)
				return nil, 0, utils.NewErrInternalServer("Failed to parse cache data")
			}

			return result, total, nil
		}
	}

	result, total, err = u.repo.GetListOrder(ctx, userID, arg)
	if err != nil {
		return nil, 0, err
	}

	jsonList, err := json.Marshal(result)
	if err != nil {
		slog.Error(
			"Failed to GetListOrder Marshal Cache",
			slog.Any("err", err),
			slog.Any("data", result),
		)
		return nil, 0, err
	}

	_ = u.cache.Set(ctx, cacheKeyCount, []byte(fmt.Sprintf("%d", total)), constant.DefaultExpiration)
	_ = u.cache.Set(ctx, cacheKeyList, jsonList, constant.DefaultExpiration)

	return result, total, nil
}

func (u *Usecase) GetListOrderItem(ctx context.Context, userID string, orderID string, arg entity.InGetListOrderItem) ([]entity.OutGetOrderItem, int64, error) {
	_, err := u.repo.GetUserOrder(ctx, userID, orderID)
	if err != nil {
		return nil, 0, err
	}

	var (
		result []entity.OutGetOrderItem
		total  int64
	)

	cacheKeyCount := fmt.Sprintf(constant.CacheKeyListOrderItemCount, orderID)
	cacheKeyList := fmt.Sprintf(constant.CacheKeyListOrderItem, orderID, arg.Limit, arg.Offset)

	cacheData, _ := u.cache.Get(ctx, cacheKeyCount)
	if cacheData != "" {
		total, err = strconv.ParseInt(cacheData, 10, 64)
		if err != nil {
			slog.Error(
				"Failed to GetListOrderItem Parse Cache",
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
					"Failed to GetListOrderItem Parse Cache",
					slog.Any("err", err),
					slog.Any("cacheKey", cacheKeyList),
					slog.Any("cacheData", cacheData),
				)
				return nil, 0, utils.NewErrInternalServer("Failed to parse cache data")
			}

			return result, total, nil
		}
	}

	result, total, err = u.repo.GetListOrderItem(ctx, orderID, arg)
	if err != nil {
		return nil, 0, err
	}

	jsonList, err := json.Marshal(result)
	if err != nil {
		slog.Error(
			"Failed to GetListOrderItem Marshal Cache",
			slog.Any("err", err),
			slog.Any("data", result),
		)
		return nil, 0, err
	}

	_ = u.cache.Set(ctx, cacheKeyCount, []byte(fmt.Sprintf("%d", total)), constant.DefaultExpiration)
	_ = u.cache.Set(ctx, cacheKeyList, jsonList, constant.DefaultExpiration)

	return result, total, nil
}
