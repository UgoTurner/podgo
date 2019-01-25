package service

import (
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Player struct {
	Control    *beep.Ctrl
	Stream     beep.StreamSeekCloser
	SampleRate beep.SampleRate
	Playing    bool
}

type playerProgress func(string)

func (p *Player) durationToString(duration time.Duration) string {

	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	return strconv.Itoa(int(hours)) + ":" +
		strconv.Itoa(int(minutes)) + ":" +
		strconv.Itoa(int(seconds))

}

func (p *Player) handleProgression(progress playerProgress) {
	if !p.Playing {
		return
	}
	totalDuration := p.durationToString(p.SampleRate.D(p.Stream.Len()))
	current := p.durationToString(p.SampleRate.D(p.Stream.Position()))
	progress(current + " / " + totalDuration)
	time.Sleep(1 * time.Second)
	go p.handleProgression(progress)

}
func (p *Player) Play(path string, progress playerProgress) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(path)
	}
	s, format, _ := mp3.Decode(f)
	p.Stream = s
	p.SampleRate = format.SampleRate
	p.Control = &beep.Ctrl{Streamer: p.Stream, Paused: false}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	p.Playing = true
	speaker.Play(beep.Seq(p.Control, beep.Callback(func() {
		p.Playing = false
	})))
	p.handleProgression(progress)
}

func (p *Player) TogglePlayPause() {
	speaker.Lock()
	p.Control.Paused = !p.Control.Paused
	speaker.Unlock()
}

func (p *Player) Seek(sec int) {
	speaker.Lock()
	p.Stream.Seek(p.Stream.Position() + sec*int(p.SampleRate))
	speaker.Unlock()
}
