package conf

import _ "embed"

// The UI configuration is embedded into the binary at build time, so the app
// needs no external files at runtime and can be installed with `go install`
// and run from anywhere.

// PanelsJSON is the embedded panel layout configuration.
//
//go:embed panels.json
var PanelsJSON []byte

// KeybindsJSON is the embedded keybinding configuration.
//
//go:embed keybinds.json
var KeybindsJSON []byte
