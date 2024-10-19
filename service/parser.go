package service

import (
	"github.com/mmcdole/gofeed"
	"github.com/ugoturner/podgo/model"

	"github.com/sirupsen/logrus"
)

// FeedParser manages the feeds and their items.
type FeedParser struct {
	Feeds            []*model.Feed
	CurrentFeedIndex int
	CurrentItemIndex int
	Logger         *logrus.Logger
}

// GetFeeds returns all feeds.
func (fp *FeedParser) GetFeeds() []*model.Feed {
	return fp.Feeds
}

// SetFeeds sets the feeds.
func (fp *FeedParser) SetFeeds(feeds []*model.Feed) {
	fp.Feeds = feeds
}

// AddFeed adds a new feed to the list.
func (fp *FeedParser) AddFeed(feed *model.Feed) {
	fp.Feeds = append(fp.Feeds, feed)
}

// extractItems extracts the items from gofeed.Item and converts them to model.Item.
func (fp *FeedParser) extractItems(gItems []*gofeed.Item) []*model.Item {
	var items []*model.Item
	for _, gItem := range gItems {
		items = append(items, &model.Item{
			Title:       gItem.Title,
			Description: gItem.Description,
			Url:         fp.extractItemLink(gItem),
		})
	}
	return items
}

// extractFeed converts a gofeed.Feed into a model.Feed.
func (fp *FeedParser) extractFeed(gFeed *gofeed.Feed) *model.Feed {
	return &model.Feed{
		Title: gFeed.Title,
		Items: fp.extractItems(gFeed.Items),
	}
}

// extractItemLink returns the audio URL from the gofeed.Item.
func (fp *FeedParser) extractItemLink(gItem *gofeed.Item) string {
	if gItem == nil {
		return ""
	}
	for _, en := range gItem.Enclosures {
		if en.Type == "audio/mpeg" {
			return en.URL
		}
	}
	return ""
}

// LoadFeedFromUrl loads and parses a feed from a URL.
func (fp *FeedParser) LoadFeedFromUrl(url string) *model.Feed {
	feedParser := gofeed.NewParser()
	gFeed, err := feedParser.ParseURL(url)
	if err != nil {
		fp.Logger.Warnf("Unable to parse URL: '%s': %w", url, err)
		return nil
	}
	feed := fp.extractFeed(gFeed)
	feed.Url = url
	return feed
}

// GetFeedNames returns the titles of all feeds.
func (fp *FeedParser) GetFeedNames() []string {
	var titles []string
	for _, feed := range fp.Feeds {
		titles = append(titles, feed.Title)
	}
	return titles
}

// GetItemNames returns the titles of all items across all feeds.
func (fp *FeedParser) GetItemNames() []string {
	var titles []string
	for _, feed := range fp.Feeds {
		for _, item := range feed.Items {
			titles = append(titles, item.Title)
		}
	}
	return titles
}

// getFeedByName returns a feed matching the given name.
func (fp *FeedParser) getFeedByName(feedName string) *model.Feed {
	for _, feed := range fp.Feeds {
		if feed.Title == feedName {
			return feed
		}
	}
	return nil
}

// GetItemNamesByFeedName returns the item titles for a specific feed by its name.
func (fp *FeedParser) GetItemNamesByFeedName(feedName string) []string {
	feed := fp.getFeedByName(feedName)
	if feed == nil {
		return nil
	}
	var titles []string
	for _, item := range feed.Items {
		titles = append(titles, item.Title)
	}
	return titles
}

// getFeedByIndex returns a feed by its index.
func (fp *FeedParser) getFeedByIndex(index int) *model.Feed {
	if index >= 0 && index < len(fp.Feeds) {
		return fp.Feeds[index]
	}
	return nil
}

// getItemByIndex returns an item by its index.
func (fp *FeedParser) getItemByIndex(index int) *model.Item {
	currentFeed := fp.getCurrentFeed()
	if currentFeed != nil && index >= 0 && index < len(currentFeed.Items) {
		return currentFeed.Items[index]
	}
	return nil
}

// getCurrentFeed returns the current feed.
func (fp *FeedParser) getCurrentFeed() *model.Feed {
	return fp.getFeedByIndex(fp.CurrentFeedIndex)
}

// getCurrentItem returns the current item.
func (fp *FeedParser) getCurrentItem() *model.Item {
	return fp.getItemByIndex(fp.CurrentItemIndex)
}

