package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/validator"
)

type MetaData struct {
	Key            string    `json: "key"` // Unique identifier for the post (hash)
	CreatedAt      time.Time `json: "created_at"`
	ExpirationDate time.Time `json: "expiration_date"`
}

type PasteModel struct {
	db *sql.DB
}

func (model *PasteModel) Insert(metadata *MetaData) error {
	query := `
        INSERT INTO metadata(key, created_at, expiration_date) 
        VALUES ($1, $2, $3)
        
    `
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, err := model.db.ExecContext(ctx, query, metadata.Key, metadata.CreatedAt, metadata.ExpirationDate)
	if err != nil {
		return err
	}

	return nil
}

func IsMetaDataValid(metadata *MetaData, v *validator.Validator) {
	// Key Validation
	v.Check(len([]rune(metadata.Key)) == 8, "key", "Key should be 8 characters long")

	// CreatedAt Validation
	v.Check(!metadata.CreatedAt.After(time.Now()), "created_at", "Paste creation date cannot be in the future")

	// ExpirationDate Validation
	v.Check(metadata.ExpirationDate.After(time.Now()), "expiration_date", "Expiration date must be in the future")
}
