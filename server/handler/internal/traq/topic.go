package traq

import (
	"context"
	"errors"
	"fmt"
	"git.trap.jp/Takeno-hito/salmon/server/model"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/go-traq"
	"github.com/traPtitech/traq-ws-bot/payload"
	"gorm.io/gorm"
	"strings"
)

// NewTopic : /topic new [topic]
func (h Handler) NewTopic(p *payload.MessageCreated) {
	channelId := p.Message.ChannelID
	text := p.Message.Text
	// remove "/topic new "
	if strings.Index(text, "/topic new ") == -1 {
		err := h.bot.PostMessage(context.Background(), channelId, "usage: /topic new [topic]")
		if err != nil {
			log.Error(err)
		}
		return
	}
	topic := text[strings.Index(text, "/topic new ")+len("/topic new "):]

	activeTopics, err := h.db.GetActiveTopics()
	if err != nil {
		err := h.bot.PostMessage(context.Background(), channelId, fmt.Sprintf("エラーが発生しました: %v", err))
		if err != nil {
			log.Error(err)
		}
	}

	occupiedChannels := make(map[string]struct{}, len(activeTopics))
	for _, activeTopic := range activeTopics {
		occupiedChannels[activeTopic.ChannelId] = struct{}{}
	}

	targetChannelId := ""
	for _, cid := range model.TopicChannelIds {
		if _, ok := occupiedChannels[cid]; ok {
			continue
		}
		targetChannelId = cid
		break
	}

	if targetChannelId == "" {
		err := h.bot.PostMessage(context.Background(), channelId, "トピックチャンネルが埋まっています。")
		if err != nil {
			log.Error(err)
		}
		return
	}

	messageId, err := h.bot.PostMessageEmbed(context.Background(), targetChannelId, fmt.Sprintf("トピック: %s", topic))

	if _, _, err := h.bot.API().MessageApi.CreatePin(context.Background(), messageId).Execute(); err != nil {
		log.Error(err)
	}

	if err := h.db.CreateTopic(topic, targetChannelId, messageId); err != nil {
		log.Error(err)
	}

	if err := h.bot.PostMessage(context.Background(), channelId, "https://q.trap.jp/messages/"+messageId); err != nil {
		log.Error(err)
	}

	if _, err := h.bot.API().ChannelApi.
		EditChannelTopic(context.Background(), targetChannelId).
		PutChannelTopicRequest(traq.PutChannelTopicRequest{
			Topic: fmt.Sprintf("現在のトピック: %s", topic),
		}).Execute(); err != nil {
		log.Error(err)
	}
}

// GetTopics : /topic list [topic]
func (h Handler) GetTopics(p *payload.MessageCreated) {
	channelId := p.Message.ChannelID

	activeTopics, err := h.db.GetActiveTopics()
	if err != nil {
		err := h.bot.PostMessage(context.Background(), channelId, fmt.Sprintf("エラーが発生しました: %v", err))
		if err != nil {
			log.Error(err)
		}
	}

	topicsByChannelId := make(map[string]model.Topic, len(activeTopics))
	for _, activeTopic := range activeTopics {
		topicsByChannelId[activeTopic.ChannelId] = activeTopic
	}

	outputs := make([]string, 0, len(model.TopicChannelIds))
	for i, cid := range model.TopicChannelIds {
		topic, ok := topicsByChannelId[cid]

		if !ok {
			outputs = append(outputs, fmt.Sprintf("!{\"type\":\"channel\",\"raw\":\"./random/%d\",\"id\":\"%s\"}: no topic", i+1, cid))
		} else {
			outputs = append(outputs, fmt.Sprintf("!{\"type\":\"channel\",\"raw\":\"./random/%d\",\"id\":\"%s\"}: %s", i+1, cid, topic.Topic))
		}
	}

	if _, err := h.bot.PostMessageEmbed(context.Background(), channelId, strings.Join(outputs, "\n\n")); err != nil {
		log.Error(err)
	}
}

// CloseTopic : /topic close
func (h Handler) CloseTopic(p *payload.MessageCreated) {
	channelId := p.Message.ChannelID

	activeTopic, err := h.db.FindActiveTopicByChannelId(channelId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := h.bot.PostMessage(context.Background(), channelId, "このチャンネルで進行中のトピックがありません。")
			if err != nil {
				log.Error(err)
			}
			return
		}
		err := h.bot.PostMessage(context.Background(), channelId, fmt.Sprintf("エラーが発生しました: %v", err))
		if err != nil {
			log.Error(err)
		}
	}

	if err := h.db.ArchiveTopic(activeTopic.Id); err != nil {
		h.bot.PostErrorMessage(context.Background(), channelId, err)
	}

	if err := h.bot.PostMessage(context.Background(), channelId, "アーカイブします。お疲れ様でした！"); err != nil {
		log.Error(err)
	}

	if _, err := h.bot.API().ChannelApi.
		EditChannelTopic(context.Background(), channelId).
		PutChannelTopicRequest(traq.PutChannelTopicRequest{
			Topic: "現在進行中のトピックはありません",
		}).Execute(); err != nil {
		log.Error(err)
	}
}

// RenameTopic : /topic rename [topic]
func (h Handler) RenameTopic(p *payload.MessageCreated) {
	channelId := p.Message.ChannelID

	text := p.Message.Text
	// remove "/topic new "
	if strings.Index(text, "/topic rename ") == -1 {
		err := h.bot.PostMessage(context.Background(), channelId, "usage: /topic rename [topic]")
		if err != nil {
			log.Error(err)
		}
		return
	}
	topic := text[strings.Index(text, "/topic rename ")+len("/topic rename "):]

	activeTopic, err := h.db.FindActiveTopicByChannelId(channelId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := h.bot.PostMessage(context.Background(), channelId, "このチャンネルで進行中のトピックがありません。")
			if err != nil {
				log.Error(err)
			}
			return
		}
		err := h.bot.PostMessage(context.Background(), channelId, fmt.Sprintf("エラーが発生しました: %v", err))
		if err != nil {
			log.Error(err)
		}
	}

	if err := h.db.RenameTopic(activeTopic.Id, topic); err != nil {
		h.bot.PostErrorMessage(context.Background(), channelId, err)
	}

	if _, err := h.bot.API().MessageApi.
		EditMessage(context.Background(), activeTopic.FirstMessageId).
		PostMessageRequest(traq.PostMessageRequest{
			Content: fmt.Sprintf("トピック: %s", topic),
		}).
		Execute(); err != nil {
		log.Error(err)
	}

	if err := h.bot.PostMessage(context.Background(), channelId, "トピック名を変更しました。"); err != nil {
		log.Error(err)
	}

	if _, err := h.bot.API().ChannelApi.
		EditChannelTopic(context.Background(), channelId).
		PutChannelTopicRequest(traq.PutChannelTopicRequest{
			Topic: fmt.Sprintf("現在のトピック: %s", topic),
		}).Execute(); err != nil {
		log.Error(err)
	}
}
