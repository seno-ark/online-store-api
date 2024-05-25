package api

import (
	"net/http"
	"online-store/internal/entity"
	"online-store/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// PaymentWebhook payment notification webhook handler
// @Summary			Payment Webhook.
// @Description		Payment Webhook from Payment Gateway.
// @Tags			Payments
// @Accept			json
// @Produce			json
// @Param			X-API-KEY	header	string	true	"API Key"
// @Param 			json	body		entity.InPaymentWebHook	true	"Payment Webhook data"
// @Success			201		{object}	utils.Response
// @Failure			400		{object}	utils.Response
// @Failure			500		{object}	utils.Response
// @Router	/v1/payments [post]
func (h *ApiHandler) PaymentWebhook(c *fiber.Ctx) error {
	ctx := c.Context()
	resp := utils.NewResponse()

	req := new(entity.InPaymentWebHook)
	if err := c.BodyParser(req); err != nil {
		c.SendStatus(http.StatusBadRequest)
		return c.JSON(resp.Set("Invalid data", nil))
	}

	err := h.validate.Struct(req)
	if err != nil {
		errs := utils.ParseValidatorErr(err)
		c.SendStatus(http.StatusBadRequest)
		return c.JSON(resp.Set("Invalid payload data", nil).AddErrValidation(errs))
	}
	err = h.validate.Struct(req.PaymentWebHookTransactionDetails)
	if err != nil {
		errs := utils.ParseValidatorErr(err)
		c.SendStatus(http.StatusBadRequest)
		return c.JSON(resp.Set("Invalid transaction data", nil).AddErrValidation(errs))
	}
	err = h.validate.Struct(req.PaymentWebHookUserDetails)
	if err != nil {
		errs := utils.ParseValidatorErr(err)
		c.SendStatus(http.StatusBadRequest)
		return c.JSON(resp.Set("Invalid user data", nil).AddErrValidation(errs))
	}

	err = h.ucase.UpdateOrderPayment(ctx, *req)
	if err != nil {
		status, msg := utils.ErrStatusCode(err)
		c.SendStatus(status)
		return c.JSON(resp.Set(msg, nil))
	}

	return c.JSON(resp.Set("success", nil))
}
