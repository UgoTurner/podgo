package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ugoturner/podgo/conf"
	"github.com/ugoturner/podgo/model"
	"github.com/ugoturner/podgo/service"
	"github.com/ugoturner/songocui"
)

// App binds songocui events and triggers actions like fetching feeds, downloading, and playing.
type App struct {
	songocui.Subscriber
	TUI            *songocui.Songocui
	FeedRepository *model.FeedRepository
	FeedParser     *service.FeedParser
	Player         *service.Player
	Logger         *logrus.Logger
}

// On listens to songocui events and triggers the corresponding actions.
func (a *App) On(eventName string) error {
	eventHandlers := map[string]func() error{
		"Launch":                           a.launch,
		"Shutdown":                         a.shutdown,
		"PreviousPodcast":                  a.previousPodcast,
		"NextPodcast":                      a.nextPodcast,
		"EnterTracksList":                  a.enterTracksList,
		"PreviousTrack":                    a.previousTrack,
		"NextTrack":                        a.nextTrack,
		"EnterPodcastsList":                a.enterPodcastsList,
		"DownloadTrack":                    func() error { return a.downloadTrack(false) },
		"EnterTrackDescription":            a.enterTrackDescription,
		"EnterPodcastsListFromDescription": a.enterPodcastsListFromDescription,
		"PlayTrack":                        a.playTrack,
		"TogglePlayPause":                  a.togglePlayPause,
		"SeekForward":                      func() error { return a.seek(10) },
		"SeekBackward":                     func() error { return a.seek(-10) },
		"AddNewFeed":                       a.addNewFeed,
		"ConfirmNewFeed":                   a.confirmNewFeed,
		"QuitNewFeed":                      a.quitNewFeed,
	}

	if handler, ok := eventHandlers[eventName]; ok {
		return handler()
	}

	a.Logger.Warnf("Unhandled event: %s", eventName)
	return nil
}

// launch initializes the feed and updates the UI.
func (a *App) launch() error {
	err := a.RefreshAllFeeds()
	if err != nil {
		a.Logger.WithError(err).Error("Error refreshing feeds")	
	}
	feeds, err := a.FeedRepository.FetchAll()
	if err != nil {
		a.Logger.WithError(err).Error("Error fetching feeds")
		return err
	}
	a.FeedParser.SetFeeds(feeds)
	a.updateUIWithFeeds()
	return nil
}

// shutdown quits the TUI.
func (a *App) shutdown() error {
	return a.TUI.Quit()
}

// previousPodcast moves to the previous podcast in the feed.
func (a *App) previousPodcast() error {
	a.FeedParser.PrevFeed()
	a.TUI.CursorUp(conf.SideViewName)
	return a.updateUIAfterFeedChange()
}

// nextPodcast moves to the next podcast in the feed.
func (a *App) nextPodcast() error {
	a.FeedParser.NextFeed()
	a.TUI.CursorDown(conf.SideViewName)
	return a.updateUIAfterFeedChange()
}

// enterTracksList focuses the UI on the list of tracks.
func (a *App) enterTracksList() error {
	a.TUI.EnableSelection(conf.MainViewName)
	a.TUI.Focus(conf.MainViewName)
	return nil
}

// previousTrack moves to the previous track.
func (a *App) previousTrack() error {
	a.FeedParser.PrevItem()
	a.TUI.CursorUp(conf.MainViewName)
	return nil
}

// nextTrack moves to the next track.
func (a *App) nextTrack() error {
	a.FeedParser.NextItem()
	a.TUI.CursorDown(conf.MainViewName)
	return nil
}

// enterPodcastsList resets and focuses the UI on the podcast list.
func (a *App) enterPodcastsList() error {
	a.TUI.ResetCursor(conf.MainViewName)
	a.TUI.DisableSelection(conf.MainViewName)
	a.FeedParser.ResetFeedIdx()
	a.FeedParser.ResetItemIdx()
	a.TUI.Focus(conf.SideViewName)
	return nil
}

// enterTrackDescription displays the description of the current track.
func (a *App) enterTrackDescription() error {
	a.TUI.Show(conf.MainDetailsViewName)
	a.TUI.UpdateTextView(conf.MainDetailsViewName, a.FeedParser.GetCurrentItemDescription())
	a.TUI.Focus(conf.MainDetailsViewName)
	return nil
}

// enterPodcastsListFromDescription hides the track description and focuses the podcast list.
func (a *App) enterPodcastsListFromDescription() error {
	a.TUI.Hide(conf.MainDetailsViewName)
	a.TUI.Focus(conf.MainViewName)
	return nil
}

// playTrack attempts to play the current track.
func (a *App) playTrack() error {
	trackName := fmt.Sprintf("%s - %s", a.FeedParser.GetCurrentFeedName(), a.FeedParser.GetCurrentItemName())
	fileName := a.FeedParser.GetCurrentItemLocalFileName()
	if fileName == "" {
		a.TUI.UpdateTextView(conf.FooterViewName, "Track not downloaded yet.")
		return a.downloadTrack(true)
	}

	trackPath := conf.TracksPath + fileName
	go a.Player.Play(trackPath, trackName, func(status string) {
		if err := a.TUI.UpdateTextView(conf.FooterViewName, fmt.Sprintf("%s ~ %s", status, a.Player.PlayingTrackName)); err != nil {
			a.Logger.Errorf("Failed to update footer view: %v", err)
		}
	})
	return nil
}

