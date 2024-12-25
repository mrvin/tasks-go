package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrvin/tasks-go/tg-bot-meme-gen/internal/app"
	"github.com/mrvin/tasks-go/tg-bot-meme-gen/internal/config"
	"github.com/mrvin/tasks-go/tg-bot-meme-gen/internal/telegram"
)

func main() {
	var conf config.Config

	conf.LoadFromEnv()

	a, err := app.New()
	if err != nil {
		log.Printf("tg-bot-meme-gen: error: %v", err)
		return
	}
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(
		ctx,
		os.Interrupt,    // SIGINT, (Control-C)
		syscall.SIGTERM, // systemd
		syscall.SIGQUIT,
	)
	defer cancel()

	tgBot, err := telegram.NewBot(&conf.Telegram, a)
	if err != nil {
		log.Printf("tg-bot-meme-gen: error: %v", err)
		return
	}

	tgBot.Run(ctx)
}
