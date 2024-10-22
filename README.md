# Podgo

**Podgo** is a terminal-based app for managing and playing podcasts. It uses the [Gocui](https://github.com/jroimartin/gocui) library, wrapped in a custom-built layer called [Songocui](https://github.com/UgoTurner/songocui), and handles feed parsing with [Gofeed](https://github.com/mmcdole/gofeed).

![Layout](./screenshot.png)

## Installation

### From source code

```bash
cd podgo
go build .
./podgo
```

## Usage

### Layout
The terminal is divided into three panels:
1. **Podcasts list (side panel)**: Displays the user's saved podcasts.
2. **Tracks list (main panel)**: Shows the track list for the selected podcast.
3. **Summary (footer panel)**: Contains the player and shows prompts when actions are required.

### Key bindings

Key combination | Description
---|---
<kbd>&uarr; and &darr;</kbd>|Navigate up and down in the list views
<kbd>&rarr;</kbd>|Enter into the next view
<kbd>Ctrl</kbd>+<kbd>a</kbd>|Add new podcast feed (then <kbd>Enter</kbd> to confirm it)
<kbd>Ctrl</kbd>+<kbd>d</kbd>|Download the selected track
<kbd>Ctrl</kbd>+<kbd>p</kbd>|Play (and download if it is not done yet) the selected track
<kbd>Ctrl</kbd>+<kbd>space</kbd>|Toggle play/pause
<kbd>Ctrl</kbd>+<kbd>f</kbd>|Seek forward
<kbd>Ctrl</kbd>+<kbd>b</kbd>|Seek backward
<kbd>Ctrl</kbd>+<kbd>c</kbd>|Exit Podgo
