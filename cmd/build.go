package cmd

import (
	"log"

	"github.com/jroimartin/gocui"
	"github.com/ugo/podcastor/event"
	"github.com/ugo/podcastor/handler"
	"github.com/ugo/podcastor/keybind"
	"github.com/ugo/podcastor/service"
	"github.com/ugo/podcastor/ui"
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
