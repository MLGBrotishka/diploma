package entity

import "time"

type User struct {
	ID        int64
	Login     string
	PassHash  []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool
}
