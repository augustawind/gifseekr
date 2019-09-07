package ui

import (
	"image"
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/dustinrohde/gifseekr/pkg/gifloader"
)

var stderr = log.New(os.Stderr, "", 0)

const (
	title    = "gifseekr"
	gridRows = 3
	gridCols = 3
	gridSize = gridRows * gridCols
)

type AppSettings struct {
	GiphyAPIKey string `mapstructure:"giphy-api-key"`
}

type App struct {
	app    fyne.App
	window fyne.Window
	images []*canvas.Image

	client *gifloader.GiphyClient
	handle *gifloader.SearchHandle
}

func NewApp(settings *AppSettings) *App {
	a := new(App)
	a.app = app.New()
	a.window = a.app.NewWindow(title)
	a.images = make([]*canvas.Image, gridSize)

	searchBar := a.makeSearchBar()
	grid := a.makeImageGrid()

	a.window.Resize(grid.Size().Add(fyne.NewSize(0, 100)))
	a.window.SetContent(widget.NewVBox(searchBar, grid))

	a.client = gifloader.NewGiphyClient(settings.GiphyAPIKey).PageSize(gridSize)

	return a
}

func (a *App) UpdateImage(i int, img image.Image) {
	a.images[i].Image = img
	canvas.Refresh(a.images[i])
}

func (a *App) Run() {
	a.window.ShowAndRun()
}

func (a *App) makeSearchBar() fyne.CanvasObject {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("What are you looking for?")
	submit := widget.NewButton("Submit", func() {
		query := entry.Text
		a.handle = a.client.Search(query)
		result, err := a.handle.Next()
		if err != nil {
			stderr.Printf("error occured fetching search results: %s\n", err.Error())
		}
		previews, err := gifloader.LoadPreviews(result)
		if err != nil {
			stderr.Printf("error occured loading GIF previews: %s\n", err.Error())
		}
		for i, preview := range previews {
			a.UpdateImage(i, preview.GIF.Image[0])
		}
	})
	return widget.NewHBox(entry, submit)
}

func (a *App) makeImageGrid() fyne.CanvasObject {
	grid := fyne.NewContainerWithLayout(layout.NewGridLayout(gridCols))
	for i := 0; i < gridSize; i++ {
		a.images[i] = &canvas.Image{FillMode: canvas.ImageFillOriginal}
		a.images[i].SetMinSize(fyne.NewSize(200, 200))
		grid.AddObject(a.images[i])
	}
	return grid
}
