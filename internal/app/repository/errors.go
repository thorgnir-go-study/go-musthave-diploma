package repository

import "errors"

var (
	ErrUserAlreadyExists     = errors.New("user with specified login already exists")
	ErrEmptyPassword         = errors.New("password can not be empty")
	ErrAuthenticationFailure = errors.New("user with specified login and password not found")
)
