package repository

import (
	"context"
	"miborchestrator/internal/entities"
	customerrors "miborchestrator/internal/entities/custom_errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthPostgres struct {
	db *pgxpool.Pool
}

func NewAuthPostgress(db *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, user entities.User) error {

	_, err := r.db.Exec(ctx, "insert into users_table (login,password_hash) values ($1,$2);", user.Login, user.Password)
	if err != nil {
		pqErr := err.(*pgconn.PgError)
		if pqErr.Code == pgerrcode.UniqueViolation {
			return &customerrors.LoginConflict{Elem: user.Login}
		}
		return err
	}

	return nil
}

func (r *AuthPostgres) GetUser(ctx context.Context, login, password string) (int, error) {
	var userID int
	if err := r.db.QueryRow(ctx, "select id from users_table where login=$1 and password_hash=$2", login, password).Scan(&userID); err != nil {

		if err == pgx.ErrNoRows {
			return 0, &customerrors.InvalidLoginOrPassword{}
		}
		return 0, err

	}
	return userID, nil
}
