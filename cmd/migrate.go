package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	_ = godotenv.Load(".env")

	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "postgres://mpb:mpb_pas@localhost:5432/mpb_db?sslmode=disable"
		log.Println("DSN not set, using default")
	}

	db, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := goose.Up(db, "./migrations"); err != nil {
		log.Fatal(err)
	}

	log.Println("Migrations completed successfully")
}
