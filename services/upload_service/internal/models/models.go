package models

import "time"

type MetaData struct {
	Key            string
	CreatedAt      time.Time
	ExpirationDate time.Time
}
