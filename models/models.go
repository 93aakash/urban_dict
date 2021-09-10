package models

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sqlx.DB
}

type Def struct {
	Definition string `json:"definition" db:"definition"`
	ThumbsUp   int    `json:"thumbs_up" db:"thumbs_up"`
	Author     string `json:"author" db:"author"`
	Word       string `json:"word" db:"word"`
	WrittenOn  string `json:"written_on" db:"written_on"`
}

func InitDB(driverName string, sourceName string) (*DB, error) {
	conn, err := sqlx.Connect(driverName, sourceName)
	if err != nil {
		return nil, err
	}
	db := &DB{conn}
	return db, nil
}

func (db *DB) IfExists(query string) bool {
	var result string
	err := db.Get(&result, "select word from dict where lower(word) = ? limit 1", query)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Fatalln(err)
	}
	return true
}

func (db *DB) FetchDef(query string) ([]Def, error) {
	result := []Def{}
	stmt := "select definition, word, author, written_on, thumbs_up from dict where lower(word) = ?"
	err := db.Select(&result, stmt, query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *DB) InsertDef(defs []Def) error {
	stmt := "insert into dict (definition, word, author, written_on, thumbs_up)" +
		"values (:definition, :word, :author, :written_on, :thumbs_up)"
	_, err := db.NamedExec(stmt, defs)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteDef(query string) error {
	if !db.IfExists(query) {
		return fmt.Errorf("\"%s\" does not exist in the database", query)
	}
	stmt := "delete from dict where lower(word) = ?"
	db.MustExec(stmt, query)
	fmt.Printf("Deleted entries for \"%s\" from the database\n", query)
	return nil
}
