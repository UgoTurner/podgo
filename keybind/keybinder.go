package keybind

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/ugo/podcastor/event"

	"github.com/jroimartin/gocui"
)

func gocuiAdapt(keyStr string) gocui.Key {
	switch keyStr {
	case "ctrlC":
		return gocui.KeyCtrlC
	case "arrowUp":
		return gocui.KeyArrowUp
	case "arrowDown":
		return gocui.KeyArrowDown
	case "arrowRight":
		return gocui.KeyArrowRight
	case "arrowLeft":
		return gocui.KeyArrowLeft
	case "ctrlD":
		return gocui.KeyCtrlD
	case "ctrlP":
		return gocui.KeyCtrlP
	case "ctrlSpace":
		return gocui.KeyCtrlSpace
	case "ctrlF":
		return gocui.KeyCtrlF
	default:
		log.Panicln("Unkown keybind : " + keyStr)
		return gocui.KeyCtrl2
	}
}

type bindingCallbackFn func(action string) error

type Keybinder struct {
	EventDispatcher *event.Dispatcher
	TUI             *gocui.Gui
}

type confViewsKeybind struct {
	ViewName string
	Keybinds []confKeybind
}

type confKeybind struct {
	Key    string
	Action string
}

func (k *Keybinder) LoadKeybinds(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}
	viewsKeybinds := []confViewsKeybind{}
	json.Unmarshal(data, &viewsKeybinds)
	for _, v := range viewsKeybinds {
		k.assigndKeybinds(v)
	}
}

func (k *Keybinder) assigndKeybinds(ckb confViewsKeybind) {
	for _, kb := range ckb.Keybinds {
		keybind := NewKeybind(kb.Key, kb.Action)
		k.TUI.SetKeybinding(
			ckb.ViewName,
			keybind.Key,
			gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return k.EventDispatcher.Dispatch(keybind.Action)
			},
		)
	}

}
