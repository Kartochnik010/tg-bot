package postgres

import (
	"context"
	"fmt"

	"github.com/Kartochnik010/tg-bot/internal/domain/models"
	"github.com/google/uuid"

	// "github.com/Kartochnik010/tg-bot/internal/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LogsRepo struct {
	db *pgxpool.Pool
}

func NewLogsRepo(db *pgxpool.Pool) *LogsRepo {
	return &LogsRepo{
		db: db,
	}
}

func (l *LogsRepo) SaveLogToDB(ctx context.Context, logRequestResponse models.Log) error {
	// log := logger.GetLoggerFromCtx(ctx).WithField("op", "LogsRepo.SaveLogToDB")
	query := `
		INSERT INTO log_request_response (id, user_id, request, response)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	_, err := l.db.Exec(ctx, query, uuid.New(), logRequestResponse.UserID, logRequestResponse.Request, logRequestResponse.Response)
	if err != nil {
		// log.WithError(err).Error("failed to save log to db")
		return fmt.Errorf("failed to save log to db: %w", err)
	}

	// log.Debugf("saved log to postgres: %+v", logRequestResponse)
	return nil
}
