package service

import (
	"context"
	"net/http"

	"github.com/Kartochnik010/tg-bot/internal/config"
	"github.com/Kartochnik010/tg-bot/internal/domain/models"
	"github.com/Kartochnik010/tg-bot/internal/repository"
)

type UserService struct {
	c    *http.Client
	cfg  *config.Config
	repo repository.Repository
	repository.LogRequestResponse
}

func NewUserService(client *http.Client, cfg *config.Config, repo repository.Repository) *UserService {
	return &UserService{
		c:                  client,
		cfg:                cfg,
		repo:               repo,
		LogRequestResponse: repo.LogRequestResponse,
	}
}

func (s *UserService) SaveUser(ctx context.Context, user models.User) (int64, error) {
	// log := logger.GetLoggerFromCtx(ctx).WithField("op", "UserService.SaveUser")
	// log.Debug()

	id, err := s.repo.User.SaveUser(ctx, user)
	if err != nil {
		// log.WithError(err).Error("failed to store music")
		return 0, err
	}
	return id, nil
}

func (s *UserService) GetUserByTgID(ctx context.Context, tgID int64) (*models.User, error) {
	// log := logger.GetLoggerFromCtx(ctx).WithField("op", "UserService.SaveUser")
	// log.Debug()

	user, err := s.repo.User.GetUserByID(ctx, tgID)
	if err != nil {
		// log.WithError(err).Error("failed to store music")
		return nil, err
	}
	return user, nil
}
