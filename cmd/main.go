package main

import (
	"github.com/lemjoe/Grapho/internal"
)

func main() {
	app := internal.NewApp()
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
