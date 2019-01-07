package main

import (
	"log"
	"github.com/jroimartin/gocui"
	"github.com/ugo/podcastor/podcast"
	"github.com/ugo/podcastor/episode"
)

func layout(g *gocui.Gui) error {
  if _, err := podcast.Init(g); err != nil {
    return err
  }
  if _, err := episode.Init(g); err != nil {
    return err
  }

  return nil
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManagerFunc(layout)
  setKeybindings(g)
	// if _, err := g.SetCurrentView(podcast.VIEW_NAME); err != nil {
	// 	log.Panicln(err)
	// }

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
