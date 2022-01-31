package server

import (
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/server/auth"
	"net/http"
)

func NewServer(service *auth.Service) *http.Server {
	router := NewRouter(service)
	server := &http.Server{
		Addr:    service.Config.ServerAddress,
		Handler: router,
	}

	return server
}
