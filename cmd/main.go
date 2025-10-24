package cmd

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
