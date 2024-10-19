// Package model defines the domain types and the Scribble-backed repository.
package model

// Feed is a podcast feed and its episodes.
type Feed struct {
	Title string
	URL   string `json:"Url"`
	Items []*Item
}

// Item is a single podcast episode.
type Item struct {
	Title         string
	Description   string
	URL           string `json:"Url"`
	LocalFileName string
	Playing       bool
}
