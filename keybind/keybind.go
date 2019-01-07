package keybind

import "github.com/jroimartin/gocui"

type Keybind struct {
	Key    gocui.Key
	Action string
}

func NewKeybind(keyStr string, action string) *Keybind {
	return &Keybind{
		Key:    gocuiAdapt(keyStr),
		Action: action,
	}
}
