package bot

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/go-traq"
	"github.com/traPtitech/traq-ws-bot/payload"
)

// エラーメッセージを柔軟に返却させるために、エラーはここでハンドリングしない
func (b *Bot) joinOrLeaveHandler(p *payload.MessageCreated) {
	m := p.Message
	log.Debug("Received MESSAGE_CREATED event: " + m.Text + " / " + m.PlainText)

	if b.env == EnvProduction {
		if m.PlainText == "@BOT_no_hito きて" {
			b.joinChannel(m)
			return
		}

		if m.PlainText == "@BOT_no_hito でてって" {
			b.leaveChannel(m)
			return
		}
	} else {
		if m.PlainText == "@BOT_no_hito_local きて2" {
			b.joinChannel(m)
			return
		}

		if m.PlainText == "@BOT_no_hito_local でてって2" {
			b.leaveChannel(m)
			return
		}
	}
}

func (b *Bot) joinChannel(m payload.Message) {
	if m.User.Name != "Takeno_hito" {
		err := b.PostMessage(context.Background(), m.ChannelID, ":dare:")
		if err != nil {
			log.Error(err)
		}
		return
	}

	_, err := b.API().BotApi.
		LetBotJoinChannel(context.Background(), b.botID).PostBotActionJoinRequest(traq.PostBotActionJoinRequest{
		ChannelId: m.ChannelID,
	}).Execute()

	if err != nil {
		log.Error(err)
		_ = b.PostMessage(context.Background(), m.ChannelID, "なんか参加できなかったかも")
		return
	}

	err = b.PostMessage(context.Background(), m.ChannelID, ":trasta_general.large:")
	return
}

func (b *Bot) leaveChannel(m payload.Message) {
	if m.User.Name != "Takeno_hito" {
		err := b.PostMessage(context.Background(), m.ChannelID, ":dare:")
		if err != nil {
			log.Error(err)
		}
		return
	}

	_, err := b.API().BotApi.
		LetBotLeaveChannel(context.Background(), b.botID).PostBotActionLeaveRequest(traq.PostBotActionLeaveRequest{
		ChannelId: m.ChannelID,
	}).Execute()

	if err != nil {
		log.Error(err)
		_ = b.PostMessage(context.Background(), m.ChannelID, "なんか退出できなかったかも")
		return
	}

	err = b.PostMessage(context.Background(), m.ChannelID, ":gomen.large: ばいばい…")
	return
}
