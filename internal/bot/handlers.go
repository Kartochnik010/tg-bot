package tgbot

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Kartochnik010/tg-bot/internal/domain/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	infoMessage = `
	Bot created as a test task. Integration with cats API https://thecatapi.com

	Author: @millionaire_go`

	userInfo = `User info:
	TgID: %d
	Firstname: %s
	Lastname: %s
	Username: %s
	Created_at: %s
	`
)

func (b *Bot) start(update tgbotapi.Update) error {
	log := b.l.WithField("chat_id", update.Message.Chat.ID)
	tgID := update.Message.Chat.ID
	user := models.User{
		TgID:      tgID,
		Username:  update.Message.From.UserName,
		Firstname: update.Message.From.FirstName,
		Lastname:  update.Message.From.LastName,
	}
	_, err := b.service.User.SaveUser(b.ctx, user)
	if err != nil {
		log.Error(err.Error())
	}

	msg := tgbotapi.NewMessage(tgID, infoMessage)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/breeds"),
			tgbotapi.NewKeyboardButton("/random"),
		),
	)
	_, err = b.api.Send(msg)
	if err != nil {
		log.WithError(err).Error("failed to send message")
		return err
	}

	return nil
}

func (b *Bot) info(update tgbotapi.Update) error {
	log := b.l.WithField("chat_id", update.Message.Chat.ID)
	tgID := update.Message.Chat.ID

	_, err := b.SendString(tgID, update.Message.Text, update.Message.Chat.ID, infoMessage)
	if err != nil {
		log.WithError(err).Error("failed to send message")
		return err
	}
	return nil
}

func (b *Bot) me(update tgbotapi.Update) error {
	log := b.l.WithField("chat_id", update.Message.Chat.ID)
	tgID := update.Message.Chat.ID

	user, err := b.service.User.GetUserByTgID(b.ctx, tgID)
	if err != nil {
		log.Error(err.Error())
		b.SendString(tgID, update.Message.Text, update.Message.Chat.ID, "Curerntly, we do not store your info, but here is what we know:")
		b.SendString(tgID, update.Message.Text, update.Message.Chat.ID, fmt.Sprintf("User ID: %v\nName: %v\nLastname: %v\nUsername: @%v\n", tgID, update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.UserName))

		return nil
	}
	_, err = b.SendString(tgID, update.Message.Text, update.Message.Chat.ID, fmt.Sprintf(userInfo, user.TgID, user.Firstname, user.Lastname, user.Username, user.CreatedAt))
	if err != nil {
		log.WithError(err).Error("failed to send message")
		return err
	}
	return nil
}

func (b *Bot) breeds(update tgbotapi.Update) error {
	log := b.l.WithField("chat_id", update.Message.Chat.ID)
	tgID := update.Message.Chat.ID

	breeds, err := b.service.CatsAPI.GetAllBreeds(b.ctx)
	if err != nil {
		log.Error(err.Error())
		b.SendString(tgID, update.Message.Text, update.Message.Chat.ID, "failed to get breeds")
		return err
	}

	l := len(breeds)

	b.breedsCounter.mu.Lock()
	i := b.breedsCounter.m[tgID] % l
	b.breedsCounter.m[tgID]++
	b.breedsCounter.mu.Unlock()

	msg := tgbotapi.NewMessage(tgID, fmt.Sprintf("%+v", breeds[i]))
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("next"),
		),
	)
	_, err = b.api.Send(msg)
	if err != nil {
		log.WithError(err).Error("failed to send message")
	}

	return nil
}

func (b *Bot) breed(update tgbotapi.Update) error {
	log := b.l.WithField("chat_id", update.Message.Chat.ID)
	tgID := update.Message.Chat.ID

	_, err := b.SendString(tgID, update.Message.Text, update.Message.Chat.ID, "Enter breed id:")
	if err != nil {
		log.WithError(err).Error("failed to send message")
		return err
	}

	for update := range b.updates {
		if update.Message.Chat.ID != tgID {
			continue
		}
		break
	}

	breedInfo, err := b.service.CatsAPI.GetBreed(b.ctx, update.Message.Text)
	if err != nil {
		log.Error(err.Error())
		b.SendString(tgID, update.Message.Text, update.Message.Chat.ID, err.Error())
		return err
	}
	_, err = b.SendString(tgID, update.Message.Text, update.Message.Chat.ID, fmt.Sprintf("%v", breedInfo))
	if err != nil {
		log.WithError(err).Error("failed to send message")
		return err
	}
	return nil

}

func (b *Bot) random(update tgbotapi.Update) error {
	log := b.l.WithField("chat_id", update.Message.Chat.ID)
	tgID := update.Message.Chat.ID

	catPic, err := b.service.CatsAPI.GetRandomCat(b.ctx)
	if err != nil {
		log.Error(err.Error())
		b.SendString(tgID, update.Message.Text, update.Message.Chat.ID, "failed to get breed info")
		return err
	}

	resp, err := http.Get(catPic.URL)
	if err != nil {
		log.WithError(err).Error("failed to fetch music")
		return err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("failed to read response body")
		return err
	}

	photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{
		Name:  catPic.ID,
		Bytes: bytes,
	})

	_, err = b.Send(tgID, "some cat", photo)
	if err != nil {
		log.WithError(err).Error("failed to send message")
		return err
	}
	return nil

}
