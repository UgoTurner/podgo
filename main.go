// Command podgo is a terminal-based app for managing and playing podcasts.
package main

import (
	"fmt"
	"os"

	"github.com/ugoturner/podgo/cmd"
)

func main() {
	if err := cmd.Build(); err != nil {
		fmt.Fprintln(os.Stderr, "podgo:", err)
		os.Exit(1)
	}
}
