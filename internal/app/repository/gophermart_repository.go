package repository

import "context"

type GophermartRepository interface {
	RegisterUser(ctx context.Context, user UserEntity, password string) (UserEntity, error)
	Authenticate(ctx context.Context, login string, password string) (UserEntity, error)
}
