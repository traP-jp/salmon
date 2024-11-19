package bot

import (
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

type Bot struct {
	botID  string
	userID string
	bot    *traqwsbot.Bot
	env    Environment
}

type Environment string

const (
	EnvProduction Environment = "production"
	EnvLocal      Environment = "local"
)

func New(botID string, traQAccessToken string, isLocal bool) Bot {
	_b, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: traQAccessToken,
	})
	if err != nil {
		panic(err)
	}

	var env Environment

	if isLocal {
		env = EnvLocal
	} else {
		env = EnvProduction
	}

	b := Bot{
		botID: botID,
		bot:   _b,
		env:   env,
	}

	b.bot.OnError(func(message string) {
		log.Error("Received ERROR message: " + message)
	})
	b.bot.OnMessageCreated(b.joinOrLeaveHandler)
	go func() {
		if err := b.bot.Start(); err != nil {
			panic(err)
		}
	}()

	m, _, err := b.API().MeApi.GetMe(nil).Execute()
	if err != nil {
		panic(err)
	}

	b.userID = m.Id

	return b
}

func (b *Bot) OnMessageCreated(h func(p *payload.MessageCreated)) {
	b.bot.OnMessageCreated(h)
}

func (b *Bot) API() *traq.APIClient {
	return b.bot.API()
}

func (b *Bot) Env() Environment {
	return b.env
}
