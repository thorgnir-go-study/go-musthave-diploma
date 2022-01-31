package auth

import (
	"context"
)

type Repository interface {
	RegisterUser(ctx context.Context, user UserDto, password string) (UserDto, error)
	Authenticate(ctx context.Context, login string, password string) (UserDto, error)
}
