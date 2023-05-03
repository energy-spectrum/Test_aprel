package util

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidPassword = errors.New("invalid password")
)
