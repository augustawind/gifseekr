package ui

import (
	"image"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
)

const (
	title    = "gifseekr"
	gridRows = 3
	gridCols = 3
	gridSize = gridRows * gridCols
)

type App struct {
	app    fyne.App
	window fyne.Window
	images []*canvas.Image
}

func NewApp() *App {
	a := new(App)
	a.app = app.New()
	a.window = a.app.NewWindow(title)
	a.images = make([]*canvas.Image, gridSize)

	grid := fyne.NewContainerWithLayout(layout.NewGridLayout(gridCols))
	for i := 0; i < gridSize; i++ {
		a.images[i] = &canvas.Image{FillMode: canvas.ImageFillOriginal}
		grid.AddObject(a.images[i])
	}

	a.window.SetContent(grid)
	return a
}

func (a *App) UpdateImage(i int, img image.Image) {
	a.images[i].Image = img
	canvas.Refresh(a.images[i])
}

func (a *App) Run() {
	a.window.ShowAndRun()
}
