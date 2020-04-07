package main

import (
	"errors"
	"log"
	"os"

	"github.com/bksworm/gst"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const appId = "com.github.bksworm.gst.videorec"

type CamUiApp struct {
	application *gtk.Application
	player      *Player
}

func NewCamUiApp() (app *CamUiApp, err error) {
	app = &CamUiApp{}
	app.application, err = gtk.ApplicationNew(appId, glib.APPLICATION_FLAGS_NONE)
	return
}

func (ca *CamUiApp) Run() int {
	return ca.application.Run(os.Args)
}
func (ca *CamUiApp) ConnectAppSignals() {

	// Connect function to application startup event, this is not required.
	ca.application.Connect("startup", func() {
		log.Println("application startup")
	})
	// Connect function to application startup event, this is not required.
	ca.application.Connect("startup", func() {
		log.Println("application startup")
	})
	// Connect function to application shutdown event, this is not required.
	ca.application.Connect("shutdown", func() {
		log.Println("application shutdown")
	})
	// Connect function to application activate event
	ca.application.Connect("activate", ca.onActivate)
}

func (ca *CamUiApp) onActivate() {
	log.Println("application activate")

	// Get the GtkBuilder UI definition in the glade file.
	builder, err := gtk.BuilderNewFromFile("ui/videorec.ui")
	errorCheck(err)

	// Map the handlers to callback functions, and connect the signals
	// to the Builder.
	signals := map[string]interface{}{
		"on_main_window_destroy": onMainWindowDestroy,
		"on_quit_btn_clicked":    ca.onQuitButton,
	}

	// Get the object with the id of "main_window".
	obj, err := builder.GetObject("main_window")
	errorCheck(err)

	// Verify that the object is a pointer to a gtk.ApplicationWindow.
	win, err := isWindow(obj)
	errorCheck(err)

	p := NewPlayer()
	err = p.Assemble()
	errorCheck(err)
	//get container for the gtksink widget
	obj, err = builder.GetObject("top_box")
	errorCheck(err)
	// Verify that the object is a pointer to a gtk.Box
	vbox, err := isBox(obj)
	errorCheck(err)
	vbox.PackStart(p.widget, true, true, 0) //pack gtksink to the ui

	//add player buttons handlers
	signals["on_show_btn_clicked"] = func() {
		log.Println("Play")
		p.Play()
	}
	signals["on_stop_btn_clicked"] = func() {
		log.Println("Pause")
		p.Pause()
	}
	signals["on_photo_btn_clicked"] = func() {
		log.Println("Take picture")
		p.TakePicture()
	}
	signals["on_photo_btn_clicked"] = func() {
		log.Println("Take picture")
		p.TakePicture()
	}
	signals["on_photo_btn_clicked"] = func() {
		log.Println("Take picture")
		p.TakePicture()
	}

	builder.ConnectSignals(signals)

	go p.PictureTaker("./out")
	ca.player = p

	ca.application.AddWindow(win)
	// Show the Window and all of its components.
	win.ShowAll()
}

func (ca *CamUiApp) onQuitButton(obj *gtk.Button) {
	log.Println("onQuitButton")
	//log.Printf("obj %#v", obj)
	//log.Printf("ca %#v", ca)
	ca.application.Quit()
}
func main() {
	// Create a new application.
	application, err := NewCamUiApp()
	errorCheck(err)
	application.ConnectAppSignals()

	// Launch the application
	ml := gst.MainLoopNew()
	ret := application.Run()
	ml.Quit()
	application.player.Close()
	os.Exit(ret)
}

func isWindow(obj glib.IObject) (*gtk.Window, error) {
	// Make type assertion (as per gtk.go).
	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.Window")
}

func isBox(obj glib.IObject) (*gtk.Box, error) {
	// Make type assertion (as per gtk.go).
	if win, ok := obj.(*gtk.Box); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.Box")
}
func errorCheck(e error) {
	if e != nil {
		// panic for any errors.
		log.Panic(e)
	}
}

// onMainWindowDestory is the callback that is linked to the
// on_main_window_destroy handler. It is not required to map this,
// and is here to simply demo how to hook-up custom callbacks.
func onMainWindowDestroy() {
	log.Println("onMainWindowDestroy")
}
