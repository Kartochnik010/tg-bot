package models

import (
	"github.com/Kartochnik010/tg-bot/internal/pkg/lib_time"
	"github.com/google/uuid"
)

type Log struct {
	ID       uuid.UUID `json:"uuid"` // could use uuid
	UserID   int64     `json:"userID"`
	Request  string    `json:"request"`
	Response string    `json:"response"`
}

type User struct {
	TgID      int64            `json:"tgID"`
	Firstname string           `json:"firstname,omitempty"`
	Lastname  string           `json:"lastname,omitempty"`
	Username  string           `json:"username"`
	CreatedAt lib_time.IntDate `json:"created_at"`
}
