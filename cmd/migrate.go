package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	db, err := goose.OpenDBWithDriver("postgres", "postgres://mpb:mpb_pas@localhost:5432/mpb_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := goose.Up(db, "./migrations"); err != nil {
		log.Fatal(err)
	}
}
