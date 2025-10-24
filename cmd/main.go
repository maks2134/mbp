package main

import (
	"mpb/configs"
	"mpb/internal/auth"
	"mpb/pkg/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)

	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	// блок репозиториев
	authRepository := auth.NewAuthRepository(conf, database)

	// блок хэндлеров
	authHandler := auth.NewAuthHandlers(authRepository)

	// блок роутов
	api := app.Group("/api")
	authRoutes := auth.NewAuthRoutes(api, authHandler)
	authRoutes.Register()

	err := app.Listen(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
