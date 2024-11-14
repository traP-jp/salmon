package traq

import (
	"context"
	"git.trap.jp/Takeno-hito/salmon/server/bot"
	"git.trap.jp/Takeno-hito/salmon/server/model"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/traq-ws-bot/payload"
	"strings"
	"time"
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

func (h Handler) StartVote(p *payload.MessageCreated) {
	channelId := p.Message.ChannelID
	text := p.Message.Text
	// remove "/vote "

	if strings.Index(text, "/vote https://q.trap.jp/messages/") == -1 {
		err := h.bot.PostMessage(context.Background(), channelId, "/vote の後ろにメッセージの URL を付けてください。\n\n"+"https://q.trap.jp/messages/"+p.Message.ID)
		if err != nil {
			log.Error(err)
		}
		return
	}
	targetId := text[strings.Index(text, "/vote https://q.trap.jp/messages/")+len("/vote https://q.trap.jp/messages/"):]

	//msgPlain := "@Takeno_hito \n役員の皆さん投票をお願いします！ 24 時間後に、投票状況に応じて自動で決議されます。投票数が足りなかったらまた来ます！"
	//msgId, err := h.bot.PostMessageEmbed(context.Background(), channelId, msgPlain)
	//if err != nil {
	//	log.Fatal(err)
	//}

	if err := h.bot.AttachVoteStamps(context.Background(), uuid.FromStringOrNil(targetId)); err != nil {
		log.Error(err)
	}

	if err := h.db.CreateScheduledTask(model.JudgeVote, targetId, p.Message.CreatedAt.Add(24*time.Hour)); err != nil {
		log.Error(err)
	}
}
