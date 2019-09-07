package main

import (
	"github.com/dustinrohde/gifseekr/pkg/ui"
)

func main() {
	settings, err := ReadConfig()
	if err != nil {
		panic(err)
	}
	app := ui.NewApp(settings)
	app.Run()
}
