package gui

import (
	"fmt"
	"log"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// DApplication struct is a derived from gtk.Application
type DApplication struct {
	*gtk.Application
}

func NewDApplication() (*DApplication, error) {
	// Create a new gtk.Application
	app, err := gtk.ApplicationNew("org.aki237.ditta", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return nil, err
	}
	app.Connect("activate", activate, nil)
	dapp := &DApplication{app}
	return dapp, nil
}

func activate(app *gtk.Application) {
	window, err := NewDWidnow(&DApplication{app})
	if len(os.Args) >= 2 {
		window.manager.SetFileName(os.Args[1])
		fmt.Println("File Opened. : ", os.Args[1])
	}
	if err != nil {
		log.Fatal(err)
	}
	window.ShowAll()
}
