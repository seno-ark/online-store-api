package api

import (
	"net/http"
	"online-store/internal/entity"
	"online-store/pkg/token"
	"online-store/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Register create user handler
// @Summary			Create User.
// @Description		Create new User.
// @Tags			Users
// @Accept			json
// @Produce			json
// @Param 			json	body		entity.InUserRegister	true	"User data"
// @Success			201		{object}	utils.Response
// @Failure			400		{object}	utils.Response
// @Failure			500		{object}	utils.Response
// @Router	/v1/users/register [post]
func (h *ApiHandler) Register(c *fiber.Ctx) error {
	ctx := c.Context()
	resp := utils.NewResponse()

	req := new(entity.InUserRegister)
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

	result, err := h.ucase.Register(ctx, entity.InUserRegister{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	})
	if err != nil {
		status, msg := utils.ErrStatusCode(err)
		return c.Status(status).JSON(resp.Set(msg, nil))
	}

	accessToken, err := h.tokenManager.GenerateJwtToken(token.Claims{
		UserID: result.ID,
	})
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return c.JSON(resp.Set("Failed to login", nil))
	}

	c.Status(http.StatusCreated)
	return c.JSON(resp.Set("success", fiber.Map{
		"user":         result,
		"access_token": accessToken,
	}))
}

// Login login handler
// @Summary			User Login.
// @Description		User Login.
// @Tags			Users
// @Accept			json
// @Produce			json
// @Param 			json	body		entity.InUserLogin	true	"Login data"
// @Success			200		{object}	utils.Response
// @Failure			400		{object}	utils.Response
// @Failure			500		{object}	utils.Response
// @Router	/v1/users/login [post]
func (h *ApiHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()
	resp := utils.NewResponse()

	req := new(entity.InUserLogin)
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

	result, err := h.ucase.Login(ctx, *req)
	if err != nil {
		status, msg := utils.ErrStatusCode(err)
		return c.Status(status).JSON(resp.Set(msg, nil))
	}

	accessToken, err := h.tokenManager.GenerateJwtToken(token.Claims{
		UserID: result.ID,
	})
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return c.JSON(resp.Set("Failed to login", nil))
	}

	return c.JSON(resp.Set("success", fiber.Map{
		"user":         result,
		"access_token": accessToken,
	}))
}

// GetLoggedInUser get active user handler
// @Summary			Get Logged in user profile.
// @Description		Get active user profile.
// @Tags			Users
// @Param			Authorization	header	string	true	"Bearer token"
// @Produce			json
// @Success			200 			{object}		utils.Response
// @Failure			404				{object}		utils.Response
// @Failure			500				{object}		utils.Response
// @Router	/v1/users/me [get]
func (h *ApiHandler) GetLoggedInUser(c *fiber.Ctx) error {
	ctx := c.Context()
	resp := utils.NewResponse()

	claims := c.Locals("claims").(*token.Claims)
	if claims == nil {
		c.SendStatus(http.StatusUnauthorized)
		return c.JSON(resp.Set("Unauthorized", nil))
	}

	result, err := h.ucase.GetUser(ctx, claims.UserID)
	if err != nil {
		status, msg := utils.ErrStatusCode(err)
		return c.Status(status).JSON(resp.Set(msg, nil))
	}

	return c.JSON(resp.Set("success", fiber.Map{
		"id":         result.ID,
		"email":      result.Email,
		"full_name":  result.FullName,
		"created_at": result.CreatedAt,
		"updated_at": result.UpdatedAt,
	}))
}
