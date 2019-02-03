package service

import (
	"github.com/mmcdole/gofeed"
	"github.com/ugo/podgo/model"
)

//"https://distorsionpodcast.com/podcasts.xml"
//"http://feeds.feedburner.com/bureaudesmysteres?fmt=xml"
//"http://feeds.soundcloud.com/playlists/soundcloud:playlists:253185017/sounds.rss"

type FeedParser struct {
	Feeds            []*model.Feed
	CurrentFeedIndex int
	CurrentItemIndex int
}

func (fh *FeedParser) GetFeeds() []*model.Feed {
	return fh.Feeds
}

func (fh *FeedParser) SetFeeds(feeds []*model.Feed) {
	fh.Feeds = feeds
}

func (fh *FeedParser) AddFeed(feed *model.Feed) {
	fh.Feeds = append(fh.Feeds, feed)
}

func (fh *FeedParser) extractItems(gItems []*gofeed.Item) []*model.Item {
	var it []*model.Item
	for _, gItem := range gItems {
		it = append(
			it,
			&model.Item{
				Title:       gItem.Title,
				Description: gItem.Description,
				Url:         fh.extractItemLink(gItem),
			})
	}

	return it
}

func (fh *FeedParser) extractFeed(gFeed *gofeed.Feed) *model.Feed {
	return &model.Feed{
		Title: gFeed.Title,
		Items: fh.extractItems(gFeed.Items),
	}
}

func (fh *FeedParser) extractItemLink(gItem *gofeed.Item) string {
	if gItem == nil {
		return ""
	}
	for _, en := range gItem.Enclosures {
		if en.Type != "audio/mpeg" {
			continue
		}

		return en.URL
	}
	return ""

}

func (fh *FeedParser) LoadFeedFromUrl(url string) *model.Feed {
	fp := gofeed.NewParser()
	f, err := fp.ParseURL(url)
	if err != nil {
		return nil
	}
	feed := fh.extractFeed(f)
	feed.Url = url

	return feed
}

func (fh *FeedParser) GetFeedNames() []string {
	var titles []string
	for _, feed := range fh.Feeds {
		titles = append(titles, feed.Title)
	}

	return titles
}

func (fh *FeedParser) GetItemNames() []string {
	var titles []string
	for _, feed := range fh.Feeds {
		for _, item := range feed.Items {
			titles = append(titles, item.Title)
		}
	}

	return titles
}

func (fh *FeedParser) getFeedByName(feedName string) *model.Feed {
	for _, feed := range fh.Feeds {
		if feed.Title != feedName {
			continue
		}
		return feed
	}

	return nil
}

func (fh *FeedParser) GetItemNamesByFeedName(feedName string) []string {
	feed := fh.getFeedByName(feedName)
	if feed == nil {
		return nil
	}
	var titles []string
	for _, item := range feed.Items {
		titles = append(titles, item.Title)
	}

	return titles
}

func (fh *FeedParser) getFeedByIndex(i int) *model.Feed {
	if i >= 0 && i < len(fh.Feeds) {
		return fh.Feeds[i]
	}

	return nil
}

func (fh *FeedParser) getItemByIndex(i int) *model.Item {
	if i >= 0 && i < len(fh.getCurrentFeed().Items) {
		return fh.getCurrentFeed().Items[i]
	}
	return nil
}

func (fh *FeedParser) getCurrentFeed() *model.Feed {
	return fh.getFeedByIndex(fh.CurrentFeedIndex)

}

func (fh *FeedParser) getCurrentItem() *model.Item {
	return fh.getItemByIndex(fh.CurrentItemIndex)
}

func (fh *FeedParser) GetCurrentFeedName() string {
	if fh.getCurrentFeed() == nil {
		return ""
	}

	return fh.getCurrentFeed().Title

}

func (fh *FeedParser) NextFeed() {
	if fh.CurrentFeedIndex < len(fh.Feeds)-1 {
		fh.CurrentFeedIndex = fh.CurrentFeedIndex + 1
	}
}

func (fh *FeedParser) PrevFeed() {
	if fh.CurrentFeedIndex > 0 {
		fh.CurrentFeedIndex = fh.CurrentFeedIndex - 1
	}
}

func (fh *FeedParser) NextItem() {
	if fh.CurrentItemIndex < len(fh.GetCurrentFeedItems())-1 {
		fh.CurrentItemIndex = fh.CurrentItemIndex + 1
	}
}

func (fh *FeedParser) PrevItem() {
	if fh.CurrentItemIndex > 0 {
		fh.CurrentItemIndex = fh.CurrentItemIndex - 1
	}
}

func (fh *FeedParser) ResetFeedIdx() {
	fh.CurrentFeedIndex = 0
}

func (fh *FeedParser) ResetItemIdx() {
	fh.CurrentItemIndex = 0
}

func (fh *FeedParser) GetCurrentFeedItems() []*model.Item {
	if fh.getCurrentFeed() == nil {
		return []*model.Item{}
	}

	return fh.getCurrentFeed().Items

}

func (fh *FeedParser) GetCurrentFeedItemsName() []string {
	if fh.getCurrentFeed() == nil {
		return []string{}
	}

	return fh.getItemsNameFromFeed(fh.getCurrentFeed())

}

func (fh *FeedParser) GetCurrentFeedItemsNameAndStatus() []string {
	if fh.getCurrentFeed() == nil {
		return []string{}
	}

	return fh.getItemsNameAndStatusFromFeed(fh.getCurrentFeed())

}

func (fh *FeedParser) getItemsNameFromFeed(feed *model.Feed) []string {
	var titles []string
	for _, item := range feed.Items {
		titles = append(titles, item.Title)
	}

	return titles
}

func (fh *FeedParser) getItemsNameAndStatusFromFeed(feed *model.Feed) []string {
	var titles []string
	for _, item := range feed.Items {
		t := item.Title
		if item.LocalFileName != "" {
			t = "[D] " + t

		}
		titles = append(titles, t)
	}

	return titles
}

func (fh *FeedParser) GetCurrentItemDescription() string {
	ci := fh.getCurrentItem()
	if ci == nil {
		return ""
	}

	return ci.Description

}

func (fh *FeedParser) GetCurrentItemName() string {
	ci := fh.getCurrentItem()
	if ci == nil {
		return ""
	}

	return ci.Title
}

func (fh *FeedParser) GetCurrentItemUrl() string {
	ci := fh.getCurrentItem()
	if ci == nil {
		return ""
	}

	return ci.Url
}

func (fh *FeedParser) SetCurrentItemLocalFileName(fileName string) {
	ci := fh.getCurrentItem()
	ci.LocalFileName = fileName
}

func (fh *FeedParser) GetCurrentItemLocalFileName() string {
	ci := fh.getCurrentItem()
	if ci == nil {
		return ""
	}

	return ci.LocalFileName
}