// togglePlayPause toggles the play/pause state of the player.
func (a *App) togglePlayPause() error {
	a.Player.TogglePlayPause()
	return nil
}

// seek moves the current playback position forward or backward.
func (a *App) seek(seconds int) error {
	a.Player.Seek(seconds)
	return nil
}

// downloadTrack downloads the current track and optionally plays it afterward.
func (a *App) downloadTrack(autoPlay bool) error {
	if a.FeedParser.GetCurrentItemLocalFileName() != "" {
		a.TUI.UpdateTextView(conf.FooterViewName, "Already downloaded")
		return nil
	}

	fileName := a.extractFileName(a.FeedParser.GetCurrentItemUrl())
	trackName := fmt.Sprintf("%s - %s", a.FeedParser.GetCurrentFeedName(), a.FeedParser.GetCurrentItemName())
	a.TUI.UpdateTextView(conf.FooterViewName, "Download will start...")

	// Start the download in a goroutine.
	go service.DownloadFile(
		conf.TracksPath+fileName,
		a.FeedParser.GetCurrentItemUrl(),
		func(progress string) {
			a.TUI.UpdateTextView(conf.FooterViewName, fmt.Sprintf("Downloading '%s' - %s", trackName, progress))
		},
		func() {
			a.TUI.UpdateTextView(conf.FooterViewName, fmt.Sprintf("Successfully downloaded '%s'!", trackName))
			a.FeedParser.SetCurrentItemLocalFileName(fileName)
			a.updateUIWithFeeds()
			if autoPlay {
				a.playTrack()
			}
		},
		func() {
			a.Logger.Errorf("Failed to download '%s'", fileName)
			a.TUI.UpdateTextView(conf.FooterViewName, fmt.Sprintf("Failed to download '%s'", fileName))
		},
		a.Logger,
	)

	return nil
}

func (a *App) RefreshAllFeeds() error {
	// Retrieve all feeds from the FeedRepository
	feeds, err := a.FeedRepository.FetchAll()
	if err != nil {
		a.Logger.Error("No feeds found in the database")
		return nil
	}

	for _, feed := range feeds {
		updatedFeed := a.FeedParser.LoadFeedFromUrl(feed.Url)
		if updatedFeed == nil {
			a.Logger.Errorf("Failed to refresh feed: %s", feed.Title)
			continue
		}

		feed.Title = updatedFeed.Title
		feed.Items = updatedFeed.Items

		if err := a.FeedRepository.Update([]*model.Feed{feed}); err != nil {
			a.Logger.Errorf("Error saving updated feed: %s, %v", feed.Title, err)
		} else {
			a.Logger.Errorf("Successfully refreshed feed: %s", feed.Title)
		}
	}

	return nil
}

// addNewFeed prepares the UI for adding a new feed.
func (a *App) addNewFeed() error {
	a.TUI.Show(conf.PromptViewName)
	a.TUI.Focus(conf.PromptViewName)
	return nil
}

// confirmNewFeed adds a new feed based on user input.
func (a *App) confirmNewFeed() error {
	url := strings.TrimSpace(a.TUI.GetCurrentBuffer(conf.PromptViewName))
	feed := a.FeedParser.LoadFeedFromUrl(url)
	if feed == nil {
		return a.triggerInfoMessage("Invalid URL", logrus.Fields{"url": url})
	}
	a.FeedRepository.Update([]*model.Feed{feed})
	a.FeedParser.AddFeed(feed)
	a.TUI.Hide(conf.PromptViewName)
	a.TUI.Focus(conf.SideViewName)
	return a.launch()
}

// quitNewFeed cancels adding a new feed.
func (a *App) quitNewFeed() error {
	a.TUI.Hide(conf.PromptViewName)
	a.TUI.Focus(conf.SideViewName)
	return nil
}

// Helper methods

// extractFileName extracts the file name from a URL.
func (a *App) extractFileName(url string) string {
	tokens := strings.Split(url, "/")
	return tokens[len(tokens)-1]
}

// updateUIWithFeeds updates the UI with the current feed and item list.
func (a *App) updateUIWithFeeds() {
	a.TUI.UpdateListView(conf.SideViewName, a.FeedParser.GetFeedNames())
	a.TUI.UpdateListView(conf.MainViewName, a.FeedParser.GetCurrentFeedItemsNameAndStatus())
}

// updateUIAfterFeedChange updates the UI after changing feeds.
func (a *App) updateUIAfterFeedChange() error {
	a.TUI.UpdateListView(conf.MainViewName, a.FeedParser.GetCurrentFeedItemsNameAndStatus())
	return nil
}

// triggerInfoMessage displays a temporary info message in the UI.
func (a *App) triggerInfoMessage(msg string, fields logrus.Fields) error {
	a.TUI.Show(conf.FlashMessageViewName)
	a.TUI.UpdateTextView(conf.FlashMessageViewName, msg)

	go func() {
		time.Sleep(2 * time.Second)
		a.TUI.Hide(conf.FlashMessageViewName)
	}()

	return nil
}
