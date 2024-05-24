package api

import (
	"net/http"
	"online-store/internal/entity"
	"online-store/pkg/token"
	"online-store/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

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
		c.SendStatus(http.StatusBadRequest)
		return c.JSON(resp.Set("Invalid data", nil).AddErrValidation(errs))
	}

	result, err := h.ucase.CreateCartItem(ctx, claims.UserID, *req)
	if err != nil {
		status, msg := utils.ErrStatusCode(err)
		c.SendStatus(status)
		return c.JSON(resp.Set(msg, nil))
	}

	return c.JSON(resp.Set("success", result))
}

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
		c.SendStatus(status)
		return c.JSON(resp.Set(msg, nil))
	}

	return c.JSON(resp.Set("success", nil))
}

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
		c.SendStatus(status)
		return c.JSON(resp.Set(msg, nil))
	}

	resp.AddMeta(page, count, total)
	return c.JSON(resp.Set("success", results))
}
