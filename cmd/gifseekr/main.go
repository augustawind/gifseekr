package main

import (
	"github.com/davecgh/go-spew/spew"

	"github.com/dustinrohde/gifseekr/pkg/gif"
)

const apiKey = "jEYeMjhDmvnJpAsKEIauvWzWyPae6wwm"

func main() {
	client := gif.NewGiphyClient(apiKey).PageSize(2)
	handle := client.Search("food")
	resp, err := handle.Next()
	if err != nil {
		spew.Dump("ERROR: ", err)
		spew.Dump("BODY: ", resp)
	} else {
		spew.Dump(resp)
	}
}
