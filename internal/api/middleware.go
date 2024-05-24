package api

import (
	"net/http"
	"online-store/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *ApiHandler) AuthBearerMiddleware(c *fiber.Ctx) error {
	resp := utils.NewResponse()

	authHeader := strings.Split(c.Get("Authorization", ""), " ")
	if len(authHeader) != 2 {
		c.Status(http.StatusUnauthorized)
		return c.JSON(resp.Set("Unauthorized", nil))
	}

	authType := strings.Trim(authHeader[0], " ")
	authToken := authHeader[1]

	if strings.ToLower(authType) != "bearer" || authToken == "" {
		c.Status(http.StatusUnauthorized)
		return c.JSON(resp.Set("Invalid Token", nil))
	}

	claims, err := h.tokenManager.ValidateJwtToken(authToken)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return c.JSON(resp.Set(strings.ToTitle(err.Error()), nil))
	}

	c.Locals("claims", claims)

	return c.Next()
}

func (h *ApiHandler) AuthApiKeyMiddleware(c *fiber.Ctx) error {
	// dummy auth

	resp := utils.NewResponse()

	authKey := c.Get("X-API-KEY", "")
	if authKey != h.conf.WebhookApiKey {
		c.Status(http.StatusUnauthorized)
		return c.JSON(resp.Set("Invalid API Key", nil))
	}

	return c.Next()
}
