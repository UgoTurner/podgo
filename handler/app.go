package handler

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/uturner/podgo/conf"
	"github.com/uturner/podgo/model"
	"github.com/uturner/podgo/service"
	"github.com/uturner/songocui"
)

// App : Binds songocui event and triggers actions (fetch feeds, dl, play ...)
type App struct {
	songocui.Subscriber
	TUI            *songocui.Songocui
	FeedRepository *model.FeedRepository
	FeedParser     *service.FeedParser
	Player         *service.Player
	Logger         *logrus.Logger
}

// On : Method from "Subscriber" interface, triggered by songocui.
// Call the action matching an event name
func (a *App) On(eventName string) error {
	switch eventName {
	case "Launch":
		return a.launch()
	case "Shutdown":
		return a.shutdown()
	case "PreviousPodcast":
		return a.previousPodcast()
	case "NextPodcast":
		return a.nextPodcast()
	case "EnterTracksList":
		return a.enterTracksList()
	case "PreviousTrack":
		return a.previousTrack()
	case "NextTrack":
		return a.nextTrack()
	case "EnterPodcastsList":
		return a.enterPodcastsList()
	case "DownloadTrack":
		return a.downloadTrack(false)
	case "EnterTrackDescription":
		return a.enterTrackDescription()
	case "EnterPodcastsListFromDescription":
		return a.enterPodcastsListFromDescription()
	case "PlayTrack":
		return a.playTrack()
	case "TogglePlayPause":
		return a.togglePlayPause()
	case "SeekForward":
		return a.seekForward()
	case "SeekBackward":
		return a.seekBackward()
	case "AddNewFeed":
		return a.addNewFeed()
	case "ConfirmNewFeed":
		return a.confirmNewFeed()
	case "QuitNewFeed":
		return a.quitNewFeed()
	default:
		return nil
	}

}

func (a *App) launch() error {
	feeds := a.FeedRepository.FetchAll()
	a.FeedParser.SetFeeds(feeds)
	a.TUI.UpdateListView(conf.SideViewName, a.FeedParser.GetFeedNames())
	a.TUI.UpdateListView(conf.MainViewName, a.FeedParser.GetCurrentFeedItemsNameAndStatus())
	return nil
}

func (a *App) shutdown() error {
	return a.TUI.Quit()
}

func (a *App) previousPodcast() error {
	a.FeedParser.PrevFeed()
	a.TUI.CursorUp(conf.SideViewName)
	a.TUI.UpdateListView(
		conf.MainViewName,
		a.FeedParser.GetCurrentFeedItemsNameAndStatus(),
	)

	return nil
}

func (a *App) nextPodcast() error {
	a.FeedParser.NextFeed()
	a.TUI.CursorDown(conf.SideViewName)
	a.TUI.UpdateListView(
		conf.MainViewName,
		a.FeedParser.GetCurrentFeedItemsNameAndStatus(),
	)

	return nil
}

func (a *App) enterTracksList() error {
	a.TUI.EnableSelection(conf.MainViewName)
	a.TUI.Focus(conf.MainViewName)

	return nil
}

func (a *App) previousTrack() error {
	a.FeedParser.PrevItem()
	a.TUI.CursorUp(conf.MainViewName)

	return nil
}

func (a *App) nextTrack() error {
	a.FeedParser.NextItem()
	a.TUI.CursorDown(conf.MainViewName)

	return nil
}

func (a *App) enterPodcastsList() error {
	a.TUI.ResetCursor(conf.MainViewName)
	a.TUI.DisableSelection(conf.MainViewName)
	a.FeedParser.ResetFeedIdx()
	a.FeedParser.ResetItemIdx()
	a.TUI.Focus(conf.SideViewName)

	return nil
}

func (a *App) extractFileName(url string) string {
	tokens := strings.Split(url, "/")

	return tokens[len(tokens)-1]

}

