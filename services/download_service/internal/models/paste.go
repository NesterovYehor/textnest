package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Paste struct {
	Key         string
	CreatedAt   time.Time
	ExpiredDate time.Time
	Content     string
}

type PasteModel struct {
	DB *sql.DB
}

func (model *PasteModel) Get(key string) (*Paste, error) {
	query := `
        SELECT key, content, created_at, updated_at FROM metadata WHERE key = $1
    `

	// Set up a context with a timeout for the database query
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Prepare a variable to hold the result
	var paste Paste

	// Execute the query and scan the result into the paste struct
	err := model.DB.QueryRowContext(ctx, query, key).Scan(
		&paste.Key,         // Assuming you have a field Key in the Paste struct
		&paste.CreatedAt,   // Assuming you have a field CreatedAt in the Paste struct
		&paste.ExpiredDate, // Assuming you have a field UpdatedAt in the Paste struct
	)
	// Handle any potential errors
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no paste found with the key: %s", key)
		}
		return nil, fmt.Errorf("query failed: %v", err)
	}

	// Return the Paste object
	return &paste, nil
}
