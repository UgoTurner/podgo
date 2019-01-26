package cmd

import (
	"log"

	"github.com/jroimartin/gocui"
	"github.com/ugo/podgo/event"
	"github.com/ugo/podgo/handler"
	"github.com/ugo/podgo/keybind"
	"github.com/ugo/podgo/service"
	"github.com/ugo/podgo/ui"
)

func Build() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
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
