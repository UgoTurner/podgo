package model

import (
	"encoding/json"
	"fmt"
	"strings"

	scribble "github.com/nanobox-io/golang-scribble"

	"github.com/kennygrant/sanitize"
)

type FeedRepository struct {
	Db *scribble.Driver
}

func (fr *FeedRepository) Update(feeds []*Feed) {
	for _, feed := range feeds {
		// log.WithFields(log.Fields{
		// 	"feed.Url":   feed.Url,
		// 	"feed.Title": feed.Title,
		// }).Info("Insert new feed")
		fr.Db.Write("feed", strings.ToLower(sanitize.BaseName(feed.Title)), &feed)
	}
}

func (fr *FeedRepository) FetchAll() []*Feed {
	records, err := fr.Db.ReadAll("feed")
	if err != nil {
		fmt.Println("Error", err)
	}

	feeds := []*Feed{}
	for _, f := range records {
		feedFound := Feed{}
		if err := json.Unmarshal([]byte(f), &feedFound); err != nil {
			fmt.Println("Error", err)
		}
		feeds = append(feeds, &feedFound)
	}

	return feeds
}
