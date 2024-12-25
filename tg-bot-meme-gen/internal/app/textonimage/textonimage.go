package textonimage

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

func DrawText(font *truetype.Font, pathToImage, topText, bottomText string) (string, error) {
	imageFile, err := os.Open(pathToImage)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}
	defer imageFile.Close()

	originalImage, _, err := image.Decode(imageFile)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	rgba := image.NewRGBA(originalImage.Bounds())
	draw.Draw(rgba, rgba.Bounds(), originalImage, image.Point{}, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(22)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.NewUniform(color.White))

	pt := freetype.Pt(10, 40)
	if _, err := c.DrawString(topText, pt); err != nil {
		return "", fmt.Errorf("failed to draw top text: %w", err)
	}
	pt = freetype.Pt(10, rgba.Bounds().Dy()-10)
	if _, err := c.DrawString(bottomText, pt); err != nil {
		return "", fmt.Errorf("failed to draw bottom text: %w", err)
	}

	resultImage, err := os.CreateTemp("", "meme")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer resultImage.Close()

	if err := jpeg.Encode(resultImage, rgba, nil); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	return resultImage.Name(), nil
}
