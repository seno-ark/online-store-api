package main

import (
	"fmt"
	"log/slog"
	"online-store/internal/api"
	"online-store/internal/repository"
	"online-store/internal/usecase"
	"online-store/pkg/cache"
	"online-store/pkg/config"
	"online-store/pkg/database"
	"online-store/pkg/token"
	"online-store/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// @title Online-Store API
// @version 1.0
// @description This is a simple api server to simulate customer order system

// @host localhost:9000
// @BasePath
func main() {
	conf := config.GetConfig()

	db, err := database.Postgres(conf)
	if err != nil {
		panic(err)
	}

	cache := cache.NewCache(conf)
	tokenManager := token.NewTokenManager(conf)

	validate := validator.New()
	utils.RegisterCustomValidator(validate)

	repo := repository.NewRepository(db)
	ucase := usecase.NewUsecase(repo, cache)
	handler := api.NewApiHandler(conf, ucase, validate, tokenManager)

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	apiV1 := app.Group("/v1")
	handler.Routes(apiV1)

	slog.Info("Server is starting...", "port", conf.Port)
	if err := app.Listen(fmt.Sprintf(":%v", conf.Port)); err != nil {
		slog.Error("Failed to start the server", err)
	}
}
