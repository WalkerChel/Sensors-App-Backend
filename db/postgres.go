package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewPostgresDB(host, port, user, dbname, password, sslmode string) (*sqlx.DB, error) {

	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, dbname, password, sslmode))

	return db, err

}
