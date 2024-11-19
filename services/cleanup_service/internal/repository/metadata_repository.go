package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/validator"
)

type metadataRepository struct {
	DB *sql.DB // Database connection
}

func NewMetadataRepository(db *sql.DB) MetadataRepository {
	return &metadataRepository{
		DB: db,
	}
}

func (repo *metadataRepository) DeleteAndReturnExpiredKeys() ([]string, error) {
	query := `DELETE FROM metadata WHERE expiration_time <= $1 RETURNING key; `
	var keys []string

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	rows, err := repo.DB.QueryContext(ctx, query, time.Now())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, err
}

func (repo *metadataRepository) DeletePasteByKey(key string) error {
	query := `  DELETE FROM metadata WHERE key = $1`
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, query, key)
	if err != nil {
		return err
	}
	return nil
}

func (repo *metadataRepository) IsKeyValid(v *validator.Validator, key string) {
	v.Check(len([]rune(key)) != 8, "key", "Provided key must be 8 chars lenth")
}
