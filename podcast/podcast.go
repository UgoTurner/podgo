package podcast

import (
	"github.com/jroimartin/gocui"
	"github.com/ugo/podcastor/utils"
)

const (
  VIEW_NAME = "side"
  VIEW_TITLE = "Podcasts"
)

func getPodcastNames() []string {
  return []string {"podcast1", "podcast2", "podcast3"}
}

func Init(g *gocui.Gui) (*gocui.View, error) {
	_, maxY := g.Size()
  v, err := g.SetView(VIEW_NAME, 0, 0, 30, maxY-1)
	if err != nil && err != gocui.ErrUnknownView {
	  return nil, err
	}
  v.Highlight = true
  v.SelBgColor = gocui.ColorGreen
  v.SelFgColor = gocui.ColorBlack
  v.Frame = true
  v.Title = VIEW_TITLE
  utils.UpdateList(v, getPodcastNames())

  return v, nil
}

