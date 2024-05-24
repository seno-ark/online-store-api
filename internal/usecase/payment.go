package usecase

import (
	"context"
	"encoding/json"
	"online-store/internal/entity"
	"online-store/internal/repository"
	"online-store/pkg/constant"
	"online-store/pkg/utils"
	"time"
)

func (u *Usecase) UpdateOrderPayment(ctx context.Context, arg entity.InPaymentWebHook) error {
	order, err := u.repo.GetOrder(ctx, arg.PaymentWebHookTransactionDetails.OrderID)
	if err != nil {
		return err
	}

	if order.Status != string(constant.OrderStatusPaymentPending) {
		return nil
	}

	user, err := u.repo.GetUser(ctx, order.UserID)
	if err != nil {
		return err
	}

	webhookPayload, err := json.Marshal(arg)
	if err != nil {
		return utils.NewErrInternalServer("invalid json payload")
	}

	err = u.repo.WithTx(ctx, func(rtx *repository.Repository) error {
		var errTx error

		paidAt := time.Now().UTC()
		log := string(webhookPayload)

		// update payment
		errTx = rtx.UpdateOrderPayment(ctx, entity.Payment{
			OrderID:       order.ID,
			PaidAmount:    arg.PaymentAmount,
			Status:        arg.Status,
			TransactionID: arg.TransactionID,
			PaidAt:        &paidAt,
			Log:           &log,
		})
		if errTx != nil {
			return errTx
		}

		paymentComplete := arg.Status == string(constant.PaymentStatusSettlement)
		validPaymentAccount := arg.PaymentWebHookUserDetails.FullName == user.FullName
		validPaymentAmount := arg.PaymentAmount >= order.TotalCost

		if paymentComplete && validPaymentAccount && validPaymentAmount {

			// update order status
			errTx = rtx.UpdateOrderStatus(ctx, order.ID, constant.OrderStatusPaid)
			if errTx != nil {
				return errTx
			}
		}

		return nil
	})

	return err
}

func (u *Usecase) GetOrderPayment(ctx context.Context, id string) (*entity.Payment, error) {
	return u.repo.GetOrderPayment(ctx, id)
}
