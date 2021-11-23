package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"time"

	ws "github.com/ChristianHering/WaveShare"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/jrydval/svatky"
)

var (
	Debug      bool = false
	faceBig    font.Face
	faceSmall  font.Face
	faceMedium font.Face
)

const (
	WIDTH           = 800
	HEIGHT          = 480
	BOTTOMBARHEIGHT = 60
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

		tempText := getTemp()
		log.Print(tempText)
		drawTemp(dc, tempText)

		drawSvatek(dc)

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
	faceMedium = truetype.NewFace(fontFaceBold, &truetype.Options{
		Size:    80,
		Hinting: font.HintingFull,
	})
	faceSmall = truetype.NewFace(fontFaceBold, &truetype.Options{
		Size:    38,
		Hinting: font.HintingFull,
	})
}

func displayImage(i *image.Image) {
	ws.DisplayImage(*i)
	ws.Sleep()
}

var dayString string
var nameString string
var datumString string

func getClockText() string {
	monthNames := []string{"", "Leden", "Únor", "Březen", "Duben", "Květen", "Červen", "Červenec", "Září", "Říjen", "Listopad", "Prosinec"}
	t := time.Now()
	h := t.Hour()
	m := t.Minute()
	text := fmt.Sprintf("%02d:%02d", h, m)
	dayStringNew := fmt.Sprintf("%02d", t.Day())
	if dayStringNew != dayString {
		dayString = dayStringNew

		monthName := monthNames[(t.Month())]
		datumString = fmt.Sprintf("%s %s", dayString, monthName)
		svatekIndex := fmt.Sprintf("%s.%02d.", dayString, t.Month())
		nameString = svatky.GetSvatekByDate(svatekIndex)
		log.Printf("Svatek and date update: %s, %s", nameString, datumString)
	}

	return text
}

func drawClock(dc *gg.Context, text string) {
	dc.SetFontFace(faceBig)
	dc.SetColor(color.Black)
	dc.DrawStringAnchored(text, WIDTH/2, HEIGHT/2, 0.5, 0.5)
}

func drawSvatek(dc *gg.Context) {
	dc.SetFontFace(faceSmall)
	dc.SetColor(color.Black)
	dc.DrawStringAnchored(nameString, 10, 90, 0, 1)
	dc.SetFontFace(faceMedium)
	dc.DrawStringAnchored(datumString, 10, 10, 0, 1)

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
	//dc.DrawRectangle(0, HEIGHT-BOTTOMBARHEIGHT, WIDTH, BOTTOMBARHEIGHT)
	//dc.Fill()
	//dc.SetColor(color.White)
	dc.SetFontFace(faceSmall)
	dc.DrawStringAnchored(text, WIDTH/2, (HEIGHT - BOTTOMBARHEIGHT/2), 0.5, 0.4)
}

func getTemp() string {
	resp, err := http.Get("http://10.0.5.10:3480/data_request?id=variableget&DeviceNum=28&serviceId=urn:upnp-org:serviceId:TemperatureSensor1&Variable=CurrentTemperature")
	if err != nil {
		log.Printf("HTTP error: %v", err)
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("HTTP stream error: %v", err)
		return ""
	}

	temp, err := strconv.ParseFloat(string(body), 32)

	if err != nil {
		log.Printf("Temperature conversion error: %v", err)
		return ""
	}

	return fmt.Sprintf("%.1f°C", temp)
}

func drawTemp(dc *gg.Context, text string) {

	dc.SetColor(color.Black)
	dc.SetFontFace(faceMedium)
	dc.DrawStringAnchored(text, WIDTH-10, 10, 1, 1)
}