func (a *App) downloadTrack(autoPlay bool) error {
	if a.FeedParser.GetCurrentItemLocalFileName() != "" {
		a.TUI.UpdateTextView(
			conf.FooterViewName,
			"Already downloaded",
		)
		return nil
	}
	fileName := a.extractFileName(a.FeedParser.GetCurrentItemUrl())
	a.TUI.UpdateTextView(
		conf.FooterViewName,
		"Download will start...",
	)
	go service.DownloadFile(
		conf.TracksPath+fileName,
		a.FeedParser.GetCurrentItemUrl(),
		func(progress string) {
			a.TUI.UpdateTextView(
				conf.FooterViewName,
				"Downloading '"+fileName+"' - "+progress,
			)
		},
		func() {
			a.TUI.UpdateTextView(
				conf.FooterViewName,
				"Successfully download '"+fileName+"' !",
			)
			a.FeedParser.SetCurrentItemLocalFileName(fileName)
			a.TUI.UpdateListView(
				conf.MainViewName,
				a.FeedParser.GetCurrentFeedItemsNameAndStatus(),
			)
			if autoPlay {
				a.playTrack()
			}
		},
		func() {
			a.TUI.UpdateTextView(
				conf.FooterViewName,
				"Fail to donwload '"+fileName+"' !",
			)
		},
	)

	return nil
}

func (a *App) enterTrackDescription() error {
	a.TUI.Show(conf.MainDetailsViewName)
	a.TUI.Focus(conf.MainDetailsViewName)
	a.TUI.UpdateTextView(
		conf.MainDetailsViewName,
		a.FeedParser.GetCurrentItemDescription(),
	)

	return nil

}

func (a *App) enterPodcastsListFromDescription() error {
	a.TUI.Hide(conf.MainDetailsViewName)
	a.TUI.Focus(conf.MainViewName)

	return nil
}

func (a *App) playTrack() error {
	track := a.FeedParser.GetCurrentFeedName() + " - " + a.FeedParser.GetCurrentItemName()
	fileName := a.FeedParser.GetCurrentItemLocalFileName()
	if fileName == "" {
		a.TUI.UpdateTextView(
			conf.FooterViewName,
			"Track not downloaded yet.",
		)
		a.downloadTrack(true)

		return nil
	}
	path := conf.TracksPath + fileName
	go a.Player.Play(
		path,
		track,
		func(s string) {
			a.TUI.UpdateTextView(
				conf.FooterViewName,
				s+" ~ "+a.Player.PlayingTrackName,
			)
		},
	)
	return nil
}

func (a *App) togglePlayPause() error {
	a.Player.TogglePlayPause()

	return nil
}

func (a *App) seekForward() error {
	a.Player.Seek(10)

	return nil
}

func (a *App) seekBackward() error {
	a.Player.Seek(-10)

	return nil
}

func (a *App) addNewFeed() error {
	a.TUI.Show(conf.PromptViewName)
	a.TUI.Focus(conf.PromptViewName)

	return nil
}

func (a *App) triggerInfoMessage(msg string) error {
	go func() error {
		a.TUI.Show(conf.FlashMessageViewName)
		a.TUI.UpdateTextView(
			conf.FlashMessageViewName,
			msg,
		)
		time.Sleep(2 * time.Second)
		a.TUI.Hide(conf.FlashMessageViewName)

		return nil
	}()

	return nil
}

func (a *App) confirmNewFeed() error {
	url := strings.TrimSpace(a.TUI.GetCurrentBuffer(conf.PromptViewName))
	feed := a.FeedParser.LoadFeedFromUrl(url)
	if feed == nil {
		a.triggerInfoMessage("Invalid url")
		a.Logger.WithFields(logrus.Fields{"url": url}).Error("Can't parse feed source")
		return nil
	}
	a.FeedRepository.Update([]*model.Feed{feed})
	a.FeedParser.AddFeed(feed)
	a.TUI.Hide(conf.PromptViewName)
	a.TUI.Focus(conf.SideViewName)
	a.launch()

	return nil
}

func (a *App) quitNewFeed() error {
	a.TUI.Hide(conf.PromptViewName)
	a.TUI.Focus(conf.SideViewName)

	return nil
}
