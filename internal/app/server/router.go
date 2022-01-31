package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/middlewares"
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
		r.Use(middlewares.JwtAuthMiddleware(service.JwtWrapper))

		r.Get("/api/user/blabla/{x}", func(writer http.ResponseWriter, request *http.Request) {
			claims, err := middlewares.GetClaimsFromContext(request.Context())
			if err != nil {
				log.Info().Err(err).Msg("Error while getting claims from context")
			}
			writer.Write([]byte(claims.Login))
			writer.WriteHeader(200)
		})
	})

	return r
}
