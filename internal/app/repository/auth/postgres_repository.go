package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type authPostgresRepository struct {
	DB *sqlx.DB
}

var (
	registerUserStmt     *sqlx.NamedStmt
	authenticateUserStmt *sqlx.Stmt
)

func NewAuthPostgresRepository(_ context.Context, dsn string) (*authPostgresRepository, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = prepareStatements(db); err != nil {
		return nil, err
	}

	return &authPostgresRepository{DB: db}, nil
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

func (g *authPostgresRepository) RegisterUser(ctx context.Context, user UserDto, password string) (UserDto, error) {
	if password == "" {
		return UserDto{}, ErrEmptyPassword
	}

	userWithPwd := userWithPassword{
		userEntity: userEntity{
			ID:    "",
			Login: user.Login,
		},
		Password: password,
	}

	row := registerUserStmt.QueryRowContext(ctx, userWithPwd)
	var userID string
	err := row.Scan(&userID)
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
	err := authenticateUserStmt.GetContext(ctx, &entity, login, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserDto{}, ErrAuthenticationFailure
		}
		return UserDto{}, err
	}
	return UserDto{
		UserID: entity.ID,
		Login:  entity.Login,
	}, nil
}
