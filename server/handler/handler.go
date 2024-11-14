package handler

import (
	"database/sql"
	"git.trap.jp/Takeno-hito/salmon/server/bot"
	"git.trap.jp/Takeno-hito/salmon/server/handler/internal/traq"
	"git.trap.jp/Takeno-hito/salmon/server/model"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/traq-ws-bot/payload"
	"time"
)

type Handler struct {
	traQHandler traq.Handler
	db          *model.Client
	bot         *bot.Bot
}

func New(b *bot.Bot, db *model.Client) Handler {
	return Handler{
		traQHandler: traq.New(b, db),
		db:          db,
		bot:         b,
	}
}

func (h Handler) TraQMessageHandler(p *payload.MessageCreated) {
	msg := p.Message.PlainText
	if msg == "/vote" || msg == "@BOT_no_hito_local /vote" {
		h.traQHandler.StartVote(p)
	}
}

// TaskConsumeHandler consumes scheduled tasks
func (h Handler) TaskConsumeHandler() {
	tasks, err := h.db.GetActiveScheduledTasks()
	if err != nil {
		log.Error(err)
	}

	for _, task := range tasks {

		err = h.db.UpdateScheduledTask(model.ScheduledTask{
			Id:          task.Id,
			Command:     task.Command,
			Arg:         task.Arg,
			ScheduledAt: task.ScheduledAt,
			CreatedAt:   task.CreatedAt,
			ExecutedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		})
		if err != nil {
			log.Error(err)
		}

		switch task.Command {
		case model.JudgeVote:
			if err = judge(h.bot, task.Arg); err != nil {
				log.Error(err)
			}
		default:
			log.Errorf("unknown command: %s", task.Command)
		}
	}
}
