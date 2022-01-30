package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/config"
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/repository"
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/server"
	systemLog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		systemLog.Fatalf("error while getting configuration: %v", err)
	}
	configureLogger(*cfg)

	gophermartService, err := createGophermartService(*cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while creating service")
	}

	errC, err := run(gophermartService)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't run")
	}

	if err = <-errC; err != nil {
		log.Fatal().Err(err).Msg("Error while running")
	}

}

func createGophermartService(cfg config.Config) (*server.GophermartService, error) {
	repo, err := repository.NewGophermartPostgresRepository(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	srv := server.NewService(repo, cfg)
	return srv, nil
}

func run(service *server.GophermartService) (<-chan error, error) {
	errC := make(chan error, 1)

	srv := server.NewServer(service)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		log.Info().Msg("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			stop()
			cancel()
			close(errC)
		}()

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}
		log.Info().Msg("Shutdown completed")

	}()

	go func() {
		log.Info().Msg("Server started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()
	return errC, nil
}

func configureLogger(_ config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// в дальнейшем можно добавить в конфиг требуемый уровень логирования, аутпут (файл или еще чего) и т.д.
	// пока пишем в консоль красивенько
	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
