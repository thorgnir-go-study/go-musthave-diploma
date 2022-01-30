package server

import (
	"net/http"
)

func NewServer(service *GophermartService) *http.Server {
	router := NewRouter(service)
	server := &http.Server{
		Addr:    service.Config.ServerAddress,
		Handler: router,
	}

	return server
}
