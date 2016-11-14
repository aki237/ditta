package main

import (
	"log"

	"github.com/aki237/ditta/gui"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)
	app, err := gui.NewDApplication()
	if err != nil {
		log.Fatal(err)
	}
	app.Run(nil)
}
