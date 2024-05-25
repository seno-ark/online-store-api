package api

import (
	"net/http"
	"online-store/internal/entity"
	"online-store/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// GetListProductByCategory get list of products by a category handler
// @Summary			Get list of Products by Category.
// @Description		Get list of GetListProductByCategory.
// @Tags			Products
// @Produce			json
// @Param			page			query			int	     false	"Pagination page number (default 1, max 500)"				example(1)
// @Param			count			query			int	     false	"Pagination data limit  (default 10, max 100)"				example(10)
// @Success			200 			{object}		utils.Response{data=[]entity.OutGetProduct}
// @Failure			500				{object}		utils.Response
// @Router	/v1/products [get]
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
