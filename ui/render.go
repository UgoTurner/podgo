package ui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ugo/podgo/conf"

	"github.com/jroimartin/gocui"
)

func getTUIcolor(color string) gocui.Attribute {
	switch color {
	case conf.ColorGreen:
		return gocui.ColorGreen
	case conf.ColorBlack:
		return gocui.ColorBlack
	case conf.ColorWhite:
		return gocui.ColorWhite
	case conf.ColorDefault:
		return gocui.ColorDefault
	default:
		return gocui.ColorBlack
	}
}

type Render struct {
	TUI    *gocui.Gui
	Panels []*Panel
}

func (r *Render) LoadPanels(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("Error when opening json file")
	}
	json.Unmarshal(data, &r.Panels)
	for i := range r.Panels {
		r.Panels[i].Coordinate.Scale(r.TUI.Size())
	}
}

func (r *Render) InitLayout() {
	r.TUI.SetManagerFunc(r.layout)
}

func (r *Render) layout(g *gocui.Gui) error {
	r.CreateViews()
	// Focus on side view if no current view
	if r.TUI.CurrentView() == nil {
		if err := r.Focus(conf.SideViewName); err != nil {
			log.Panicln(err)
		}
	}

	return nil
}

func (r *Render) CreateViews() []*gocui.View {
	var views []*gocui.View
	for _, pan := range r.Panels {
		if pan.Hidden {
			continue
		}
		lv := r.CreateView(pan)
		if lv == nil {
			continue
		}
		views = append(views, lv)

	}

	return views
}

func (r *Render) CreateView(p *Panel) *gocui.View {

	v, err := r.TUI.SetView(
		p.Name,
		p.Coordinate.TopLeftXabs,
		p.Coordinate.TopLeftYabs,
		p.Coordinate.BottomRightXabs,
		p.Coordinate.BottomRightYabs,
	)

	if err != nil && err != gocui.ErrUnknownView {
		return nil
	}
	v.SelBgColor = getTUIcolor(p.SelectionColor.BgColorCurrent)
	v.SelFgColor = getTUIcolor(p.SelectionColor.FgColorCurrent)
	v.Highlight = p.Highlight
	v.Frame = p.Frame
	v.Title = p.Title
	v.Wrap = true
	v.Overwrite = p.Overwrite

	return v
}

func (r *Render) UpdateListView(viewName string, data []string) {
	go r.TUI.Update(func(g *gocui.Gui) error {
		v, err := g.View(viewName)
		if err != nil {
			// handle error
		}
		v.Clear()
		for _, item := range data {
			fmt.Fprintln(v, item)
		}

		return nil
	})
}

func (r *Render) UpdateTextView(viewName string, data string) error {
	go r.TUI.Update(func(g *gocui.Gui) error {
		v, err := g.View(viewName)
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprintln(v, data)

		return nil
	})
	return nil
}

func (r *Render) CursorDown(viewName string) error {
	v, err := r.TUI.View(viewName)
	if err != nil {
		return err
	}
	if v != nil && r.GetNextLine(v) != "" {
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

func (r *Render) CursorUp(viewName string) error {
	v, err := r.TUI.View(viewName)
	if err != nil {
		return err
	}
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}

func (r *Render) ResetCursor(viewName string) error {
	v, err := r.TUI.View(viewName)
	if err != nil {
		return err
	}
	if err := v.SetCursor(0, 0); err != nil {
		return err
	}

	return nil
}

func (r *Render) getPanelByViewName(viewName string) *Panel {
	for i, p := range r.Panels {
		if p.Name == viewName {
			return r.Panels[i]
		}
	}

	return nil
}

func (r *Render) EnableSelection(viewName string) error {
	pan := r.getPanelByViewName(viewName)
	if pan == nil {
		return nil
	}
	pan.EnableSelection()
	return nil
}

func (r *Render) DisableSelection(viewName string) error {
	pan := r.getPanelByViewName(viewName)
	if pan == nil {
		return nil
	}
	pan.DisableSelection()
	return nil
}

func (r *Render) Quit() error {
	return gocui.ErrQuit
}

func (r *Render) GetCurrentLine(v *gocui.View) string {
	_, cursorY := v.Cursor()
	l, _ := v.Line(cursorY)

	return l
}

func (r *Render) GetNextLine(v *gocui.View) string {
	_, cursorY := v.Cursor()
	l, _ := v.Line(cursorY + 1)

	return l
}

func (r *Render) Focus(viewName string) error {
	if _, err := r.TUI.SetCurrentView(viewName); err != nil {
		log.Panicln("Try to focus on " + viewName)
		return err
	}

	return nil
}

func (r *Render) Show(viewName string) error {
	pan := r.getPanelByViewName(viewName)
	if pan == nil {
		return nil
	}
	pan.Hidden = false
	r.CreateViews()

	return nil
}

func (r *Render) Hide(viewName string) error {
	pan := r.getPanelByViewName(viewName)
	if pan == nil {
		return nil
	}
	pan.Hidden = true
	r.TUI.DeleteView(viewName)

	return nil
}
