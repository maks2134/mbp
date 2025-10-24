package db

import (
	"mpb/configs"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Db struct {
	Db *sqlx.DB
}

func NewDb(conf *configs.Config) *Db {
	db, err := sqlx.Connect("postgres", conf.Db.Dsn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	return &Db{db}
}
