package config

import (
	"log"
	"os"

	"github.com/mrvin/tasks-go/tg-bot-meme-gen/internal/telegram"
)

type Config struct {
	Telegram telegram.Conf
}

func (c *Config) LoadFromEnv() {
	if tgToken := os.Getenv("TG_TOKEN"); tgToken != "" {
		c.Telegram.Token = tgToken
	} else {
		log.Println("Empty telegram token")
	}
}
