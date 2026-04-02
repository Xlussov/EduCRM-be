package domain

import "errors"

var (
	ErrAlreadyExists      = errors.New("already exists")
	ErrNotFound           = errors.New("not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInternal           = errors.New("internal error")
	ErrPhoneAlreadyExists = errors.New("phone already exists")
)
