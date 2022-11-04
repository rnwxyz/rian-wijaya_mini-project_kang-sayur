package utils

import "errors"

var (
	ErrUserNotFound                 = errors.New("user not found")
	ErrEmailAlredyExist             = errors.New("email is used")
	ErrBadRequestBody               = errors.New("bad request body")
	ErrInvalidId                    = errors.New("invalid id")
	ErrNotAllowedDeleteDefaultAdmin = errors.New("default admin user not allowed to delete")
	ErrFailedDeleteUser             = errors.New("failed delete user")
	ErrInvalidPassword              = errors.New("invalid password")
	ErrPermission                   = errors.New("not have permission to access")
)