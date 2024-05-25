package api

import (
	"net/http"
	"online-store/internal/entity"
	"online-store/pkg/token"
	"online-store/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// CreateOrder create order handler
// @Summary			Create Order.
// @Description		Create new Order.
// @Tags			Orders
// @Accept			json
// @Produce			json
// @Param			Authorization	header	string	true	"Bearer token"
// @Param 			json	body		entity.InCreateOrder	true	"Order data"
// @Success			201		{object}	utils.Response{data=entity.Order}
// @Failure			400		{object}	utils.Response
// @Failure			500		{object}	utils.Response
// @Router	/v1/orders [post]
func (h *ApiHandler) CreateOrder(c *fiber.Ctx) error {
	ctx := c.Context()
	resp := utils.NewResponse()

	claims := c.Locals("claims").(*token.Claims)
	if claims == nil {
		c.SendStatus(http.StatusUnauthorized)
		return c.JSON(resp.Set("Unauthorized", nil))
	}

	req := new(entity.InCreateOrder)
	if err := c.BodyParser(req); err != nil {
		c.SendStatus(http.StatusBadRequest)
		return c.JSON(resp.Set("Invalid data", nil))
	}

	err := h.validate.Struct(req)
	if err != nil {
		errs := utils.ParseValidatorErr(err)
		c.SendStatus(http.StatusBadRequest)
		return c.JSON(resp.Set("Invalid data", nil).AddErrValidation(errs))
	}

	err = h.validate.Struct(req.Payment)
	if err != nil {
		errs := utils.ParseValidatorErr(err)
		c.SendStatus(http.StatusBadRequest)
		return c.JSON(resp.Set("Invalid data", nil).AddErrValidation(errs))
	}

	for _, v := range req.Items {
		err = h.validate.Struct(v)
		if err != nil {
			errs := utils.ParseValidatorErr(err)
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(resp.Set("Invalid data", nil).AddErrValidation(errs))
		}
	}

	result, err := h.ucase.CreateOrder(ctx, claims.UserID, *req)
	if err != nil {
		status, msg := utils.ErrStatusCode(err)
		return c.Status(status).JSON(resp.Set(msg, nil))
	}

	c.Status(http.StatusCreated)
	return c.JSON(resp.Set("success", result))
}

// GetOrder get user order handler
// @Summary			Get user Order.
// @Description		Get user Order.
// @Tags			Orders
// @Param			Authorization	header	string	true	"Bearer token"
// @Param			order_id		path			string	 true	"Order ID"
// @Produce			json
// @Success			200 			{object}		utils.Response{data=entity.Order}
// @Failure			404				{object}		utils.Response
// @Failure			500				{object}		utils.Response
// @Router	/v1/orders/{order_id} [get]
func (h *ApiHandler) GetOrder(c *fiber.Ctx) error {
	ctx := c.Context()
	resp := utils.NewResponse()

	claims := c.Locals("claims").(*token.Claims)
	if claims == nil {
		c.SendStatus(http.StatusUnauthorized)
		return c.JSON(resp.Set("Unauthorized", nil))
	}

	orderID := c.Params("order_id")
	err := h.validate.Var(orderID, "uuid")
	if err != nil {
		errs := utils.ParseValidatorErr(err)
		c.SendStatus(http.StatusBadRequest)
		return c.JSON(resp.Set("Invalid data", nil).AddErrValidation(errs))
	}

	orderDetail, err := h.ucase.GetOrder(ctx, claims.UserID, orderID)
	if err != nil {
		status, msg := utils.ErrStatusCode(err)
		return c.Status(status).JSON(resp.Set(msg, nil))
	}

	orderItemList, _, err := h.ucase.GetListOrderItem(ctx, claims.UserID, orderID, entity.InGetListOrderItem{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		status, msg := utils.ErrStatusCode(err)
		return c.Status(status).JSON(resp.Set(msg, nil))
	}

	paymentDetail, err := h.ucase.GetOrderPayment(ctx, orderID)
	if err != nil {
		status, msg := utils.ErrStatusCode(err)
		return c.Status(status).JSON(resp.Set(msg, nil))
	}

	return c.JSON(resp.Set("success", fiber.Map{
		"order":   orderDetail,
		"payment": paymentDetail,
		"items":   orderItemList,
	}))
}

// GetListOrder get list of user orders handler
// @Summary			Get list of user Orders.
// @Description		Get list of user Orders..
// @Tags			Orders
// @Produce			json
// @Param			Authorization	header	string	true	"Bearer token"
// @Param			page			query			int	     false	"Pagination page number (default 1, max 500)"				example(1)
// @Param			count			query			int	     false	"Pagination data limit  (default 10, max 100)"				example(10)
// @Success			200 			{object}		utils.Response{data=[]entity.Order}
// @Failure			500				{object}		utils.Response
// @Router	/v1/orders [get]
func (h *ApiHandler) GetListOrder(c *fiber.Ctx) error {
	ctx := c.Context()
	resp := utils.NewResponse()

	claims := c.Locals("claims").(*token.Claims)
	if claims == nil {
		c.SendStatus(http.StatusUnauthorized)
		return c.JSON(resp.Set("Unauthorized", nil))
	}

	queries := c.Queries()
	page, count := utils.Pagination(queries["page"], queries["count"])

	results, total, err := h.ucase.GetListOrder(ctx, claims.UserID, entity.InGetListOrder{
		Limit:  count,
		Offset: (page - 1) * count,
	})
	if err != nil {
		status, msg := utils.ErrStatusCode(err)
		return c.Status(status).JSON(resp.Set(msg, nil))
	}

	resp.AddMeta(page, count, total)
	return c.JSON(resp.Set("success", results))
}
