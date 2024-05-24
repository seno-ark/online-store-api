package usecase

import (
	"context"
	"errors"
	"fmt"
	"online-store/internal/entity"
	"online-store/internal/repository"
	"online-store/pkg/utils"
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

	return u.repo.GetOrder(ctx, userID, orderID)
}

func (u *Usecase) GetOrder(ctx context.Context, userID, orderID string) (*entity.Order, error) {
	return u.repo.GetOrder(ctx, userID, orderID)
}

func (u *Usecase) GetListOrder(ctx context.Context, userID string, arg entity.InGetListOrder) ([]entity.Order, int64, error) {
	return u.repo.GetListOrder(ctx, userID, arg)
}

func (u *Usecase) GetListOrderItem(ctx context.Context, userID string, orderID string, arg entity.InGetListOrderItem) ([]entity.OutGetOrderItem, int64, error) {
	_, err := u.repo.GetOrder(ctx, userID, orderID)
	if err != nil {
		return nil, 0, err
	}

	return u.repo.GetListOrderItem(ctx, orderID, arg)
}
