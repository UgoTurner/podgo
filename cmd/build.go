package cmd

import (
	"os"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
	"github.com/ugo/podgo/conf"
	"github.com/ugo/podgo/event"
	"github.com/ugo/podgo/handler"
	"github.com/ugo/podgo/keybind"
	"github.com/ugo/podgo/service"
	"github.com/ugo/podgo/ui"
)

func init() {
	file, err := os.OpenFile(conf.logFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		log.SetFormatter(&log.JSONFormatter{})
		/* the log level is set to "Info" which allows
		to log Info(), Warn(), Error() and Fatal()*/
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
}

func Build() {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln("Error when creating GUI")
	}
	defer g.Close()
	render := &ui.Render{TUI: g}
	render.InitLayout()
	render.LoadPanels("./conf/panels.json")

	eventDispatcher := &event.Dispatcher{}
	keybinder := &keybind.Keybinder{
		EventDispatcher: eventDispatcher,
		TUI:             g,
	}
	keybinder.LoadKeybinds("./conf/keybinds.json")

	appSubscriber := &event.Subscriber{
		Handler: &handler.AppHandler{
			FeedParser: &service.FeedParser{},
			Render:     render,
			Player:     &service.Player{},
		},
	}

	eventDispatcher.AddSubscriber(appSubscriber)
	eventDispatcher.Dispatch("Launch")

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
