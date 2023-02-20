package repository

import (
	"context"
	"miborchestrator/internal/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Authorization interface {
	CreateUser(ctx context.Context, user entities.User) error
	GetUser(ctx context.Context, login, password string) (int, error)
}

// type TxLog interface {
// 	Record(userId int, tx entities.Transaction) (int, error)
// 	GetById(userId, txID string) (todo.TodoList, error)
// 	Delete(userId, listId int) error
// 	Update(userId, listId int, input todo.UpdateListInput) error
// }

type Repository struct {
	Authorization
	// TxLog
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Authorization: NewAuthPostgress(db),
	}
}
