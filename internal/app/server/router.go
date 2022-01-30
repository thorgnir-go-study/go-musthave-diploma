package server

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(service *GophermartService) chi.Router {
	r := chi.NewRouter()

	r.Get("/dummy", service.DummyHandler())

	return r
}
