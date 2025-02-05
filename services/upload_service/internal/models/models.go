package models

import "time"

type MetaData struct {
	Key            string
	UserId         string
	CreatedAt      time.Time
	ExpirationDate time.Time
}
