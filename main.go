package main

import (
	"fmt"

	"image"
	"image/color"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"

	"time"

	ws "github.com/ChristianHering/WaveShare"
)

func main() {
	fmt.Println("Ahoj")
	ws.Initialize()
	pallete := color.Palette{color.White, color.Black}

	for {
		t := time.Now()
		h := t.Hour()
		m := t.Minute()
		i := image.NewPaletted(image.Rect(0, 0, 800, 480), pallete)
		text := fmt.Sprintf("%02d:%02d", h, m)
		fmt.Println(text)
		drawText(i, text)
		ws.DisplayImage(i)
		ws.Sleep()
		time.Sleep(60 * time.Second)

	}

}

func drawText(canvas *image.Paletted, text string) error {
	var (
		fgColor  image.Image
		fontFace *truetype.Font
		err      error
		fontSize = 290.0
	)
	fgColor = image.Black
	fontFace, err = freetype.ParseFont(goregular.TTF)
	fontDrawer := &font.Drawer{
		Dst: canvas,
		Src: fgColor,
		Face: truetype.NewFace(fontFace, &truetype.Options{
			Size:    fontSize,
			Hinting: font.HintingFull,
		}),
	}
	textBounds, _ := fontDrawer.BoundString(text)
	xPosition := (fixed.I(canvas.Rect.Max.X) - fontDrawer.MeasureString(text)) / 2
	textHeight := textBounds.Max.Y - textBounds.Min.Y
	yPosition := fixed.I((canvas.Rect.Max.Y)-textHeight.Ceil())/2 + fixed.I(textHeight.Ceil())
	fontDrawer.Dot = fixed.Point26_6{
		X: xPosition,
		Y: yPosition,
	}
	fontDrawer.DrawString(text)
	return err
}
