package main

import (
	"fmt"
	"github.com/rilinor/nkb"
)

func main () {
	app, err := nkb.New()
	if err != nil {
		fmt.Print(err)
		return
	}
	app.Run()
}