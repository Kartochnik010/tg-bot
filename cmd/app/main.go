package main

import (
	l "log"
	"net/http"
	"os"

	"github.com/Kartochnik010/tg-bot/internal/app"
	"github.com/Kartochnik010/tg-bot/internal/config"
	"github.com/Kartochnik010/tg-bot/internal/pkg/logger"
	"github.com/Kartochnik010/tg-bot/internal/repository/postgres"
)

// @title Currency API
// @version 1.0
// @description This is a simple currency API
// @host localhost:8080
// @BasePath /
// @schemes http
// @produce json
// @consumes json
// @contact.name Telegram
// @contact.url https://t.me/ilyas_amantaev
func main() {
	cfg, err := config.ReadConfigFile("config.json")
	if err != nil {
		l.Println("failed to read config file")
		return
	}

	log, err := logger.NewLogger(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		l.Println("failed to create logger:", err)
		os.Exit(1)
	}

	db, err := postgres.New(cfg.DSN)
	if err != nil {
		log.WithError(err).Error("failed to create db")
		os.Exit(1)
	}

	app := app.NewApp(cfg, db, log)
	if err != nil {
		log.WithError(err).Error("failed to create app")
		os.Exit(1)
	}

	if err := app.Run(); err != nil && err != http.ErrServerClosed {
		log.WithError(err).Error("failed to run app")
		os.Exit(1)
	}
}
