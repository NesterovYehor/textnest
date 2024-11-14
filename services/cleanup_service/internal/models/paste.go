package models

import "time"

type Paste struct {
	Key            string
	CreatedAt      time.Time
	ExpirationDate time.Time
}
