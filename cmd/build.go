package cmd

import (
	"log"
	"os"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/sirupsen/logrus"
	"github.com/ugoturner/podgo/conf"
	"github.com/ugoturner/podgo/handler"
	"github.com/ugoturner/podgo/model"
	"github.com/ugoturner/podgo/service"
	"github.com/ugoturner/songocui"
)

// initLogger initializes a logrus logger with the specified log file.
func initLogger(logFile string) *logrus.Logger {
	logger := logrus.New()

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Unable to open log file %s: %v", logFile, err)
	}

	logger.SetOutput(file)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}

// initDb initializes the database and returns the Scribble driver.
func initDb() *scribble.Driver {
	db, err := scribble.New(conf.DbPath, nil)
	if err != nil {
		log.Fatalf("Error initializing DB: %v", err)
	}
	return db
}

// Build initializes the application components and starts the TUI.
func Build() {
	// Initialize loggers
	songocuiLogger := initLogger(conf.SgLogFile)
	appLogger := initLogger(conf.AppLogFile)

	// Initialize the Songocui TUI
	sg := songocui.NewWithLogger(songocuiLogger)
	sg.Configure(
		conf.ConfPanels,
		conf.ConfKeybinds,
		conf.SideViewName,
	)

	// Set up the application subscriber
	appSubscriber := &handler.App{
		TUI:            sg,
		FeedRepository: &model.FeedRepository{Db: initDb(), Logger: appLogger},
		FeedParser:     &service.FeedParser{Logger: appLogger},
		Player:         &service.Player{Logger: appLogger},
		Logger:         appLogger,
	}

	// Register the subscriber and boot the application
	sg.RegisterSubscribers([]songocui.Subscriber{appSubscriber})
	sg.Boot()
}
