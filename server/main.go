package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/traP-jp/requestan/server/bot"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Info("Welcome to requestan!")

	bot.New(os.Getenv("TRAQ_BOT_ID"), os.Getenv("TRAQ_ACCESS_TOKEN"))

	log.Info("application has started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-quit

	log.Warn("Shutting down...")
}
