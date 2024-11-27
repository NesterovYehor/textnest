package models

import "time"

type MetaData struct {
	Key            string    `json: "key"` // Unique identifier for the post (hash)
	CreatedAt      time.Time `json: "created_at"`
	ExpirationDate time.Time `json: "expiration_date"`
}
