package tgbot

import (
	"fmt"

	"github.com/Kartochnik010/tg-bot/internal/domain/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	infoMessage = `
	Bot created as a test task. Integration with cats API https://api.thecatapi.com
	
	Command list:
	/start - start bot
	/info - show this message
	/me - user info
	/breeds - all breeds
	/breed - breed info
	/random - random cat picture

	Author: @millionaire_go
	`

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
		return err
	}
	fmt.Println("(tgID, update.Message.Text, update.Message.Chat.ID, infoMessage)")
	fmt.Println(tgID, update.Message.Text, update.Message.Chat.ID, infoMessage)
	_, err = b.SendString(tgID, update.Message.Text, update.Message.Chat.ID, infoMessage)
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
		b.SendString(tgID, update.Message.Text, update.Message.Chat.ID, "failed to get user")
		return err
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
	res := "List of all avalable breeds:\n"
	for _, breed := range breeds {
		res += breed + "\n"
	}
	_, err = b.SendString(tgID, update.Message.Text, update.Message.Chat.ID, res)
	if err != nil {
		log.WithError(err).Error("failed to send message")
		return err
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

func (b *Bot) randomPic(update tgbotapi.Update) error {
	log := b.l.WithField("chat_id", update.Message.Chat.ID)
	tgID := update.Message.Chat.ID

	bytes, err := b.service.CatsAPI.GetRandomCatPicture(b.ctx)
	if err != nil {
		log.Error(err.Error())
		b.SendString(tgID, update.Message.Text, update.Message.Chat.ID, "failed to get breed info")
		return err
	}

	photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{
		Name:  "cat",
		Bytes: bytes,
	})

	_, err = b.Send(tgID, "random cat picture", photo)
	if err != nil {
		log.WithError(err).Error("failed to send message")
		return err
	}
	return nil

}
