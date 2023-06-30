package cmd

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/ugoturner/podgo/conf"
	"github.com/ugoturner/podgo/model"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/ugoturner/podgo/handler"
	"github.com/ugoturner/podgo/service"
	"github.com/ugoturner/songocui"
)

func initSongocuiLogger() *logrus.Logger {
	var logger = logrus.New()
	file, err := os.OpenFile(conf.SgLogFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Panicln(err)
	}
	logger.SetOutput(file)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}

func initAppLogger() *logrus.Logger {
	var logger = logrus.New()
	file, err := os.OpenFile(conf.AppLogFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Panicln(err)
	}
	logger.SetOutput(file)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}

func initDb() *scribble.Driver {
	db, err := scribble.New(conf.DbPath, nil)
	if err != nil {
		log.Panicln("Error during init DB", err)
	}

	return db
}

func Build() {
	sg := songocui.NewWithLogger(initSongocuiLogger())
	sg.Configure(
		conf.ConfPanels,
		conf.ConfKeybinds,
		conf.SideViewName,
	)
	appSubscriber := &handler.App{
		TUI:            sg,
		FeedRepository: &model.FeedRepository{Db: initDb()},
		FeedParser:     &service.FeedParser{},
		Player:         &service.Player{},
		Logger:         initSongocuiLogger(),
	}
	sg.RegisterSubscribers([]songocui.Subscriber{appSubscriber})
	sg.Boot()
}
