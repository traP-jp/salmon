package bot

import (
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

type Bot struct {
	botID string
	bot   *traqwsbot.Bot
}

func New(botID string, traQAccessToken string) Bot {
	_b, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: traQAccessToken,
	})
	if err != nil {
		panic(err)
	}

	b := Bot{botID, _b}

	b.bot.OnError(func(message string) {
		log.Error("Received ERROR message: " + message)
	})
	b.bot.OnMessageCreated(b.onMessageCreated)
	go func() {
		if err := b.bot.Start(); err != nil {
			panic(err)
		}
	}()

	return b
}

func (b *Bot) api() *traq.APIClient {
	return b.bot.API()
}
