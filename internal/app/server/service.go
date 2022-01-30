package server

import (
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/auth"
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/config"
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/repository"
)

type GophermartService struct {
	Repository repository.GophermartRepository
	Config     config.Config
	JwtWrapper *auth.JwtWrapper
}

func NewService(repository repository.GophermartRepository, config config.Config) *GophermartService {
	jwtWrapper := auth.NewJwtWrapper(config.JWTSecret, "gophermart", 24)

	s := &GophermartService{
		Repository: repository,
		Config:     config,
		JwtWrapper: jwtWrapper,
	}

	return s
}
