package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"time"

	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"

	ws "github.com/ChristianHering/WaveShare"
)

var (
	Debug     bool = false
	faceBig   font.Face
	faceSmall font.Face
)

const (
	WIDTH           = 800
	HEIGHT          = 480
	BOTTOMBARHEIGHT = 40
)

func main() {
	log.Println("Starting")
	initFonts()

	ws.Initialize()

	for {

		dc := gg.NewContext(WIDTH, HEIGHT)
		img := dc.Image()

		dc.SetColor(color.White)
		dc.DrawRectangle(0, 0, WIDTH, HEIGHT)
		dc.Fill()

		clockText := getClockText()
		log.Print(clockText)
		drawClock(dc, clockText)

		btcText := getBtcText()
		log.Print(btcText)
		drawBtc(dc, btcText)

		if Debug {
			dc.SavePNG("debug.png")
			os.Exit(0)
		}

		displayImage(&img)
		time.Sleep(60 * time.Second)

	}

}

func initFonts() {
	fontFace, _ := freetype.ParseFont(goregular.TTF)
	fontFaceBold, _ := freetype.ParseFont(gobold.TTF)
	faceBig = truetype.NewFace(fontFace, &truetype.Options{
		Size:    290,
		Hinting: font.HintingFull,
	})
	faceSmall = truetype.NewFace(fontFaceBold, &truetype.Options{
		Size:    20,
		Hinting: font.HintingFull,
	})
}

func displayImage(i *image.Image) {
	ws.DisplayImage(*i)
	ws.Sleep()
}

func getClockText() string {
	t := time.Now()
	h := t.Hour()
	m := t.Minute()
	text := fmt.Sprintf("%02d:%02d", h, m)
	return text
}

func drawClock(dc *gg.Context, text string) {
	dc.SetFontFace(faceBig)
	dc.SetColor(color.Black)
	dc.DrawStringAnchored(text, WIDTH/2, HEIGHT/2, 0.5, 0.5)
}

func getBtcText() string {
	resp, err := http.Get("https://www.bitstamp.net/api/v2/ticker/btceur/")
	if err != nil {
		log.Printf("HTTP error: %v", err)
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("HTTP stream error: %v", err)
		return ""
	}
	//log.Println(string(body))

	type BTCPrices struct {
		High float64 `json:"high,string"`
		Low  float64 `json:"low,string"`
		Last float64 `json:"last,string"`
		Open float64 `json:"open,string"`
	}

	var prices BTCPrices
	err = json.Unmarshal(body, &prices)
	if err != nil {
		log.Printf("JSON decoding error: %v", err)
		return ""
	}

	//log.Printf("Values: %+v", prices)

	text := fmt.Sprintf("1 BTC = %.0f EUR  @%.0f ↑%.0f ↓%.0f", prices.Last, prices.Open, prices.High, prices.Low)
	return text

}

func drawBtc(dc *gg.Context, text string) {
	dc.SetColor(color.Black)
	dc.DrawRectangle(0, HEIGHT-BOTTOMBARHEIGHT, WIDTH, BOTTOMBARHEIGHT)
	dc.Fill()
	dc.SetColor(color.White)
	dc.SetFontFace(faceSmall)
	dc.DrawStringAnchored(text, WIDTH/2, (HEIGHT - BOTTOMBARHEIGHT/2), 0.5, 0.4)
}
