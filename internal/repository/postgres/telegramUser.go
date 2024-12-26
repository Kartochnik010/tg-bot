package postgres

import (
	"context"
	"fmt"

	"github.com/Kartochnik010/tg-bot/internal/domain/models"
	"github.com/Kartochnik010/tg-bot/internal/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (l *UserRepo) SaveUser(ctx context.Context, user models.User) (int64, error) {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "UserRepo.SaveUser")

	query := `
		INSERT INTO users (tg_id, firstname, lastname, username)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (tg_id) DO UPDATE SET firstname = EXCLUDED.firstname, lastname = EXCLUDED.lastname, username = EXCLUDED.username
		RETURNING tg_id
	`
	var tgID int64
	err := l.db.QueryRow(ctx, query, user.TgID, user.Firstname, user.Lastname, user.Username).Scan(&tgID)
	if err != nil {
		log.WithError(err).Error("failed to save user")
		return 0, fmt.Errorf("failed to save user: %w", err)
	}

	log.Debugf("saved user: %+v", user)
	return tgID, nil
}

func (c *UserRepo) GetUserByID(ctx context.Context, tgID int64) (*models.User, error) {
	// // // log := logger.GetLoggerFromCtx(ctx).WithField("op", "UserRepo.GetUser")

	query := "select tg_id, firstname, lastname, username, created_at from users where tg_id = $1"
	var user models.User
	err := c.db.QueryRow(ctx, query, tgID).Scan(&user.TgID, &user.Firstname, &user.Lastname, &user.Username, &user.CreatedAt)
	if err != nil {
		// log.WithError(err).Error("failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
