package models

import "time"

type MetaData struct {
	Key            string
	Title          string
	UserId         string
	CreatedAt      time.Time
	ExpirationDate time.Time
}
