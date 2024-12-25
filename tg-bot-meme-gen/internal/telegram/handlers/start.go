package handlers

import (
	tele "gopkg.in/telebot.v3"
)

func Start(c tele.Context) error {
	return c.Send("Send image") //nolint:wrapcheck
}
