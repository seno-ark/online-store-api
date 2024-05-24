package api

import (
	"online-store/internal/usecase"
	"online-store/pkg/cache"
	"online-store/pkg/token"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ApiHandler struct {
	ucase        *usecase.Usecase
	validate     *validator.Validate
	cache        *cache.Cache
	tokenManager *token.TokenManager
}

func NewApiHandler(
	ucase *usecase.Usecase,
	validate *validator.Validate,
	cache *cache.Cache,
	tokenManager *token.TokenManager,
) *ApiHandler {
	return &ApiHandler{
		ucase:        ucase,
		validate:     validate,
		cache:        cache,
		tokenManager: tokenManager,
	}
}

func (h *ApiHandler) Routes(apiV1 fiber.Router) {
	apiV1.Get("/health", h.healthCheck)

	users := apiV1.Group("users")

	users.Post("/register", h.Register)
	users.Post("/login", h.Login)
	users.Get("/me", h.AuthMiddleware, h.GetLoggedInUser)
}

func (h *ApiHandler) healthCheck(c *fiber.Ctx) error {
	return c.SendString("I'm fine, thanks.")
}
