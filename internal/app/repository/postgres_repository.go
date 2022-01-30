package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type gophermartPostgresRepository struct {
	DB *sqlx.DB
}

var (
	registerUserStmt     *sqlx.NamedStmt
	authenticateUserStmt *sqlx.Stmt
)

func NewGophermartPostgresRepository(_ context.Context, dsn string) (*gophermartPostgresRepository, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = prepareStatements(db); err != nil {
		return nil, err
	}

	return &gophermartPostgresRepository{DB: db}, nil
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func prepareStatements(db *sqlx.DB) error {
	var err error

	if registerUserStmt, err = db.PrepareNamed(`
INSERT INTO gophermart.users (login, password) 
VALUES (:login, crypt(:password, gen_salt('bf'))) 
RETURNING id`); err != nil {
		return err
	}

	if authenticateUserStmt, err = db.Preparex(`
SELECT id, login 
FROM gophermart.users 
WHERE login=$1 AND password = crypt($2, password)`); err != nil {
		return err
	}
	return nil
}

func (g *gophermartPostgresRepository) RegisterUser(ctx context.Context, user UserEntity, password string) (UserEntity, error) {
	if password == "" {
		return UserEntity{}, ErrEmptyPassword
	}

	userWithPwd := UserWithPassword{
		UserEntity: user,
		Password:   password,
	}

	row := registerUserStmt.QueryRowContext(ctx, userWithPwd)
	var userID string
	err := row.Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "login_unique" {
				return UserEntity{}, ErrUserAlreadyExists
			}
		}
		return UserEntity{}, err
	}
	return UserEntity{
		ID:    userID,
		Login: user.Login,
	}, nil
}

func (g *gophermartPostgresRepository) Authenticate(ctx context.Context, login string, password string) (UserEntity, error) {
	if password == "" {
		return UserEntity{}, ErrEmptyPassword
	}
	var entity UserEntity
	err := authenticateUserStmt.GetContext(ctx, &entity, login, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserEntity{}, ErrAuthenticationFailure
		}
		return UserEntity{}, err
	}
	return entity, nil
}
