package api

import (
	"online-store/internal/entity"
	"online-store/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *ApiHandler) GetListCategory(c *fiber.Ctx) error {
	ctx := c.Context()
	resp := utils.NewResponse()

	queries := c.Queries()
	page, count := utils.Pagination(queries["page"], queries["count"])

	results, total, err := h.ucase.GetListCategory(ctx, entity.InGetListCategory{
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
