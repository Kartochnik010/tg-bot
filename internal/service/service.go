package service

import (
	"net/http"

	"github.com/Kartochnik010/tg-bot/internal/config"
	"github.com/Kartochnik010/tg-bot/internal/repository"
)

type Service struct {
	User    *UserService
	CatsAPI *CatsApiService
}

func NewService(repo repository.Repository, cfg *config.Config, c *http.Client) *Service {
	return &Service{
		User:    NewUserService(c, cfg, repo),
		CatsAPI: NewCatsApiService(c),
	}
}
