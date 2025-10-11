package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sqlx.DB
}

func NewDB(sourceName string) (*DB, error) {
	conn, err := sqlx.Connect("sqlite3", sourceName)
	if err != nil {
		return nil, err
	}
	db := &DB{conn}
	return db, nil
}
