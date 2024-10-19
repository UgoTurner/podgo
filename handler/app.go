// Package handler maps TUI events to application actions: fetching feeds,
// downloading tracks and controlling playback.
package handler

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
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
	TracksDir      string
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
	if err := a.RefreshAllFeeds(); err != nil {
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
	if err := a.TUI.CursorUp(conf.SideViewName); err != nil {
		return err
	}
	return a.updateUIAfterFeedChange()
}

// nextPodcast moves to the next podcast in the feed.
func (a *App) nextPodcast() error {
	a.FeedParser.NextFeed()
	if err := a.TUI.CursorDown(conf.SideViewName); err != nil {
		return err
	}
	return a.updateUIAfterFeedChange()
}

// enterTracksList focuses the UI on the list of tracks.
func (a *App) enterTracksList() error {
	if err := a.TUI.EnableSelection(conf.MainViewName); err != nil {
		return err
	}
	return a.TUI.Focus(conf.MainViewName)
}

// previousTrack moves to the previous track.
func (a *App) previousTrack() error {
	a.FeedParser.PrevItem()
	return a.TUI.CursorUp(conf.MainViewName)
}

// nextTrack moves to the next track.
func (a *App) nextTrack() error {
	a.FeedParser.NextItem()
	return a.TUI.CursorDown(conf.MainViewName)
}

// enterPodcastsList resets and focuses the UI on the podcast list.
func (a *App) enterPodcastsList() error {
	if err := a.TUI.ResetCursor(conf.MainViewName); err != nil {
		return err
	}
	if err := a.TUI.DisableSelection(conf.MainViewName); err != nil {
		return err
	}
	a.FeedParser.ResetFeedIdx()
	a.FeedParser.ResetItemIdx()
	return a.TUI.Focus(conf.SideViewName)
}

// enterTrackDescription displays the description of the current track.
func (a *App) enterTrackDescription() error {
	if err := a.TUI.Show(conf.MainDetailsViewName); err != nil {
		return err
	}
	if err := a.TUI.UpdateTextView(conf.MainDetailsViewName, a.FeedParser.GetCurrentItemDescription()); err != nil {
		return err
	}
	return a.TUI.Focus(conf.MainDetailsViewName)
}

// enterPodcastsListFromDescription hides the track description and focuses the podcast list.
func (a *App) enterPodcastsListFromDescription() error {
	if err := a.TUI.Hide(conf.MainDetailsViewName); err != nil {
		return err
	}
	return a.TUI.Focus(conf.MainViewName)
}

// playTrack attempts to play the current track, downloading it first if needed.
func (a *App) playTrack() error {
	trackName := fmt.Sprintf("%s - %s", a.FeedParser.GetCurrentFeedName(), a.FeedParser.GetCurrentItemName())
	fileName := a.FeedParser.GetCurrentItemLocalFileName()
	if fileName == "" {
		if err := a.TUI.UpdateTextView(conf.FooterViewName, "Track not downloaded yet."); err != nil {
			a.Logger.WithError(err).Warn("Failed to update footer view")
		}
		return a.downloadTrack(true)
	}

	trackPath := filepath.Join(a.TracksDir, fileName)
	go func() {
		if err := a.Player.Play(trackPath, trackName, func(status string) {
			if err := a.TUI.UpdateTextView(conf.FooterViewName, fmt.Sprintf("%s ~ %s", status, a.Player.PlayingTrackName)); err != nil {
				a.Logger.Errorf("Failed to update footer view: %v", err)
			}
		}); err != nil {
			a.Logger.WithError(err).Errorf("Failed to play track %q", trackName)
		}
	}()
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
		return a.TUI.UpdateTextView(conf.FooterViewName, "Already downloaded")
	}

	trackURL := a.FeedParser.GetCurrentItemURL()
	fileName := extractFileName(trackURL)
	if fileName == "" {
		return a.TUI.UpdateTextView(conf.FooterViewName, "Track has no downloadable URL")
	}
	trackName := fmt.Sprintf("%s - %s", a.FeedParser.GetCurrentFeedName(), a.FeedParser.GetCurrentItemName())
	if err := a.TUI.UpdateTextView(conf.FooterViewName, "Download will start..."); err != nil {
		a.Logger.WithError(err).Warn("Failed to update footer view")
	}

	destPath := filepath.Join(a.TracksDir, fileName)
	go func() {
		_ = service.DownloadFile(
			destPath,
			trackURL,
			func(progress string) {
				a.updateFooter(fmt.Sprintf("Downloading '%s' - %s", trackName, progress))
			},
			func() {
				a.updateFooter(fmt.Sprintf("Successfully downloaded '%s'!", trackName))
				a.FeedParser.SetCurrentItemLocalFileName(fileName)
				a.persistCurrentFeed()
				a.updateUIWithFeeds()
				if autoPlay {
					if err := a.playTrack(); err != nil {
						a.Logger.WithError(err).Error("Failed to auto-play track")
					}
				}
			},
			func() {
				a.Logger.Errorf("Failed to download '%s'", fileName)
				a.updateFooter(fmt.Sprintf("Failed to download '%s'", fileName))
			},
			a.Logger,
		)
	}()

	return nil
}

