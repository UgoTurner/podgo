package episode

import (
	"github.com/jroimartin/gocui"
	"github.com/ugo/podcastor/utils"
)

const (
  VIEW_NAME = "main"
  VIEW_TITLE = "Episodes"
)

func getEpisodeNames() []string {
  return []string {"episode1", "episode2"}
}

func Init(g *gocui.Gui) (*gocui.View, error) {
	maxX, maxY := g.Size()
	v, err := g.SetView(VIEW_NAME, 31, 0, maxX-1, maxY-1)
	if err != nil && err != gocui.ErrUnknownView {
	  return nil, err
	}
  v.Highlight = true
  v.SelBgColor = gocui.ColorGreen
  v.SelFgColor = gocui.ColorBlack
  v.Frame = true
  v.Title = VIEW_TITLE
  utils.UpdateList(v, getEpisodeNames())

  return v, nil
}
