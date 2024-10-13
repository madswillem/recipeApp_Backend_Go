package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectToDB(conf *sqlx.Conn, connection_string string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", connection_string)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
