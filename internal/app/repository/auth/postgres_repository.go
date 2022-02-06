package auth

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type authPostgresRepository struct {
	dbpool *pgxpool.Pool
}

func NewAuthPostgresRepository(dbpool *pgxpool.Pool) (*authPostgresRepository, error) {
	return &authPostgresRepository{dbpool: dbpool}, nil
}

func (g *authPostgresRepository) RegisterUser(ctx context.Context, user UserDto, password string) (UserDto, error) {
	if password == "" {
		return UserDto{}, ErrEmptyPassword
	}

	var userID string
	err := g.dbpool.QueryRow(ctx, `INSERT INTO gophermart.users (login, password) 
VALUES ($1, crypt($2, gen_salt('bf'))) 
RETURNING id`, user.Login, password).Scan(&userID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "login_unique" {
				return UserDto{}, ErrUserAlreadyExists
			}
		}
		return UserDto{}, err
	}
	return UserDto{
		UserID: userID,
		Login:  user.Login,
	}, nil
}

func (g *authPostgresRepository) Authenticate(ctx context.Context, login string, password string) (UserDto, error) {
	if password == "" {
		return UserDto{}, ErrEmptyPassword
	}
	var entity userEntity

	if err := pgxscan.Get(ctx, g.dbpool, &entity, `SELECT id, login
FROM gophermart.users
WHERE login=$1 AND password = crypt($2, password)`, login, password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserDto{}, ErrAuthenticationFailure
		}
		return UserDto{}, err
	}

	return UserDto{
		UserID: entity.ID,
		Login:  entity.Login,
	}, nil
}
