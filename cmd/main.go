package main

import (
	"mpb/configs"
	_ "mpb/docs"
	"mpb/internal/auth"
	"mpb/pkg/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title MPB Blog Auth API
// @version 1.0
// @description Authentication service for MPB blog platform
// @host localhost:8000
// @BasePath /api
func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)

	app := fiber.New()

	api := app.Group("/api")

	// auth блок
	authRepo := auth.NewAuthRepository(conf, database)
	authService := auth.NewAuthService(authRepo, []byte(conf.JWT.SecretKey), conf.JWT.AccessTokenTTL)
	authHandler := auth.NewAuthHandlers(authService)
	authRoutes := auth.NewAuthRoutes(api, authHandler)
	authRoutes.Register()

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err)
	}
}
