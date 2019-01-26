package service

import (
	"github.com/mmcdole/gofeed"
)

type Item struct {
	Title         string
	Description   string
	Url           string
	LocalFileName string
	Playing       bool
}

type Feed struct {
	Title string
	Items []*Item
}

type FeedParser struct {
	Feeds            []*Feed
	CurrentFeedIndex int
	CurrentItemIndex int
}

func (fh *FeedParser) getFeedUrls() []string {
	return []string{
		"https://distorsionpodcast.com/podcasts.xml",
		"http://feeds.feedburner.com/bureaudesmysteres?fmt=xml",
		//"http://feeds.soundcloud.com/playlists/soundcloud:playlists:253185017/sounds.rss",
	}
}

func (fh *FeedParser) extractItems(gItems []*gofeed.Item) []*Item {
	var it []*Item
	for _, gItem := range gItems {
		it = append(
			it,
			&Item{
				Title:       gItem.Title,
				Description: gItem.Description,
				Url:         fh.extractItemLink(gItem),
			})
	}

	return it
}

func (fh *FeedParser) extractFeed(gFeed *gofeed.Feed) *Feed {
	return &Feed{
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

func (fh *FeedParser) LoadFeeds() {
	/* parse local file
	file, _ := os.Open("/path/to/a/file.xml")
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	fmt.Println(feed.Title)
	*/

	fp := gofeed.NewParser()
	for _, furl := range fh.getFeedUrls() {
		f, _ := fp.ParseURL(furl)
		fh.Feeds = append(fh.Feeds, fh.extractFeed(f))
	}
	if len(fh.Feeds) > 0 {
		fh.CurrentFeedIndex = 0
		fh.CurrentItemIndex = 0
	}
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

func (fh *FeedParser) getFeedByName(feedName string) *Feed {
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

func (fh *FeedParser) getFeedByIndex(i int) *Feed {
	// todo: test key exist
	return fh.Feeds[i]

}

func (fh *FeedParser) getItemByIndex(i int) *Item {
	// todo: test key exist
	//fmt.Printf("%+v\n", fh.getCurrentFeed().Items[i])
	return fh.getCurrentFeed().Items[i]

}

func (fh *FeedParser) getCurrentFeed() *Feed {
	// todo: test key exist
	return fh.getFeedByIndex(fh.CurrentFeedIndex)

}

func (fh *FeedParser) getCurrentItem() *Item {
	// todo: test key exist
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
	// todo: test key exist
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
	// todo: test key exist
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

func (fh *FeedParser) GetCurrentFeedItems() []*Item {
	if fh.getCurrentFeed() == nil {
		return []*Item{}
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

func (fh *FeedParser) getItemsNameFromFeed(feed *Feed) []string {
	var titles []string
	for _, item := range feed.Items {
		titles = append(titles, item.Title)
	}

	return titles
}

func (fh *FeedParser) getItemsNameAndStatusFromFeed(feed *Feed) []string {
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
