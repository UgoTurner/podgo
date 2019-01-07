package utils

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

// List utils :
func displayItems(v *gocui.View, items []string) *gocui.View {
  for _, item := range items {
		fmt.Fprintln(v, item)
  }

  return v
}

func UpdateList(v *gocui.View, items []string) *gocui.View {
  v.Clear()
  displayItems(v, items)

  return v
}
