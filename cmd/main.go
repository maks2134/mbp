package main

import (
	"mpb/internal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	//conf := configs.LoadConfig()
	//database := db.NewDb(conf)

	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	api := app.Group("/api")
	authRoutes := auth.NewAuthRoutes(api)
	authRoutes.Register()

	err := app.Listen(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
