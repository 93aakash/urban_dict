package database

import (
	"fmt"
	"time"

	"github.com/93aakash/urban_dict/internal/udict/model"
)

func (db *DB) GetDefinitions(query string) ([]model.Definition, error) {
	definitions := []model.Definition{}

	stmt := "SELECT definition, word, author, written_on, thumbs_up FROM definitions WHERE LOWER(word) = ?"
	if err := db.db.Select(&definitions, stmt, query); err != nil {
		return nil, err
	}

	return definitions, nil
}

func (db *DB) InsertDefinitions(defs []model.Definition) error {
	if len(defs) == 0 {
		return nil
	}

	stmt := "INSERT INTO definitions (definition, word, author, written_on, thumbs_up)" +
		"VALUES (:definition, :word, :author, :written_on, :thumbs_up)"

	values := make([]map[string]any, 0, len(defs))

	now := time.Now()

	for _, def := range defs {
		values = append(values, map[string]any{
			"definition": def.Definition,
			"word":       def.Word,
			"author":     def.Author,
			"written_on": def.WrittenOn,
			"thumbs_up":  def.ThumbsUp,
			"created_at": now,
			"updated_at": now,
		})
	}

	_, err := db.db.NamedExec(stmt, values)
	if err != nil {
		return fmt.Errorf("failed to insert definition for '%s': %w", defs[0].Word, err)
	}
	return nil
}

func (db *DB) DeleteDefinition(query string) error {
	stmt := "DELETE FROM definitions WHERE LOWER(word) = ?"
	if _, err := db.db.Exec(stmt, query); err != nil {
		return fmt.Errorf("failed to delete definitions for '%s': %w", query, err)
	}

	return nil
}
