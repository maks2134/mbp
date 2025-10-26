package db

import (
	"mpb/configs"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Db struct {
	Conn *sqlx.DB
}

func NewDb(conf *configs.Config) *Db {
	db, err := sqlx.Connect("postgres", conf.Db.Dsn)
	if err != nil {
		panic(err)
	}

	return &Db{Conn: db}
}
