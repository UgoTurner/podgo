package service

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ugo/beep"
	"github.com/ugo/beep/mp3"
	"github.com/ugo/beep/speaker"
)

type Player struct {
	Control *beep.Ctrl
	Stream  beep.StreamSeekCloser
	Playing bool
}

type playerProgress func(string)

func (p *Player) handleProgression(progress playerProgress) {
	if !p.Playing {
		return
	}
	progress("Durée : " + strconv.Itoa(p.Stream.Len()) + " -- Pos : " + strconv.Itoa(p.Stream.Position()))
	time.Sleep(1 * time.Second)
	go p.handleProgression(progress)

}
func (p *Player) Play(path string, progress playerProgress) {
	f, err := os.Open(path)
	// Check for errors when opening the file
	if err != nil {
		log.Fatal(path)
	}

	// Decode the .mp3 File, if you have a .wav file, use wav.Decode(f)
	s, format, _ := mp3.Decode(f)
	p.Stream = s
	p.Control = &beep.Ctrl{Streamer: s, Paused: false}

	// Init the Speaker with the SampleRate of the format and a buffer size of 1/10s
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Channel, which will signal the end of the playback.
	//playing := make(chan struct{})
	p.Playing = true
	// Now we Play our Streamer on the Speaker
	speaker.Play(beep.Seq(p.Control, beep.Callback(func() {
		// Callback after the stream Ends
		//close(playing)
		p.Playing = false
	})))
	p.handleProgression(progress)

	//	<-playing
}

func (p *Player) TogglePlayPause() {
	speaker.Lock()
	p.Control.Paused = !p.Control.Paused
	speaker.Unlock()
}

func (p *Player) Seek(sec int) {
	speaker.Lock()
	p.Stream.Seek(1)
	speaker.Unlock()
}
