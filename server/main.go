package main

import (
	"git.trap.jp/Takeno-hito/salmon/server/bot"
	"git.trap.jp/Takeno-hito/salmon/server/database"
	"git.trap.jp/Takeno-hito/salmon/server/handler"
	"github.com/go-co-op/gocron/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Info("Welcome to salmon!")

	traQBotId := os.Getenv("TRAQ_BOT_ID")
	traQAccessToken := os.Getenv("TRAQ_ACCESS_TOKEN")

	dbUser := os.Getenv("NS_MARIADB_USER")
	dbPass := os.Getenv("NS_MARIADB_PASSWORD")
	dbHost := os.Getenv("NS_MARIADB_HOSTNAME")
	dbPort := os.Getenv("NS_MARIADB_PORT")
	dbName := os.Getenv("NS_MARIADB_DATABASE")
	isLocal := os.Getenv("IS_LOCAL")

	b := bot.New(traQBotId, traQAccessToken, isLocal == "true")

	log.SetLevel(log.DebugLevel)

	db, err := database.NewClientAndMigrate(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s, err := gocron.NewScheduler()

	if err != nil {
		log.Fatal(err)
	}

	h := handler.New(&b, db)
	b.OnMessageCreated(h.HandleBotMessage)

	_, err = s.NewJob(
		gocron.DurationJob(1*time.Minute),
		gocron.NewTask(h.JudgeVote),
	)

	if err != nil {
		log.Fatal(err)
	}

	s.Start()
	defer func() {
		err := s.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Info("application has started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-quit

	log.Warn("Shutting down...")
}
