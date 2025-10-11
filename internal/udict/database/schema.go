package database

import "fmt"

func (db *DB) InitDB() error {
	schema := `
CREATE TABLE IF NOT EXISTS definitions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    word TEXT NOT NULL,
    definition TEXT NOT NULL,
    thumbs_up INTEGER DEFAULT 0,
    author TEXT NOT NULL,
    written_on TEXT DEFAULT CURRENT_TIMESTAMP,
	created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP 
);`

	if _, err := db.db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	// Create index for case-insensitive word lookups
	createIndexQuery := `
	CREATE INDEX IF NOT EXISTS idx_definitions_word_lower ON definitions (LOWER(word));
	`

	if _, err := db.db.Exec(createIndexQuery); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}
