package main

import (
	"image/gif"
	"net/http"

	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"github.com/davecgh/go-spew/spew"

	"github.com/dustinrohde/gifseekr/pkg/gifloader"
)

const gifURL = "http://giphygifs.s3.amazonaws.com/media/M8xmO5ZcLPtAY/giphy.gifloader"

func main() {
	settings, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	client := gifloader.NewGiphyClient(settings.GiphyAPIKey).PageSize(2)
	handle := client.Search("food")
	page, err := handle.Next()
	if err != nil {
		println("ERROR ===============V")
		spew.Dump(err)
	}
	if page != nil {
		println("PAGE JSON ===============V")
		spew.Dump(page)
	}

	response, err := http.Get(gifURL)
	if err != nil {
		panic(err)
	}

	gifimg, err := gif.Decode(response.Body)
	image := &canvas.Image{FillMode: canvas.ImageFillOriginal}
	image.Image = gifimg

	a := app.New()
	win := a.NewWindow("Hello")
	win.SetContent(widget.NewVBox(
		widget.NewButton("Quit", func() {
			a.Quit()
		}),
		image,
	))
	canvas.Refresh(image)
	win.ShowAndRun()
}
