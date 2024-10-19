// Package conf holds the app's constants, runtime paths and the embedded
// panel/keybinding configuration.
package conf

// View names and titles used across the UI.
const (
	SideViewName  = "side"
	SideViewTitle = "Podcasts"

	MainViewName  = "main-primary"
	MainViewTitle = "Tracks"

	MainDetailsViewName  = "main-secondary"
	MainDetailsViewTitle = ""

	FooterViewName  = "footer"
	FooterViewTitle = ""

	PromptViewName  = "prompt"
	PromptViewTitle = ""

	FlashMessageViewName  = "flashmsg"
	FlashMessageViewTitle = ""
)

// FeedCollection is the Scribble collection name used to store feeds.
const FeedCollection = "feed"
