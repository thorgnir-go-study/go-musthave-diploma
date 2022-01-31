package auth

import (
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/auth"
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/config"
	authRepo "github.com/thorgnir-go-study/go-musthave-diploma/internal/app/repository/auth"
)

type Service struct {
	AuthRepository authRepo.Repository
	Config         config.Config
	JwtWrapper     *auth.JwtWrapper
}

func New(authRepo authRepo.Repository, config config.Config) *Service {
	jwtWrapper := auth.NewJwtWrapper(config.JWTSecret, "gophermart", 24)

	s := &Service{
		AuthRepository: authRepo,
		Config:         config,
		JwtWrapper:     jwtWrapper,
	}

	return s
}
