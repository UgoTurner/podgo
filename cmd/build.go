// Package cmd wires the application components together and boots the TUI.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/sirupsen/logrus"
	"github.com/ugoturner/podgo/conf"
	"github.com/ugoturner/podgo/handler"
	"github.com/ugoturner/podgo/model"
	"github.com/ugoturner/podgo/service"
	"github.com/ugoturner/songocui"
)

// Build wires the application components together and starts the TUI. It
// returns an error instead of exiting so main can decide how to report it.
func Build() error {
	paths, err := conf.ResolvePaths()
	if err != nil {
		return fmt.Errorf("resolving paths: %w", err)
	}

	appLog, err := openLogFile(paths.AppLogFile())
	if err != nil {
		return err
	}
	defer func() { _ = appLog.Close() }()

	sgLog, err := openLogFile(paths.SongocuiLogFile())
	if err != nil {
		return err
	}
	defer func() { _ = sgLog.Close() }()

	appLogger := newLogger(appLog)
	songocuiLogger := newLogger(sgLog)

	db, err := scribble.New(paths.DbDir, nil)
	if err != nil {
		return fmt.Errorf("initializing database: %w", err)
	}

	cfg, err := loadUIConfig()
	if err != nil {
		return err
	}

	sg := songocui.NewWithLogger(songocuiLogger)

	appSubscriber := &handler.App{
		TUI:            sg,
		FeedRepository: &model.FeedRepository{Db: db, Logger: appLogger},
		FeedParser:     &service.FeedParser{Logger: appLogger},
		Player:         &service.Player{Logger: appLogger},
		TracksDir:      paths.TracksDir,
		Logger:         appLogger,
	}

	sg.RegisterSubscribers([]songocui.Subscriber{appSubscriber})

	if err := sg.Configure(cfg); err != nil {
		return fmt.Errorf("configuring TUI: %w", err)
	}
	return sg.Boot()
}

// loadUIConfig unmarshals the embedded panel and keybind definitions.
func loadUIConfig() (songocui.Config, error) {
	var panels []*songocui.Panel
	if err := json.Unmarshal(conf.PanelsJSON, &panels); err != nil {
		return songocui.Config{}, fmt.Errorf("parsing embedded panels config: %w", err)
	}
	var keybinds []songocui.KeybindGroup
	if err := json.Unmarshal(conf.KeybindsJSON, &keybinds); err != nil {
		return songocui.Config{}, fmt.Errorf("parsing embedded keybinds config: %w", err)
	}
	return songocui.Config{
		Panels:       panels,
		Keybinds:     keybinds,
		DefaultFocus: conf.SideViewName,
	}, nil
}

// openLogFile opens (creating if needed) a log file for appending.
func openLogFile(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("opening log file %q: %w", path, err)
	}
	return f, nil
}

// newLogger returns a JSON logrus logger writing to w.
func newLogger(w *os.File) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(w)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}
