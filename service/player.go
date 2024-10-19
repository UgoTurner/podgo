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

// Player decodes and plays MP3 tracks through the ALSA speaker, exposing
// play/pause and seek controls.
type Player struct {
	Control          *beep.Ctrl
	Stream           beep.StreamSeekCloser
	SampleRate       beep.SampleRate
	Playing          bool
	PlayingTrackName string
	Logger           *logrus.Logger

	speakerInit bool
	speakerRate beep.SampleRate
}

type playerProgress func(string)

func (p *Player) durationToString(duration time.Duration) string {
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// handleProgression reports playback position once per second until the track
// finishes or is replaced by another one.
func (p *Player) handleProgression(progress playerProgress) {
	stream := p.Stream
	for p.Playing && p.Control != nil && p.Stream == stream {
		totalDuration := p.durationToString(p.SampleRate.D(stream.Len()))
		current := p.durationToString(p.SampleRate.D(stream.Position()))
		progress(current + " / " + totalDuration)
		time.Sleep(1 * time.Second)
	}
}

// Play decodes the MP3 at path and plays it, reporting playback position
// through the progress callback until the track ends or is replaced.
func (p *Player) Play(path, trackName string, progress playerProgress) error {
	// Stop any track currently playing so we don't overlap streams.
	if p.speakerInit {
		speaker.Clear()
	}
	p.Playing = false
	p.Control = nil

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
		return fmt.Errorf("failed to decode MP3 file '%s': %v", path, err)
	}
	defer func() {
		if closeErr := s.Close(); closeErr != nil {
			p.Logger.WithFields(logrus.Fields{"trackName": trackName, "error": closeErr}).Warn("Failed to close stream for track")
		}
	}()

	p.Stream = s
	p.SampleRate = format.SampleRate
	p.Control = &beep.Ctrl{Streamer: p.Stream, Paused: false}

	// The speaker can only be initialized once per sample rate; re-init only
	// when it changes to avoid resetting the audio device on every track.
	if !p.speakerInit || p.speakerRate != format.SampleRate {
		if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
			p.Logger.WithError(err).Error("Failed to initialize speaker")
			return fmt.Errorf("failed to initialize speaker: %w", err)
		}
		p.speakerInit = true
		p.speakerRate = format.SampleRate
	}

	p.Playing = true
	speaker.Play(beep.Seq(p.Control, beep.Callback(func() {
		p.Playing = false
	})))
	p.handleProgression(progress)

	return nil
}

// TogglePlayPause pauses or resumes the current track.
func (p *Player) TogglePlayPause() {
	if p.Control == nil {
		return
	}
	speaker.Lock()
	p.Control.Paused = !p.Control.Paused
	speaker.Unlock()
	p.Logger.Infof("Track '%s' playback paused: %v", p.PlayingTrackName, p.Control.Paused)
}

// Seek moves the playback position by sec seconds (negative to rewind).
func (p *Player) Seek(sec int) {
	if p.Stream == nil {
		return
	}
	speaker.Lock()
	target := p.Stream.Position() + sec*int(p.SampleRate)
	if target < 0 {
		target = 0
	}
	if target >= p.Stream.Len() {
		target = p.Stream.Len() - 1
	}
	if err := p.Stream.Seek(target); err != nil {
		p.Logger.WithError(err).Warnf("Failed to seek track '%s'", p.PlayingTrackName)
	}
	speaker.Unlock()
	p.Logger.Infof("Track '%s' seeked by %d seconds", p.PlayingTrackName, sec)
}
