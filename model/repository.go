package model

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/kennygrant/sanitize"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/sirupsen/logrus"
	"github.com/ugoturner/podgo/conf"
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
		if err := fr.Db.Write(conf.FeedCollection, feedKey, feed); err != nil {
			fr.Logger.Errorf("Failed to update feed '%s': %v", feed.Title, err)
			return err
		}
	}
	return nil
}

// FetchAll returns every feed stored in the database.
func (fr *FeedRepository) FetchAll() ([]*Feed, error) {
	records, err := fr.Db.ReadAll(conf.FeedCollection)
	if err != nil {
		// No collection yet (e.g. first launch, no feeds added) is not an error.
		if os.IsNotExist(err) {
			return nil, nil
		}
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
