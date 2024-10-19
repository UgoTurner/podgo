package service

import (
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/ugoturner/podgo/model"
)

func newParserWithFeeds() *FeedParser {
	return &FeedParser{
		Feeds: []*model.Feed{
			{Title: "Feed A", URL: "http://a", Items: []*model.Item{
				{Title: "A1", URL: "http://a/1"},
				{Title: "A2", URL: "http://a/2", LocalFileName: "a2.mp3"},
			}},
			{Title: "Feed B", URL: "http://b", Items: []*model.Item{
				{Title: "B1", URL: "http://b/1"},
			}},
		},
	}
}

func TestFeedNavigationStaysInBounds(t *testing.T) {
	fp := newParserWithFeeds()

	fp.PrevFeed() // already at 0, must stay
	if fp.CurrentFeedIndex != 0 {
		t.Fatalf("expected feed index 0, got %d", fp.CurrentFeedIndex)
	}

	fp.NextFeed()
	fp.NextFeed() // only two feeds, must clamp at 1
	if fp.CurrentFeedIndex != 1 {
		t.Fatalf("expected feed index 1, got %d", fp.CurrentFeedIndex)
	}
}

func TestItemNavigationStaysInBounds(t *testing.T) {
	fp := newParserWithFeeds()

	fp.PrevItem()
	if fp.CurrentItemIndex != 0 {
		t.Fatalf("expected item index 0, got %d", fp.CurrentItemIndex)
	}

	fp.NextItem()
	fp.NextItem() // feed A has two items, must clamp at 1
	if fp.CurrentItemIndex != 1 {
		t.Fatalf("expected item index 1, got %d", fp.CurrentItemIndex)
	}
}

func TestGetCurrentFeedItemsNameAndStatus(t *testing.T) {
	fp := newParserWithFeeds()
	got := fp.GetCurrentFeedItemsNameAndStatus()
	want := []string{"A1", "[D] A2"}
	if len(got) != len(want) {
		t.Fatalf("expected %d titles, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("title %d: expected %q, got %q", i, want[i], got[i])
		}
	}
}

func TestGetFeedNames(t *testing.T) {
	fp := newParserWithFeeds()
	got := fp.GetFeedNames()
	if len(got) != 2 || got[0] != "Feed A" || got[1] != "Feed B" {
		t.Fatalf("unexpected feed names: %v", got)
	}
}

func TestGettersOnEmptyParser(t *testing.T) {
	fp := &FeedParser{}
	if fp.GetCurrentFeedName() != "" {
		t.Error("expected empty feed name on empty parser")
	}
	if fp.GetCurrentItemName() != "" {
		t.Error("expected empty item name on empty parser")
	}
	if fp.GetCurrentItemURL() != "" {
		t.Error("expected empty item url on empty parser")
	}
	if len(fp.GetCurrentFeedItems()) != 0 {
		t.Error("expected no items on empty parser")
	}
}

func TestExtractItemLinkPrefersAudioMpeg(t *testing.T) {
	fp := &FeedParser{}
	item := &gofeed.Item{Enclosures: []*gofeed.Enclosure{
		{URL: "http://x/cover.jpg", Type: "image/jpeg"},
		{URL: "http://x/ep.mp3", Type: "audio/mpeg"},
	}}
	if got := fp.extractItemLink(item); got != "http://x/ep.mp3" {
		t.Fatalf("expected mp3 url, got %q", got)
	}
	if got := fp.extractItemLink(nil); got != "" {
		t.Fatalf("expected empty url for nil item, got %q", got)
	}
}
