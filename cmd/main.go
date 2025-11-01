package main

import (
	"log"
	"mpb/configs"
	"mpb/internal/app"
)

func main() {
	conf := configs.LoadConfig()

	a, err := app.New(conf)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	if err := a.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
