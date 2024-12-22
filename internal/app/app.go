package app

import (
	"net/http"
	"time"

	tgbot "github.com/Kartochnik010/tg-bot/internal/bot"
	"github.com/Kartochnik010/tg-bot/internal/config"
	"github.com/Kartochnik010/tg-bot/internal/repository"
	"github.com/Kartochnik010/tg-bot/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type App struct {
	Logger *logrus.Logger
	Repo   repository.Repository
	TgBot  *tgbot.Bot
}

func NewApp(cfg *config.Config, db *pgxpool.Pool, l *logrus.Logger) *App {
	const op = "app.NewApp"

	repo := repository.NewRepository(db)

	client := &http.Client{
		Timeout: time.Second * cfg.HTTPClient.RequestTimeout,
	}

	s := service.NewService(repo, cfg, client)

	bot := tgbot.New(cfg.BotKey, s, l)
	return &App{
		TgBot:  bot,
		Logger: l,
		Repo:   repo,
	}
}

func (a *App) Run() error {
	return a.TgBot.ListenUpdates()
}
