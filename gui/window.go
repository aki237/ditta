// package gui provides gui utilities to the app
package gui

import (
	"github.com/aki237/ditta/manager"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// DWindow struct is a derived from gtk.Window  and contains some other unexported widgets.
type DWindow struct {
	*gtk.ApplicationWindow
	dArea   *gtk.DrawingArea
	manager *manager.Manager
}

// NewDWidnow creates a new Ditta Window Struct and returns it with errors if any
func NewDWidnow(app *DApplication) (*DWindow, error) {

	// Initialize new gtk.ApplicationWindow
	window, err := gtk.ApplicationWindowNew(app.Application)
	if err != nil {
		return nil, err
	}
	// Wrap the gtk.Window in a DWindow struct
	dwindow := &DWindow{ApplicationWindow: window}

	// Initialize the drawing area widget
	dwindow.dArea, err = gtk.DrawingAreaNew()
	if err != nil {
		return nil, err
	}

	// Setup Manager backend
	dwindow.manager = manager.NewManager("ditta.mainWindow")

	// Connect necessary signals to the drawing area
	dwindow.dArea.Connect("draw", drawLoop, dwindow.manager)
	dwindow.Connect("key-press-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		dwindow.manager.ReadKey(win, ev)
	})
	dwindow.Connect("key-release-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		dwindow.manager.CheckModifier(win, ev)
	})
	// Add gtk.DrawingArea to the window
	dwindow.Add(dwindow.dArea)
	return dwindow, nil
}
