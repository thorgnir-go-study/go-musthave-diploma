package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
)

func NewRouter(service *GophermartService) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)

	// public
	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", service.RegisterHandler())
		r.Post("/api/user/login", service.LoginHandler())
	})

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Info().Msg("hello from protected route middleware")
				next.ServeHTTP(w, r)
			})
		})
		r.Get("/api/user/blabla/{x}", func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte(request.RequestURI))
			writer.WriteHeader(200)
		})
	})

	return r
}
