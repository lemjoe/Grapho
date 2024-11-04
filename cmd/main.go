package main

import (
	"fmt"

	"github.com/lemjoe/Grapho/internal"
)

var Version = "development"
var BuildTime string

func main() {
	fmt.Println("Grapho version:\t", Version)
	fmt.Println("Build time:\t", BuildTime)
	app := internal.NewApp()
	err := app.Run(Version)
	if err != nil {
		panic(err)
	}
}
