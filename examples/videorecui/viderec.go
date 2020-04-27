package main

import (
	"context"
	"log"
	"os"

	"github.com/bksworm/gst"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
)

const appId = "com.github.bksworm.gst.videorec"

type CamUiApp struct {
	application *gtk.Application
	player      *Player
	ctx         context.Context
}

func NewCamUiApp() (app *CamUiApp, err error) {
	app = &CamUiApp{}
	app.application, err = gtk.ApplicationNew(appId, glib.APPLICATION_FLAGS_NONE)
	app.ctx = context.Background()
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
	signals["on_pause_show_btn_clicked"] = func() {
		log.Println("Pause")
		p.Pause()
	}
	signals["on_record_btn_clicked"] = func() {
		log.Println("Start recording")
		err := p.rec.PullCtrl(-1)
		if err != nil {
			log.Println(err.Error())
		}
	}
	signals["on_stop_recording_btn_clicked"] = func() {
		log.Println("Stop recording")
		err := p.rec.PullCtrl(0)
		if err != nil {
			log.Println(err.Error())
		}
	}

	builder.ConnectSignals(signals)

	go p.PictureTaker(ca.ctx, "./out")
	p.MovieMaker(ca.ctx, "./out")
	ca.player = p

	ca.application.AddWindow(win)
	// Show the Window and all of its components.
	win.ShowAll()
}

func (ca *CamUiApp) onQuitButton(obj *gtk.Button) {
	log.Println("onQuitButton")
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
		log.Printf("%v", e) //with frame
		os.Exit(-1)
	}
}

// onMainWindowDestory is the callback that is linked to the
// on_main_window_destroy handler. It is not required to map this,
// and is here to simply demo how to hook-up custom callbacks.
func onMainWindowDestroy() {
	log.Println("onMainWindowDestroy")
}
