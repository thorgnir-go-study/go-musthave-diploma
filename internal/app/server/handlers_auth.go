package server

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/repository"
	"io"
	"net/http"
)

type registerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *GophermartService) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyContent, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("could not read request body")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var req registerRequest
		if err = json.Unmarshal(bodyContent, &req); err != nil {
			log.Info().Err(err).Msg("invalid json")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		user, err := s.Repository.RegisterUser(r.Context(), repository.UserEntity{Login: req.Login}, req.Password)
		if err != nil {
			if errors.Is(err, repository.ErrUserAlreadyExists) {
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}
			if errors.Is(err, repository.ErrEmptyPassword) {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			log.Error().Err(err).Msg("error creating user")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		token, err := s.JwtWrapper.GenerateToken(user.ID, user.Login)
		if err != nil {
			log.Error().Err(err).Msg("error generating token")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		setJwtCookie(w, token)
		w.WriteHeader(http.StatusOK)

	}
}

func (s *GophermartService) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyContent, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("could not read request body")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var req loginRequest
		if err = json.Unmarshal(bodyContent, &req); err != nil {
			log.Info().Err(err).Msg("invalid json")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		user, err := s.Repository.Authenticate(r.Context(), req.Login, req.Password)
		if err != nil {
			if errors.Is(err, repository.ErrEmptyPassword) {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			if errors.Is(err, repository.ErrAuthenticationFailure) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			log.Error().Err(err).Msg("error authenticating user")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		token, err := s.JwtWrapper.GenerateToken(user.ID, user.Login)

		if err != nil {
			log.Error().Err(err).Msg("error generating token")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		setJwtCookie(w, token)
		w.WriteHeader(http.StatusOK)
	}
}

func setJwtCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}
