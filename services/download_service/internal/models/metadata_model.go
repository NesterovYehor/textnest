package models

import (
	"time"
)

type Metadata struct {
	Key         string
	CreatedAt   time.Time
	ExpiredDate time.Time
}

