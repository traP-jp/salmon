package handler

import (
	"database/sql"
	"fmt"
	"git.trap.jp/Takeno-hito/salmon/server/bot"
	"git.trap.jp/Takeno-hito/salmon/server/handler/internal/traq"
	"git.trap.jp/Takeno-hito/salmon/server/model"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/traq-ws-bot/payload"
	"strings"
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
	// msg contain @BOT_salmon /vote

	if h.bot.Env() == bot.EnvProduction {
		if strings.Index(msg, "/vote") == 0 || strings.Index(msg, "@BOT_salmon /vote") == 0 {
			h.traQHandler.StartVote(p)
		} else if strings.Index(msg, "/topic new") == 0 || strings.Index(msg, "@BOT_salmon /topic new") == 0 {
			h.traQHandler.NewTopic(p)
		} else if strings.Index(msg, "/topic list") == 0 || strings.Index(msg, "@BOT_salmon /topic list") == 0 {
			h.traQHandler.GetTopics(p)
		} else if strings.Index(msg, "/topic close") == 0 || strings.Index(msg, "@BOT_salmon /topic close") == 0 {
			h.traQHandler.CloseTopic(p)
		} else if strings.Index(msg, "/topic rename") == 0 || strings.Index(msg, "@BOT_salmon /topic rename") == 0 {
			h.traQHandler.RenameTopic(p)
		}
	} else {
		fmt.Println(msg)
		if strings.Index(msg, "@BOT_no_hito_local /vote") == 0 {
			h.traQHandler.StartVote(p)
		} else if strings.Index(msg, "/topic new") == 0 {
			h.traQHandler.NewTopic(p)
		} else if strings.Index(msg, "/topic list") == 0 {
			h.traQHandler.GetTopics(p)
		} else if strings.Index(msg, "/topic close") == 0 {
			h.traQHandler.CloseTopic(p)
		} else if strings.Index(msg, "/topic rename") == 0 {
			h.traQHandler.RenameTopic(p)
		}
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
