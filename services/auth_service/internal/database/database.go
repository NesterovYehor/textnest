package database

import (
	"database/sql"
	"fmt"

	"github.com/NesterovYehor/textnest/services/auth_service/config"
	_ "github.com/lib/pq"
)

type DB struct {
	Conn  *sql.DB
	Close func() error
}

func New(cfg *config.DBConfig) (*DB, error) {
	db, err := sql.Open("postgres", cfg.Link)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Ensure database is reachable
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxOpenConns(cfg.MaxOpenConns) // Corrected here
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return &DB{Conn: db, Close: db.Close}, nil
}
