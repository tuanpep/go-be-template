package user

import "errors"

var (
	ErrNotFound       = errors.New("user not found")
	ErrDuplicateEmail = errors.New("email already exists")
	ErrInvalidInput   = errors.New("invalid input")
)
