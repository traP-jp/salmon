package bot

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/traPtitech/go-traq"
)

func (b *Bot) PostMessage(ctx context.Context, channelID string, content string) error {
	_, _, err := b.api().
		MessageApi.
		PostMessage(ctx, channelID).
		PostMessageRequest(traq.PostMessageRequest{
			Content: content,
		}).
		Execute()
	return err
}

func (b *Bot) GetMessageFromMessageId(ctx context.Context, id uuid.UUID) (*traq.Message, error) {
	message, _, err := b.api().
		MessageApi.
		GetMessage(ctx, id.String()).
		Execute()
	if err != nil {
		return nil, err
	}

	return message, nil
}
