package entity

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrLoginAlreadyExists = errors.New("login already exists")
	ErrInvalidToken       = errors.New("invalid token")
)