// RefreshAllFeeds re-fetches every stored feed, preserving the download state
// of tracks that were already downloaded.
func (a *App) RefreshAllFeeds() error {
	feeds, err := a.FeedRepository.FetchAll()
	if err != nil {
		a.Logger.WithError(err).Error("Unable to fetch feeds from the database")
		return nil
	}

	for _, feed := range feeds {
		updatedFeed := a.FeedParser.LoadFeedFromURL(feed.URL)
		if updatedFeed == nil {
			a.Logger.Errorf("Failed to refresh feed: %s", feed.Title)
			continue
		}

		a.preserveDownloadState(feed.Items, updatedFeed.Items)
		feed.Title = updatedFeed.Title
		feed.Items = updatedFeed.Items

		if err := a.FeedRepository.Update([]*model.Feed{feed}); err != nil {
			a.Logger.Errorf("Error saving updated feed: %s, %v", feed.Title, err)
		} else {
			a.Logger.Infof("Successfully refreshed feed: %s", feed.Title)
		}
	}

	return nil
}

// preserveDownloadState carries the local file name of already-downloaded
// tracks from the old item list onto the freshly fetched one, dropping entries
// whose downloaded file no longer exists on disk.
func (a *App) preserveDownloadState(oldItems, newItems []*model.Item) {
	downloaded := make(map[string]string, len(oldItems))
	for _, item := range oldItems {
		if item.LocalFileName == "" {
			continue
		}
		if _, err := os.Stat(filepath.Join(a.TracksDir, item.LocalFileName)); err != nil {
			continue
		}
		downloaded[item.URL] = item.LocalFileName
	}
	for _, item := range newItems {
		if name, ok := downloaded[item.URL]; ok {
			item.LocalFileName = name
		}
	}
}

// addNewFeed prepares the UI for adding a new feed.
func (a *App) addNewFeed() error {
	if err := a.TUI.Show(conf.PromptViewName); err != nil {
		return err
	}
	return a.TUI.Focus(conf.PromptViewName)
}

// confirmNewFeed adds a new feed based on user input.
func (a *App) confirmNewFeed() error {
	feedURL := strings.TrimSpace(a.TUI.GetCurrentBuffer(conf.PromptViewName))
	feed := a.FeedParser.LoadFeedFromURL(feedURL)
	if feed == nil {
		return a.triggerInfoMessage("Invalid URL", logrus.Fields{"url": feedURL})
	}
	if err := a.FeedRepository.Update([]*model.Feed{feed}); err != nil {
		a.Logger.WithError(err).Error("Failed to store new feed")
	}
	a.FeedParser.AddFeed(feed)
	if err := a.TUI.Hide(conf.PromptViewName); err != nil {
		return err
	}
	if err := a.TUI.Focus(conf.SideViewName); err != nil {
		return err
	}
	return a.launch()
}

// quitNewFeed cancels adding a new feed.
func (a *App) quitNewFeed() error {
	if err := a.TUI.Hide(conf.PromptViewName); err != nil {
		return err
	}
	return a.TUI.Focus(conf.SideViewName)
}

// Helper methods

// extractFileName derives a filesystem-safe file name from a track URL,
// stripping any query string so the name is portable across filesystems.
func extractFileName(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	if u, err := url.Parse(rawURL); err == nil && u.Path != "" {
		return path.Base(u.Path)
	}
	if i := strings.IndexAny(rawURL, "?#"); i >= 0 {
		rawURL = rawURL[:i]
	}
	return path.Base(rawURL)
}

// persistCurrentFeed saves the current feed (including download state) to the DB.
func (a *App) persistCurrentFeed() {
	feed := a.FeedParser.GetCurrentFeed()
	if feed == nil {
		return
	}
	if err := a.FeedRepository.Update([]*model.Feed{feed}); err != nil {
		a.Logger.WithError(err).Error("Failed to persist download state")
	}
}

// updateFooter updates the footer view, logging any failure.
func (a *App) updateFooter(msg string) {
	if err := a.TUI.UpdateTextView(conf.FooterViewName, msg); err != nil {
		a.Logger.WithError(err).Warn("Failed to update footer view")
	}
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
	if err := a.TUI.Show(conf.FlashMessageViewName); err != nil {
		return err
	}
	if err := a.TUI.UpdateTextView(conf.FlashMessageViewName, msg); err != nil {
		return err
	}
	a.Logger.WithFields(fields).Info(msg)

	go func() {
		time.Sleep(2 * time.Second)
		if err := a.TUI.Hide(conf.FlashMessageViewName); err != nil {
			a.Logger.WithError(err).Warn("Failed to hide flash message")
		}
	}()

	return nil
}
