package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectToDB(conf *sqlx.Conn) *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=mads password=1234 database=test_unexpected_behavior sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	db.MustExec(schema)
	return db
}
