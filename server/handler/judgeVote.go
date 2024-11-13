package handler

import (
	"context"
	"database/sql"
	"fmt"
	"git.trap.jp/Takeno-hito/salmon/server/bot"
	"git.trap.jp/Takeno-hito/salmon/server/database"
	"github.com/gofrs/uuid"
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
		messages = append(messages, ":warning: 反対票があります！自動決議はキャンセルされました。")
	}

	if len(messages) == 0 {
		messages = append(messages, fmt.Sprintf("賛成票 %d 票で可決されました！", len(agreedUsersId)))
	}

	messages = append(messages, fmt.Sprintf("debug: agreed count: %s, disagreed count: %s", len(agreedUsersId), len(disagreedUsersId)))

	_, err = b.PostMessageEmbed(context.Background(), msg.ChannelId, fmt.Sprintf("@Takeno_hito \n\n %s", strings.Join(messages, "\n")))
	return err
}

func (h Handler) JudgeVote() error {
	tasks, err := h.db.GetScheduledTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.Command != "judge-vote" || task.ScheduledAt.After(time.Now()) || task.ExecutedAt.Valid {
			continue
		}

		err = h.db.UpdateScheduledTask(database.ScheduledTask{
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
			return err
		}

		if err = judge(h.bot, task.Arg); err != nil {
			return err
		}
	}

	return nil
}
