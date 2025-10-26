package main

import (
	"log"
	"mpb/configs"
	"mpb/internal/auth"
	"mpb/pkg/db"

	"github.com/gofiber/fiber/v2"
)

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

	log.Fatal(app.Listen(":8000"))
}
