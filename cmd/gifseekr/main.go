package main

import (
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

	previews, err := gifloader.LoadPreviews(page)
	if err != nil {
		panic(err)
	}
	spew.Config.MaxDepth = 3
	spew.Dump(previews)

	image := &canvas.Image{FillMode: canvas.ImageFillOriginal}
	image.Image = previews[0].GIF.Image[0]

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
