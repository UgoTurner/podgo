package service

import (
	"log"
	"os"
	"time"

	"github.com/ugo/beep"
	"github.com/ugo/beep/mp3"
	"github.com/ugo/beep/speaker"
)

type Player struct {
	Control *beep.Ctrl
}

func (p *Player) Play(path string) {
	f, err := os.Open(path)
	// Check for errors when opening the file
	if err != nil {
		log.Fatal(path)
	}

	// Decode the .mp3 File, if you have a .wav file, use wav.Decode(f)
	s, format, _ := mp3.Decode(f)
	p.Control = &beep.Ctrl{Streamer: s, Paused: false}

	// Init the Speaker with the SampleRate of the format and a buffer size of 1/10s
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Channel, which will signal the end of the playback.
	//playing := make(chan struct{})

	// Now we Play our Streamer on the Speaker
	speaker.Play(beep.Seq(p.Control, beep.Callback(func() {
		// Callback after the stream Ends
		//close(playing)
	})))
	//	<-playing

}

func (p *Player) TogglePlayPause() {
	speaker.Lock()
	p.Control.Paused = !p.Control.Paused
	speaker.Unlock()

}
