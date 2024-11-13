package handler

import (
	"git.trap.jp/Takeno-hito/salmon/server/bot"
	"git.trap.jp/Takeno-hito/salmon/server/database"
	"git.trap.jp/Takeno-hito/salmon/server/handler/traq"
	"github.com/traPtitech/traq-ws-bot/payload"
)

type Handler struct {
	traQHandler traq.Handler
	db          *database.Client
	bot         *bot.Bot
}

func New(b *bot.Bot, db *database.Client) Handler {
	return Handler{
		traQHandler: traq.New(b, db),
		db:          db,
		bot:         b,
	}
}

func (h Handler) HandleBotMessage(p *payload.MessageCreated) {
	msg := p.Message.PlainText
	if msg == "/vote" || msg == "@BOT_no_hito_local /vote" {
		h.traQHandler.StartVote(p)
	}
}
