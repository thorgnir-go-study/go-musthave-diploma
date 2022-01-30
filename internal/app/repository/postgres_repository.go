package repository

import (
	"context"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type gophermartPostgresRepository struct {
	DB *sqlx.DB
}

func (g gophermartPostgresRepository) Dummy() error {
	//TODO implement me
	panic("implement me")
}

func NewGophermartPostgresRepository(ctx context.Context, dsn string) (*gophermartPostgresRepository, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &gophermartPostgresRepository{DB: db}, nil
}
