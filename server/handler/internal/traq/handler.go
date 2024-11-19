package traq

import (
	"git.trap.jp/Takeno-hito/salmon/server/bot"
	"git.trap.jp/Takeno-hito/salmon/server/model"
)

type Handler struct {
	bot *bot.Bot
	db  *model.Client
}

func New(b *bot.Bot, db *model.Client) Handler {
	return Handler{
		bot: b,
		db:  db,
	}
}
