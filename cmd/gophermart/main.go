package main

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/config"
	authRepo "github.com/thorgnir-go-study/go-musthave-diploma/internal/app/repository/auth"
	"github.com/thorgnir-go-study/go-musthave-diploma/internal/app/server"
	authService "github.com/thorgnir-go-study/go-musthave-diploma/internal/app/server/auth"

	shopspring "github.com/jackc/pgtype/ext/shopspring-numeric"

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

	dbpool, err := createDbPool(*cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while creating dbpool")
	}
	authSrvc, err := createAuthService(dbpool, *cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while creating auth service")
	}

	srv := server.NewServer(authSrvc)

	errC, err := run(srv)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't run")
	}

	if err = <-errC; err != nil {
		log.Fatal().Err(err).Msg("Error while running")
	}

}

func createDbPool(cfg config.Config) (*pgxpool.Pool, error) {
	dbconfig, err := pgxpool.ParseConfig(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &shopspring.Numeric{},
			Name:  "numeric",
			OID:   pgtype.NumericOID,
		})
		return nil
	}
	dbpool, err := pgxpool.ConnectConfig(context.Background(), dbconfig)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}

func createAuthService(dbpool *pgxpool.Pool, cfg config.Config) (*authService.Service, error) {
	authRepository, err := authRepo.NewAuthPostgresRepository(dbpool)
	if err != nil {
		return nil, err
	}

	srv := authService.New(authRepository, cfg)
	return srv, nil
}

func run(srv *http.Server) (<-chan error, error) {
	errC := make(chan error, 1)

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
