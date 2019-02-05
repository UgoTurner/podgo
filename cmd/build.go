package cmd

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/uturner/podgo/conf"
	"github.com/uturner/podgo/model"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/uturner/podgo/handler"
	"github.com/uturner/podgo/service"
	"github.com/uturner/sangocui"
)

func initSangocuiLogger() *logrus.Logger {
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
	sg := sangocui.NewWithLogger(initSangocuiLogger())
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
		Logger:         initSangocuiLogger(),
	}
	sg.RegisterSubscribers([]sangocui.Subscriber{appSubscriber})
	sg.Boot()
}
