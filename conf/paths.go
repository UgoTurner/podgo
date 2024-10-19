package conf

import (
	"fmt"
	"os"
	"path/filepath"
)

// Paths holds the resolved runtime directories used by the application. They
// live under standard per-user locations so podgo can be installed with
// `go install` and run from any working directory.
type Paths struct {
	TracksDir string // downloaded audio files
	DbDir     string // Scribble database root
	LogDir    string // application and TUI log files
}

// AppLogFile returns the path to the application log file.
func (p Paths) AppLogFile() string { return filepath.Join(p.LogDir, "app.log") }

// SongocuiLogFile returns the path to the Songocui (TUI) log file.
func (p Paths) SongocuiLogFile() string { return filepath.Join(p.LogDir, "songocui.log") }

// ResolvePaths computes the runtime directories following the XDG base
// directory conventions and ensures they exist.
//
//	tracks -> $XDG_CACHE_HOME/podgo/tracks   (default ~/.cache/podgo/tracks)
//	db     -> $XDG_DATA_HOME/podgo/db        (default ~/.local/share/podgo/db)
//	logs   -> $XDG_DATA_HOME/podgo/logs      (default ~/.local/share/podgo/logs)
func ResolvePaths() (Paths, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return Paths{}, fmt.Errorf("resolving user cache dir: %w", err)
	}
	dataDir, err := userDataDir()
	if err != nil {
		return Paths{}, fmt.Errorf("resolving user data dir: %w", err)
	}

	p := Paths{
		TracksDir: filepath.Join(cacheDir, "podgo", "tracks"),
		DbDir:     filepath.Join(dataDir, "podgo", "db"),
		LogDir:    filepath.Join(dataDir, "podgo", "logs"),
	}

	for _, dir := range []string{p.TracksDir, p.DbDir, p.LogDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return Paths{}, fmt.Errorf("creating directory %q: %w", dir, err)
		}
	}

	return p, nil
}

// userDataDir returns $XDG_DATA_HOME or falls back to ~/.local/share.
func userDataDir() (string, error) {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share"), nil
}
