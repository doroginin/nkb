package main

import (
	"fmt"
	"nkb"
)

func main () {
	app, err := nuclear_kb.New()
	if err != nil {
		fmt.Print(err)
		return
	}
	app.Run()
}