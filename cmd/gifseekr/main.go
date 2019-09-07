package main

import (
	"github.com/davecgh/go-spew/spew"

	"github.com/dustinrohde/gifseekr/pkg/gifloader"
	"github.com/dustinrohde/gifseekr/pkg/ui"
)

func main() {
	settings, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	client := gifloader.NewGiphyClient(settings.GiphyAPIKey).PageSize(9)
	handle := client.Search("food")
	page, err := handle.Next()

	previews, err := gifloader.LoadPreviews(page)
	if err != nil {
		panic(err)
	}
	spew.Config.MaxDepth = 3
	spew.Dump(previews)

	app := ui.NewApp()
	for i, preview := range previews {
		app.UpdateImage(i, preview.GIF.Image[0])
	}
	app.Run()
}
