package api

import (
	"net/http"
	"online-store/internal/entity"
	"online-store/pkg/token"
	"online-store/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AddCartItem add product to cart handler
// @Summary			Add Cart Item.
// @Description		Add product to cart.
// @Tags			Carts
// @Accept			json
// @Produce			json
// @Param			Authorization	header	string	true	"Bearer token"
// @Param 			json			body	entity.InCreateCartItem	true	"Cart Item data"
// @Success			201		{object}	utils.Response{data=entity.CartItem}
// @Failure			400		{object}	utils.Response
// @Failure			500		{object}	utils.Response
// @Router	/v1/carts [post]
func (h *ApiHandler) AddCartItem(c *fiber.Ctx) error {
	ctx := c.Context()
	resp := utils.NewResponse()

	claims := c.Locals("claims").(*token.Claims)
	if claims == nil {
		c.SendStatus(http.StatusUnauthorized)
		return c.JSON(resp.Set("Unauthorized", nil))
	}

	req := new(entity.InCreateCartItem)
	if err := c.BodyParser(req); err != nil {
		c.SendStatus(http.StatusBadRequest)
		return c.JSON(resp.Set("Invalid data", nil))
	}

	err := h.validate.Struct(req)
	if err != nil {
		errs := utils.ParseValidatorErr(err)
		c.Status(http.StatusBadRequest)
		return c.JSON(resp.Set("Invalid data", nil).AddErrValidation(errs))
	}

	result, err := h.ucase.CreateCartItem(ctx, claims.UserID, *req)
	if err != nil {
		status, msg := utils.ErrStatusCode(err)
		return c.Status(status).JSON(resp.Set(msg, nil))
	}

	return c.Status(http.StatusCreated).JSON(resp.Set("success", result))
}

// DeleteCartItem delete product from cart handler
// @Summary			Delete Cart Item.
// @Description		Delete product from cart.
// @Tags			Carts
// @Param			Authorization	header	string	true	"Bearer token"
// @Param			cart_id			path	string true		"Cart ID" example(02a1a6a3-1c9c-4f46-ae18-162e2b0d7a9a)
// @Produce			json
// @Success			200 			{object}		utils.Response
// @Failure			404				{object}		utils.Response
// @Failure			500				{object}		utils.Response
// @Router	/v1/carts/{cart_id} [delete]
func (h *ApiHandler) DeleteCartItem(c *fiber.Ctx) error {
	ctx := c.Context()
	resp := utils.NewResponse()

	claims := c.Locals("claims").(*token.Claims)
	if claims == nil {
		c.SendStatus(http.StatusUnauthorized)
		return c.JSON(resp.Set("Unauthorized", nil))
	}

	cartID := c.Params("cart_id")
	err := h.validate.Var(cartID, "uuid")
	if err != nil {
		errs := utils.ParseValidatorErr(err)
		c.SendStatus(http.StatusBadRequest)
		return c.JSON(resp.Set("Invalid data", nil).AddErrValidation(errs))
	}

	err = h.ucase.DeleteCartItem(ctx, claims.UserID, cartID)
	if err != nil {
		status, msg := utils.ErrStatusCode(err)
		return c.Status(status).JSON(resp.Set(msg, nil))
	}

	return c.JSON(resp.Set("success", nil))
}

// GetListCartItem get list of cart items handler
// @Summary			Get list of Cart Items.
// @Description		Get list of Cart Items.
// @Tags			Carts
// @Produce			json
// @Param			Authorization	header	string	true	"Bearer token"
// @Param			page			query			int	     false	"Pagination page number (default 1, max 500)"				example(1)
// @Param			count			query			int	     false	"Pagination data limit  (default 10, max 100)"				example(10)
// @Success			200 			{object}		utils.Response{data=[]entity.CartItem}
// @Failure			500				{object}		utils.Response
// @Router	/v1/carts [get]
func (h *ApiHandler) GetListCartItem(c *fiber.Ctx) error {
	ctx := c.Context()
	resp := utils.NewResponse()

	claims := c.Locals("claims").(*token.Claims)
	if claims == nil {
		c.SendStatus(http.StatusUnauthorized)
		return c.JSON(resp.Set("Unauthorized", nil))
	}

	queries := c.Queries()
	page, count := utils.Pagination(queries["page"], queries["count"])

	results, total, err := h.ucase.GetListCartItem(ctx, claims.UserID, entity.InGetListCartItem{
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
