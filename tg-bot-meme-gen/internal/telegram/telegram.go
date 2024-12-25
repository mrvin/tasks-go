package telegram

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mrvin/tasks-go/tg-bot-meme-gen/internal/app"
	"github.com/mrvin/tasks-go/tg-bot-meme-gen/internal/telegram/handlers"
	tele "gopkg.in/telebot.v3"
)

type Conf struct {
	Token string
}

type TgBot struct {
	bot *tele.Bot
	app *app.Application
}

func NewBot(conf *Conf, a *app.Application) (*TgBot, error) {
	pref := tele.Settings{
		Token:  conf.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create new tg bot: %w", err)
	}

	b.Handle("/hello", handlers.Health)
	b.Handle("/start", handlers.Start)
	b.Handle(tele.OnPhoto, handlers.NewUploadImage(b))
	b.Handle(tele.OnText, handlers.NewUploadText(a))

	return &TgBot{b, a}, nil
}

func (t *TgBot) Run(ctx context.Context) {
	go func() {
		log.Println("tg-bot-meme-gen: telegram bot starting...")
		t.bot.Start()
	}()

	<-ctx.Done()
	log.Println("tg-bot-meme-gen: telegram bot stopping...")
	t.bot.Stop()
	log.Println("tg-bot-meme-gen: telegram bot stopped")
}
