package utils

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailAlredyExist = errors.New("email is used")
)
