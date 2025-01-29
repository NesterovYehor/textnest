package models

import "time"

type MetaData struct {
	Key            string
	UserId         int64
	CreatedAt      time.Time
	ExpirationDate time.Time
}
