package main

import (
	"github.com/davecgh/go-spew/spew"

	"github.com/dustinrohde/gifseekr/pkg/gif"
)

func main() {
	settings, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	client := gif.NewGiphyClient(settings.GiphyAPIKey).PageSize(2)
	handle := client.Search("food")
	resp, err := handle.Next()
	if err != nil {
		spew.Dump("ERROR: ", err)
		spew.Dump("BODY: ", resp)
	} else {
		spew.Dump(resp)
	}
}
