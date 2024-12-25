package handlers

import (
	"io"
	"os"

	tele "gopkg.in/telebot.v3"
)

type memeData struct {
	imagePath  string
	topText    string
	bottomText string
}

var storage = map[int64]memeData{}

func NewUploadImage(bot *tele.Bot) tele.HandlerFunc {
	return func(c tele.Context) error {
		photo := c.Message().Photo
		rcImage, err := bot.File(&tele.File{FileID: photo.FileID})
		if err != nil {
			return err
		}
		uploadImage, err := os.CreateTemp("", "telegram-upload")
		if err != nil {
			return err
		}

		if _, err := io.Copy(uploadImage, rcImage); err != nil {
			return err
		}
		userID := c.Sender().ID
		storage[userID] = memeData{imagePath: uploadImage.Name()}

		return c.Send("Send top text") 
	}
}
