package service

import (
	"testing"
	"time"
)

func TestDurationToString(t *testing.T) {
	p := &Player{}
	cases := []struct {
		in   time.Duration
		want string
	}{
		{0, "00:00:00"},
		{5 * time.Second, "00:00:05"},
		{90 * time.Second, "00:01:30"},
		{time.Hour + 2*time.Minute + 3*time.Second, "01:02:03"},
	}
	for _, c := range cases {
		if got := p.durationToString(c.in); got != c.want {
			t.Errorf("durationToString(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestTogglePlayPauseNoPanicWhenIdle(_ *testing.T) {
	p := &Player{Logger: silentLogger()}
	p.TogglePlayPause() // must be a no-op, not a nil panic
	p.Seek(10)          // must be a no-op, not a nil panic
}
