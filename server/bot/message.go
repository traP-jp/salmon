package bot

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/traPtitech/go-traq"
)

func (b *Bot) PostMessage(ctx context.Context, channelID string, content string) error {
	_, _, err := b.API().
		MessageApi.
		PostMessage(ctx, channelID).
		PostMessageRequest(traq.PostMessageRequest{
			Content: content,
		}).
		Execute()
	return err
}

func (b *Bot) AttachVoteStamps(ctx context.Context, messageID uuid.UUID) error {
	_, err := b.API().
		MessageApi.
		AddMessageStamp(ctx, messageID.String(), AgreeStampId).PostMessageStampRequest(traq.PostMessageStampRequest{
		Count: 0,
	}).Execute()
	if err != nil {
		return err
	}

	_, err = b.API().
		MessageApi.
		AddMessageStamp(ctx, messageID.String(), DisagreeStampId).PostMessageStampRequest(traq.PostMessageStampRequest{
		Count: 0,
	}).Execute()
	return err
}

// PostMessageEmbed return messageId or Error
func (b *Bot) PostMessageEmbed(ctx context.Context, channelID string, content string) (string, error) {
	msg, _, err := b.API().
		MessageApi.
		PostMessage(ctx, channelID).
		PostMessageRequest(traq.PostMessageRequest{
			Content: content,
			Embed:   traq.PtrBool(true),
		}).
		Execute()
	if err != nil {
		return "", err
	}

	return msg.Id, nil
}

func (b *Bot) GetMessageFromMessageId(ctx context.Context, id uuid.UUID) (*traq.Message, error) {
	message, _, err := b.API().
		MessageApi.
		GetMessage(ctx, id.String()).
		Execute()
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (b *Bot) SendDirectMessage(ctx context.Context, userID string, content string) error {
	_, _, err := b.API().
		MessageApi.
		PostDirectMessage(ctx, userID).
		PostMessageRequest(traq.PostMessageRequest{
			Content: content,
		}).
		Execute()
	return err
}
