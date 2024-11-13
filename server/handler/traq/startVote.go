package traq

import (
	"context"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/traP-jp/salmon/server/bot"
	"github.com/traP-jp/salmon/server/database"
	"github.com/traPtitech/traq-ws-bot/payload"
	"time"
)

type Handler struct {
	bot *bot.Bot
	db  *database.Client
}

func New(b *bot.Bot, db *database.Client) Handler {
	return Handler{
		bot: b,
		db:  db,
	}
}

func (h Handler) StartVote(p *payload.MessageCreated) {
	channelId := p.Message.ChannelID
	msgPlain := "@Takeno_hito \n役員の皆さん投票をお願いします！ 24 時間後に、投票状況に応じて自動で決議されます。投票数が足りなかったらまた来ます！"
	msgId, err := h.bot.PostMessageEmbed(context.Background(), channelId, msgPlain)
	if err != nil {
		log.Fatal(err)
	}

	if err := h.bot.AttachVoteStamps(context.Background(), uuid.FromStringOrNil(msgId)); err != nil {
		log.Fatal(err)
	}

	if err := h.db.CreateScheduledTask("judge-vote", msgId, p.Message.CreatedAt.Add(24*time.Hour)); err != nil {
		log.Fatal(err)
	}
}
