package app

import (
	"fmt"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/mrvin/tasks-go/tg-bot-meme-gen/internal/app/textonimage"
)

const pathToFont = "../fonts/Roboto-Regular.ttf"

type Application struct {
	font *truetype.Font
}

func New() (*Application, error) {
	fontBytes, err := os.ReadFile(pathToFont)
	if err != nil {
		return nil, fmt.Errorf("failed to read font file: %w", err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}
	return &Application{font}, nil
}

func (a *Application) DrawText(pathToImage, topText, bottomText string) (string, error) {
	return textonimage.DrawText(a.font, pathToImage, topText, bottomText)
}
