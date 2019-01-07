package main

import (
	"github.com/jroimartin/gocui"
	"github.com/ugo/podcastor/podcast"
	"github.com/ugo/podcastor/episode"
)

func setKeybindings(g *gocui.Gui) {
  g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
  g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, nextView)

  g.SetKeybinding(podcast.VIEW_NAME, gocui.KeyArrowRight, gocui.ModNone, nextView)
  g.SetKeybinding(podcast.VIEW_NAME, gocui.KeyArrowDown, gocui.ModNone, cursorDown)
  g.SetKeybinding(podcast.VIEW_NAME, gocui.KeyArrowUp, gocui.ModNone, cursorUp)

  g.SetKeybinding(episode.VIEW_NAME, gocui.KeyArrowRight, gocui.ModNone, nextView)
  g.SetKeybinding(episode.VIEW_NAME, gocui.KeyArrowDown, gocui.ModNone, cursorDown)
  g.SetKeybinding(episode.VIEW_NAME, gocui.KeyArrowUp, gocui.ModNone, cursorUp)
}
