package cmd

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/ugo/podgo/conf"
	"github.com/ugo/podgo/model"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/ugo/podgo/handler"
	"github.com/ugo/podgo/service"
	"github.com/ugo/sangocui"
)

/*
func init() {
	file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		log.SetFormatter(&log.JSONFormatter{})
		//the log level is set to "Info" which allows
		//to log Info(), Warn(), Error() and Fatal()
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// pkg log :

		file, err := os.OpenFile("pkg.log", os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			log.SetOutput(file)
			log.SetFormatter(&log.JSONFormatter{})
			//the log level is set to "Info" which allows
			//to log Info(), Warn(), Error() and Fatal()
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
}
*/

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

func Build() {
	// init db
	db, err := scribble.New("./storage/db/", nil)
	if err != nil {
		log.Panicln("Error when init DB", err)
	}

	sg := sangocui.NewWithLogger(initSangocuiLogger())
	sg.Configure(
		"conf/panels.json",
		"conf/keybinds.json",
		conf.SideViewName,
	)
	appSubscriber := &handler.App{
		TUI:            sg,
		FeedRepository: &model.FeedRepository{Db: db},
		FeedParser:     &service.FeedParser{},
		Player:         &service.Player{},
		Logger:         initSangocuiLogger(),
	}
	sg.RegisterSubscribers([]sangocui.Subscriber{appSubscriber})
	sg.Boot()
}
