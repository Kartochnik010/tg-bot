package repository

import (
	"context"

	"github.com/Kartochnik010/tg-bot/internal/domain/models"
	"github.com/Kartochnik010/tg-bot/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User interface {
	SaveUser(ctx context.Context, music models.User) (id int64, err error)
	GetUserByID(ctx context.Context, tgID int64) (*models.User, error)
}
type LogRequestResponse interface {
	SaveLogToDB(ctx context.Context, log models.Log) error
}

type Repository struct {
	User
	LogRequestResponse
}

func NewRepository(db *pgxpool.Pool) Repository {
	return Repository{
		User:               postgres.NewUserRepo(db),
		LogRequestResponse: postgres.NewLogsRepo(db),
	}
}
