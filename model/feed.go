package model

type Feed struct {
	Title, Url string
	Items      []*Item
}

type Item struct {
	Title         string
	Description   string
	Url           string
	LocalFileName string
	Playing       bool
}
