package api

import (
	"online-store/internal/entity"
	"online-store/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// GetListCategory get list of categories handler
// @Summary			Get list of Categories.
// @Description		Get list of Categories.
// @Tags			Categories
// @Produce			json
// @Param			page			query			int	     false	"Pagination page number (default 1, max 500)"				example(1)
// @Param			count			query			int	     false	"Pagination data limit  (default 10, max 100)"				example(10)
// @Success			200 			{object}		utils.Response{data=[]entity.Category}
// @Failure			500				{object}		utils.Response
// @Router	/v1/categories [get]
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
		return c.Status(status).JSON(resp.Set(msg, nil))
	}

	resp.AddMeta(page, count, total)
	return c.JSON(resp.Set("success", results))
}
