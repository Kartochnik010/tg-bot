package tgbot

import (
	"context"
	"fmt"
	"log"

	"github.com/Kartochnik010/tg-bot/internal/domain/models"
	"github.com/Kartochnik010/tg-bot/internal/pkg/logger"
	"github.com/Kartochnik010/tg-bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

var (
	UNKNOWN_COMMAND = "Sorry, I didn't understand that. Tap /info for more information"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
	service *service.Service
	l       *logrus.Logger
	router  router
	ctx     context.Context // updated while listening to updates
}

type router struct {
	routes map[string]func(update tgbotapi.Update) error
}

func NewTgRouter() router {
	return router{routes: make(map[string]func(update tgbotapi.Update) error)}
}

func (r *router) addRoute(path string, f func(update tgbotapi.Update) error) {
	r.routes[path] = func(update tgbotapi.Update) error {
		fmt.Println("executing route with path", path)
		return f(update)
	}
}

func (r *router) getFunc(path string) (func(update tgbotapi.Update) error, bool) {
	return r.routes[path], r.routes[path] != nil
}

func New(botKey string, service *service.Service, l *logrus.Logger) *Bot {

	botAPI, err := tgbotapi.NewBotAPI(botKey)
	if err != nil {
		log.Panic(err)
	}
	// botAPI.Debug = true
	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	updatesConfig := tgbotapi.NewUpdate(0)
	updatesConfig.Timeout = 60
	updates := botAPI.GetUpdatesChan(updatesConfig)

	b := &Bot{
		api:     botAPI,
		updates: updates,
		service: service,
		l:       l,
	}
	router := NewTgRouter()
	router.addRoute("/start", b.start)
	router.addRoute("/info", b.info)
	router.addRoute("/me", b.me)
	router.addRoute("/breeds", b.breeds)
	router.addRoute("/breed", b.breed)
	router.addRoute("/random", b.randomPic)

	b.router = router

	return b
}

func (b *Bot) SendString(userID int64, command string, chatID int64, messageText string) (tgbotapi.Message, error) {
	err := b.service.User.LogRequestResponse.SaveLogToDB(b.ctx, models.Log{
		UserID:   userID,
		Request:  command,
		Response: messageText,
	})
	if err != nil {
		return tgbotapi.Message{}, fmt.Errorf("failed to save log: %w", err)
	}

	return b.api.Send(tgbotapi.NewMessage(chatID, messageText))
}

func (b *Bot) Send(userID int64, command string, message any) (tgbotapi.Message, error) {
	switch msg := message.(type) {
	case tgbotapi.MessageConfig:
		err := b.service.User.LogRequestResponse.SaveLogToDB(b.ctx, models.Log{
			UserID:   userID,
			Request:  command,
			Response: msg.Text,
		})
		if err != nil {
			return tgbotapi.Message{}, fmt.Errorf("failed to save log: %w", err)
		}
	case tgbotapi.PhotoConfig:
		err := b.service.User.LogRequestResponse.SaveLogToDB(b.ctx, models.Log{
			UserID:   userID,
			Request:  command,
			Response: "some photo",
		})
		if err != nil {
			return tgbotapi.Message{}, fmt.Errorf("failed to save log: %w", err)
		}
	}

	return b.api.Send(message.(tgbotapi.Chattable))
}

func (b *Bot) ListenUpdates() error {
	for update := range b.updates {
		if update.Message == nil {
			continue
		}
		log := b.l.WithField("chat_id", update.Message.Chat.ID)
		b.ctx = context.WithValue(context.Background(), logger.ContextKeyLogger, b.l.WithField("chat_id", update.Message.Chat.ID))

		log.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
		userID := update.Message.Chat.ID
		messageText := update.Message.Text

		// handle updates
		if f, ok := b.router.routes[update.Message.Text]; ok {
			if err := f(update); err != nil {
				log.WithError(err).Error("failed to handle update")
				continue
			}
			continue
		}

		// UNKNOWN_COMMAND
		if update.Message != nil {
			newMessage := tgbotapi.NewMessage(update.Message.Chat.ID, UNKNOWN_COMMAND)
			newMessage.ReplyToMessageID = update.Message.MessageID

			_, err := b.Send(userID, messageText, newMessage)
			if err != nil {
				log.WithError(err).Error("failed to send message")
				continue
			}
		}
	}
	return nil
}
