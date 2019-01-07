package main

import (
	"github.com/jroimartin/gocui"
	"github.com/ugo/podcastor/podcast"
  "fmt"
)

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}

	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
  if v == nil {
    if _, err := g.SetCurrentView(podcast.VIEW_NAME); err != nil {
      return err
    }
  }
  if v != nil {
		fmt.Fprintln(v, "toto")
  }
  // for _, view := range g.Views() {
		// fmt.Fprintln(view, view.Name())
    // // fmt.Sprintf("%v", v.Name())
  // }

  return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
