package service

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/sirupsen/logrus"
)

type Player struct {
	Control          *beep.Ctrl
	Stream           beep.StreamSeekCloser
	SampleRate       beep.SampleRate
	Playing          bool
	PlayingTrackName string
	Logger           *logrus.Logger
}

type playerProgress func(string)

func (p *Player) durationToString(duration time.Duration) string {
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func (p *Player) handleProgression(progress playerProgress) {
	if p.Control == nil {
		return
	}
	totalDuration := p.durationToString(p.SampleRate.D(p.Stream.Len()))
	current := p.durationToString(p.SampleRate.D(p.Stream.Position()))
	progress(current + " / " + totalDuration)
	time.Sleep(1 * time.Second)
	p.handleProgression(progress)
}

func (p *Player) Play(path, trackName string, progress playerProgress) error {
	if p.Playing {
		p.Control = nil
	}

	// Attempt to open the file and log if there is an error
	f, err := os.Open(path)
	if err != nil {
		p.Logger.WithFields(logrus.Fields{"path": path, "error": err}).Error("Failed to open file")
		return fmt.Errorf("failed to open file '%s': %w", path, err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			p.Logger.WithFields(logrus.Fields{"path": path, "error": closeErr}).Warn("Failed to close file")
		}
	}()

	p.PlayingTrackName = trackName

	// Attempt to decode the MP3 file
	s, format, err := mp3.Decode(f)
	if err != nil {
		p.Logger.WithFields(logrus.Fields{"path": path, "error": err}).Error("Failed to decode MP3 file")
		return fmt.Errorf("failed to decode MP3 file '%s': %w", path, err)
	}
	defer func() {
		if closeErr := s.Close(); closeErr != nil {
			p.Logger.WithFields(logrus.Fields{"trackName": trackName, "error": closeErr}).Warn("Failed to close stream for track")
		}
	}()

	p.Stream = s
	p.SampleRate = format.SampleRate
	p.Control = &beep.Ctrl{Streamer: p.Stream, Paused: false}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	p.Playing = true
	speaker.Play(beep.Seq(p.Control, beep.Callback(func() {
		p.Playing = false
	})))
	p.handleProgression(progress)

	return nil
}

func (p *Player) TogglePlayPause() {
	speaker.Lock()
	p.Control.Paused = !p.Control.Paused
	speaker.Unlock()
	p.Logger.Infof("Track '%s' playback paused: %v", p.PlayingTrackName, p.Control.Paused)
}

func (p *Player) Seek(sec int) {
	speaker.Lock()
	p.Stream.Seek(p.Stream.Position() + sec*int(p.SampleRate))
	speaker.Unlock()
	p.Logger.Infof("Track '%s' seeked by %d seconds", p.PlayingTrackName, sec)
}
