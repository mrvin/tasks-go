package handlers

import (
	tele "gopkg.in/telebot.v3"
)

func Health(c tele.Context) error {
	return c.Send("Hello!") //nolint:wrapcheck
}
