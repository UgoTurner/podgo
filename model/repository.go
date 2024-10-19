package model

import (
	"encoding/json"
	"strings"

	"github.com/kennygrant/sanitize"
	"github.com/nanobox-io/golang-scribble"
	"github.com/sirupsen/logrus"
)

// FeedRepository handles CRUD operations for feeds using a Scribble DB.
type FeedRepository struct {
	Db     *scribble.Driver
	Logger *logrus.Logger
}

// Update inserts or updates feeds in the Scribble DB.
func (fr *FeedRepository) Update(feeds []*Feed) error {
	for _, feed := range feeds {
		feedKey := strings.ToLower(sanitize.BaseName(feed.Title))
		if err := fr.Db.Write("feed", feedKey, feed); err != nil {
			fr.Logger.Errorf("Failed to update feed '%s': %v", feed.Title, err)
			return err
		}
	}
	return nil
}

func (fr *FeedRepository) FetchAll() ([]*Feed, error) {
	records, err := fr.Db.ReadAll("feed")
	if err != nil {
		fr.Logger.Errorf("Failed to fetch feeds: %v", err)
		return nil, err
	}

	var feeds []*Feed
	for _, record := range records {
		if !json.Valid([]byte(record)) {
			fr.Logger.Warn("Skipping invalid JSON record")
			continue
		}

		var feed Feed
		if err := json.Unmarshal([]byte(record), &feed); err != nil {
			fr.Logger.Errorf("Failed to parse feed data: %v", err)
			return nil, err
		}
		feeds = append(feeds, &feed)
	}

	return feeds, nil
}
