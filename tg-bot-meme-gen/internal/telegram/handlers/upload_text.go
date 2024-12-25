package handlers

import (
	tele "gopkg.in/telebot.v3"
)

type Drawer interface {
	DrawText(pathToImage, topText, bottomText string) (string, error)
}

func NewUploadText(drawer Drawer) tele.HandlerFunc {
	return func(c tele.Context) error {
		text := c.Text()
		userID := c.Sender().ID
		if storage[userID].topText == "" {
			meme := storage[userID]
			meme.topText = text
			storage[userID] = meme

			return c.Send("Send bottom text")
		}
		meme := storage[userID]
		meme.bottomText = text

		delete(storage, userID)

		resultImagePath, err := drawer.DrawText(meme.imagePath, meme.topText, meme.bottomText)
		if err != nil {
			return err
		}

		return c.Send(&tele.Photo{File: tele.FromDisk(resultImagePath)})
	}
}
