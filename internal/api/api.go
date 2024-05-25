package api

import (
	"online-store/internal/usecase"
	"online-store/pkg/config"
	"online-store/pkg/token"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	_ "online-store/internal/api/docs"

	"github.com/gofiber/swagger"
)

type ApiHandler struct {
	conf         *config.Config
	ucase        *usecase.Usecase
	validate     *validator.Validate
	tokenManager *token.TokenManager
}

func NewApiHandler(
	conf *config.Config,
	ucase *usecase.Usecase,
	validate *validator.Validate,
	tokenManager *token.TokenManager,
) *ApiHandler {
	return &ApiHandler{
		conf:         conf,
		ucase:        ucase,
		validate:     validate,
		tokenManager: tokenManager,
	}
}

func (h *ApiHandler) Routes(apiV1 fiber.Router) {
	apiV1.Get("/health", h.healthCheck)
	apiV1.Get("/swagger/*", swagger.HandlerDefault)

	apiV1.Get("/dummy-data", h.DummyProductAndCategory)

	users := apiV1.Group("users")
	users.Post("/register", h.Register)
	users.Post("/login", h.Login)
	users.Get("/me", h.AuthBearerMiddleware, h.GetLoggedInUser)

	categories := apiV1.Group("categories")
	categories.Get("/", h.GetListCategory)

	products := apiV1.Group("products")
	products.Get("/category/:category_id", h.GetListProductByCategory)

	carts := apiV1.Group("carts")
	carts.Use(h.AuthBearerMiddleware)
	carts.Post("/", h.AddCartItem)
	carts.Get("/", h.GetListCartItem)
	carts.Delete("/:cart_id", h.DeleteCartItem)

	orders := apiV1.Group("orders")
	orders.Use(h.AuthBearerMiddleware)
	orders.Post("/", h.CreateOrder)
	orders.Get("/", h.GetListOrder)
	orders.Get("/:order_id", h.GetOrder)

	payments := apiV1.Group("payments")
	payments.Use(h.AuthApiKeyMiddleware)
	payments.Post("/webhook", h.PaymentWebhook)
}

func (h *ApiHandler) healthCheck(c *fiber.Ctx) error {
	return c.SendString("I'm fine, thanks.")
}