// GetCurrentFeedName returns the title of the current feed.
func (fp *FeedParser) GetCurrentFeedName() string {
	currentFeed := fp.getCurrentFeed()
	if currentFeed == nil {
		return ""
	}
	return currentFeed.Title
}

// NextFeed moves to the next feed.
func (fp *FeedParser) NextFeed() {
	if fp.CurrentFeedIndex < len(fp.Feeds)-1 {
		fp.CurrentFeedIndex++
	}
}

// PrevFeed moves to the previous feed.
func (fp *FeedParser) PrevFeed() {
	if fp.CurrentFeedIndex > 0 {
		fp.CurrentFeedIndex--
	}
}

// NextItem moves to the next item.
func (fp *FeedParser) NextItem() {
	if fp.CurrentItemIndex < len(fp.GetCurrentFeedItems())-1 {
		fp.CurrentItemIndex++
	}
}

// PrevItem moves to the previous item.
func (fp *FeedParser) PrevItem() {
	if fp.CurrentItemIndex > 0 {
		fp.CurrentItemIndex--
	}
}

// ResetFeedIdx resets the current feed index to zero.
func (fp *FeedParser) ResetFeedIdx() {
	fp.CurrentFeedIndex = 0
}

// ResetItemIdx resets the current item index to zero.
func (fp *FeedParser) ResetItemIdx() {
	fp.CurrentItemIndex = 0
}

// GetCurrentFeedItems returns the items of the current feed.
func (fp *FeedParser) GetCurrentFeedItems() []*model.Item {
	currentFeed := fp.getCurrentFeed()
	if currentFeed == nil {
		return []*model.Item{}
	}
	return currentFeed.Items
}

// GetCurrentFeedItemsName returns the item titles of the current feed.
func (fp *FeedParser) GetCurrentFeedItemsName() []string {
	currentFeed := fp.getCurrentFeed()
	if currentFeed == nil {
		return []string{}
	}
	return fp.getItemsNameFromFeed(currentFeed)
}

// GetCurrentFeedItemsNameAndStatus returns the item titles and their statuses (downloaded or not).
func (fp *FeedParser) GetCurrentFeedItemsNameAndStatus() []string {
	currentFeed := fp.getCurrentFeed()
	if currentFeed == nil {
		return []string{}
	}
	return fp.getItemsNameAndStatusFromFeed(currentFeed)
}

// getItemsNameFromFeed returns the item titles of a specific feed.
func (fp *FeedParser) getItemsNameFromFeed(feed *model.Feed) []string {
	var titles []string
	for _, item := range feed.Items {
		titles = append(titles, item.Title)
	}
	return titles
}

// getItemsNameAndStatusFromFeed returns the item titles and their statuses (downloaded or not).
func (fp *FeedParser) getItemsNameAndStatusFromFeed(feed *model.Feed) []string {
	var titles []string
	for _, item := range feed.Items {
		title := item.Title
		if item.LocalFileName != "" {
			title = "[D] " + title
		}
		titles = append(titles, title)
	}
	return titles
}

// GetCurrentItemDescription returns the description of the current item.
func (fp *FeedParser) GetCurrentItemDescription() string {
	currentItem := fp.getCurrentItem()
	if currentItem == nil {
		return ""
	}
	return currentItem.Description
}

// GetCurrentItemName returns the title of the current item.
func (fp *FeedParser) GetCurrentItemName() string {
	currentItem := fp.getCurrentItem()
	if currentItem == nil {
		return ""
	}
	return currentItem.Title
}

// GetCurrentItemUrl returns the URL of the current item.
func (fp *FeedParser) GetCurrentItemUrl() string {
	currentItem := fp.getCurrentItem()
	if currentItem == nil {
		return ""
	}
	return currentItem.Url
}

// SetCurrentItemLocalFileName sets the local file name for the current item.
func (fp *FeedParser) SetCurrentItemLocalFileName(fileName string) {
	currentItem := fp.getCurrentItem()
	if currentItem != nil {
		currentItem.LocalFileName = fileName
	}
}

// GetCurrentItemLocalFileName returns the local file name of the current item.
func (fp *FeedParser) GetCurrentItemLocalFileName() string {
	currentItem := fp.getCurrentItem()
	if currentItem == nil {
		return ""
	}
	return currentItem.LocalFileName
}
