package handler

import (
	"strings"

	"github.com/ugo/podcastor/conf"
	"github.com/ugo/podcastor/service"
	"github.com/ugo/podcastor/ui"
)

type Handler interface {
	Handle(eventName string) error
}

type AppHandler struct {
	Handler
	FeedParser *service.FeedParser
	Render     *ui.Render
	Player     *service.Player
}

func (a *AppHandler) Handle(eventName string) error {
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
		return a.downloadTrack()
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
	default:
		return nil
	}

}

func (a *AppHandler) launch() error {
	a.FeedParser.LoadFeeds()
	a.Render.UpdateListView(conf.SideViewName, a.FeedParser.GetFeedNames())
	a.Render.UpdateListView(conf.MainViewName, a.FeedParser.GetCurrentFeedItemsNameAndStatus())
	return nil
}

func (a *AppHandler) shutdown() error {
	return a.Render.Quit()
}

func (a *AppHandler) previousPodcast() error {
	a.FeedParser.PrevFeed()
	a.Render.CursorUp(conf.SideViewName)
	a.Render.UpdateListView(
		conf.MainViewName,
		a.FeedParser.GetCurrentFeedItemsNameAndStatus(),
	)

	return nil
}

func (a *AppHandler) nextPodcast() error {
	a.FeedParser.NextFeed()
	a.Render.CursorDown(conf.SideViewName)
	a.Render.UpdateListView(
		conf.MainViewName,
		a.FeedParser.GetCurrentFeedItemsNameAndStatus(),
	)

	return nil
}

func (a *AppHandler) enterTracksList() error {
	a.Render.EnableSelection(conf.MainViewName)
	a.Render.Focus(conf.MainViewName)

	return nil
}

func (a *AppHandler) previousTrack() error {
	a.FeedParser.PrevItem()
	a.Render.CursorUp(conf.MainViewName)

	return nil
}

func (a *AppHandler) nextTrack() error {
	a.FeedParser.NextItem()
	a.Render.CursorDown(conf.MainViewName)

	return nil
}

func (a *AppHandler) enterPodcastsList() error {
	a.Render.ResetCursor(conf.MainViewName)
	a.Render.DisableSelection(conf.MainViewName)
	a.FeedParser.ResetFeedIdx()
	a.FeedParser.ResetItemIdx()
	a.Render.Focus(conf.SideViewName)

	return nil
}

func (a *AppHandler) extractFileName(url string) string {
	tokens := strings.Split(url, "/")

	return tokens[len(tokens)-1]

}

func (a *AppHandler) downloadTrack() error {
	if a.FeedParser.GetCurrentItemLocalFileName() != "" {
		a.Render.UpdateTextView(
			conf.FooterViewName,
			"Already downloaded",
		)
		return nil
	}
	fileName := a.extractFileName(a.FeedParser.GetCurrentItemUrl())
	a.Render.UpdateTextView(
		conf.FooterViewName,
		"Download will start...",
	)
	go service.DownloadFile(
		conf.TracksPath+fileName,
		a.FeedParser.GetCurrentItemUrl(),
		func(progress string) {
			a.Render.UpdateTextView(
				conf.FooterViewName,
				"Downloading '"+fileName+"' - "+progress,
			)
		},
		func() {
			a.Render.UpdateTextView(
				conf.FooterViewName,
				"Successfully donwload '"+fileName+"' !",
			)
			a.FeedParser.SetCurrentItemLocalFileName(fileName)
			a.Render.UpdateListView(
				conf.MainViewName,
				a.FeedParser.GetCurrentFeedItemsNameAndStatus(),
			)
		},
		func() {
			a.Render.UpdateTextView(
				conf.FooterViewName,
				"Fail to donwload '"+fileName+"' !",
			)
		},
	)

	return nil
}

func (a *AppHandler) enterTrackDescription() error {
	a.Render.Show(conf.MainDetailsViewName)
	a.Render.Focus(conf.MainDetailsViewName)
	a.Render.UpdateTextView(
		conf.MainDetailsViewName,
		a.FeedParser.GetCurrentItemDescription(),
	)

	return nil

}

func (a *AppHandler) enterPodcastsListFromDescription() error {
	a.Render.Hide(conf.MainDetailsViewName)
	a.Render.Focus(conf.MainViewName)

	return nil
}

func (a *AppHandler) playTrack() error {
	a.Player.Play(
		conf.TracksPath+a.FeedParser.GetCurrentItemLocalFileName(),
		func(s string) {
			a.Render.UpdateTextView(
				conf.FooterViewName,
				s,
			)
		},
	)

	return nil
}

func (a *AppHandler) togglePlayPause() error {
	a.Player.TogglePlayPause()

	return nil
}

func (a *AppHandler) seekForward() error {
	a.Player.Seek(10)

	return nil
}
