package server

import (
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/config"
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/repository"
)

type GophermartService struct {
	Repository repository.GophermartRepository
	Config     config.Config
}

func NewService(repository repository.GophermartRepository, config config.Config) *GophermartService {
	s := &GophermartService{
		Repository: repository,
		Config:     config,
	}

	return s
}
