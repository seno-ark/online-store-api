package api

import (
	"net/http"
	"online-store/internal/entity"
	"online-store/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *ApiHandler) GetListProductByCategory(c *fiber.Ctx) error {
	ctx := c.Context()
	resp := utils.NewResponse()

	categoryID := c.Params("category_id")
	err := h.validate.Var(categoryID, "uuid")
	if err != nil {
		errs := utils.ParseValidatorErr(err)
		c.SendStatus(http.StatusBadRequest)
		return c.JSON(resp.Set("Invalid data", nil).AddErrValidation(errs))
	}

	queries := c.Queries()
	page, count := utils.Pagination(queries["page"], queries["count"])

	results, total, err := h.ucase.GetListProduct(ctx, entity.InGetListProduct{
		CategoryID: categoryID,
		Limit:      count,
		Offset:     (page - 1) * count,
	})
	if err != nil {
		status, msg := utils.ErrStatusCode(err)
		c.SendStatus(status)
		return c.JSON(resp.Set(msg, nil))
	}

	resp.AddMeta(page, count, total)
	return c.JSON(resp.Set("success", results))
}
