package main

import (
	"fmt"
	"keywatcher"
)

func main () {
	app, err := keywatcher.New()
	if err != nil {
		fmt.Print(err)
		return
	}
	app.Run()
}