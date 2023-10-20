package main

import (
	"log"

	"github.com/lemjoe/md-blog/internal"
)

func main() {
	app := internal.NewApp()
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
