package model

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

var TopicChannelIds = []string{
	//"0678c212-bc4d-411c-a347-037b8ef7df46", // for debug
	//"0f9be3e1-96ed-4b4e-9d47-ad4232596bb7", // for debug 2
	"9f9ddc31-5aee-483a-8f7a-13b37b1637a6", // #1
	"7013dc29-7b15-4b45-b2af-89f546b50f70", // #2
	"412f14ba-0730-4e7c-a415-61bae5880d73", // #3
	"ecd3ec82-5f95-47ab-aa34-2f88de1140ee", // #4
	"696b6761-4d5f-4ac8-bf75-e45bfb826303", // #5
	"0cb9bd6d-ec7f-4c9c-804c-dab8759e1dfd", // #6
	"25e1b813-eb94-41ef-ae27-75373c8b7325", // #7
	"18e6c506-28ce-44ea-b4ba-e26629b58f34", // #8
	"01911252-0ba1-7ace-8057-19857d239379", // #9
}

type Topic struct {
	Id             string
	Topic          string
	FirstMessageId string
	ChannelId      string
	CreatedAt      time.Time
	ArchivedAt     sql.NullTime
}

func (c Client) CreateTopic(topic string, channelId string, messageId string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	t := Topic{
		Id:             id.String(),
		Topic:          topic,
		ChannelId:      channelId,
		FirstMessageId: messageId,
		CreatedAt:      time.Now(),
		ArchivedAt:     sql.NullTime{},
	}
	return c.db.Create(&t).Error
}

func (c Client) GetActiveTopics() ([]Topic, error) {
	var topics []Topic
	err := c.db.Where("archived_at IS NULL").Find(&topics).Error
	return topics, err
}

func (c Client) FindActiveTopicByChannelId(channelId string) (Topic, error) {
	var topic Topic
	err := c.db.Where("channel_id = ? AND archived_at IS NULL", channelId).First(&topic).Error
	return topic, err
}

func (c Client) FindActiveTopicById(id string) (Topic, error) {
	var topic Topic
	err := c.db.Where("id = ?", id).First(&topic).Error
	return topic, err
}

func (c Client) RenameTopic(id string, topic string) error {
	return c.db.Model(&Topic{}).Where("id = ?", id).Update("topic", topic).Error
}

func (c Client) ArchiveTopic(id string) error {
	return c.db.Model(&Topic{}).Where("id = ?", id).Update("archived_at", time.Now()).Error
}
