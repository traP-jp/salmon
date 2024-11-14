package handler

import (
	"context"
	"fmt"
	"git.trap.jp/Takeno-hito/salmon/server/bot"
	"github.com/gofrs/uuid"
	"github.com/traPtitech/go-traq"
	"strings"
	"time"
)

func judge(b *bot.Bot, messageId string) error {
	msg, err := b.GetMessageFromMessageId(context.Background(), uuid.FromStringOrNil(messageId))
	if err != nil {
		return err
	}

	agreedUsersId := make([]string, 0)
	disagreedUsersId := make([]string, 0)

	groupMembers, _, err := b.API().GroupApi.GetUserGroupMembers(context.Background(), bot.ExecutiveGroupId).Execute()
	groupMembersSet := make(map[string]struct{})
	for _, member := range groupMembers {
		groupMembersSet[member.Id] = struct{}{}
	}

	if err != nil {
		return err
	}

	for _, stamp := range msg.Stamps {
		if _, ok := groupMembersSet[stamp.UserId]; !ok {
			continue
		}

		switch stamp.StampId {
		case bot.AgreeStampId:
			agreedUsersId = append(agreedUsersId, stamp.UserId)
		case bot.DisagreeStampId:
			disagreedUsersId = append(disagreedUsersId, stamp.UserId)
		}
	}

	messages := make([]string, 0)

	if len(agreedUsersId) < (len(groupMembers)*6)/10 {
		messages = append(messages, "投票数が足りないみたいです…。まだの方は投票をお願いします！")
	}

	if len(disagreedUsersId) != 0 {
		messages = append(messages, ":warning: 反対票があります！")
	}

	if len(messages) == 0 {
		messages = append(messages, fmt.Sprintf("賛成票 %d 票で可決されました！", len(agreedUsersId)))
	} else {
		messages = append(messages, "自動で決議されませんでした。")
	}

	messages = append(messages, fmt.Sprintf("debug: agreed count: %d, disagreed count: %d", len(agreedUsersId), len(disagreedUsersId)))

	_, err = b.PostMessageEmbed(context.Background(), msg.ChannelId, fmt.Sprintf("@Takeno_hito\n\n%s\n%s", strings.Join(messages, "\n"), "https://q.trap.jp/messages/"+msg.Id))

	if err != nil {
		return err
	}

	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	_, err = b.API().
		MessageApi.EditMessage(context.Background(), msg.Id).
		PostMessageRequest(traq.PostMessageRequest{
			Content: fmt.Sprintf("%s\n\n【%s 更新】\n%s", msg.Content, time.Now().In(location).Format("01-02 15:04"), strings.Join(messages, "\n")),
		}).
		Execute()

	if err != nil {
		return err
	}

	return err
}
