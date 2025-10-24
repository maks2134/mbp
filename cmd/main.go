package main

import (
	"fmt"
	"mpb/configs"
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

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
