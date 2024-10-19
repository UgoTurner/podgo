package handler

import "testing"

func TestExtractFileName(t *testing.T) {
	cases := []struct {
		url  string
		want string
	}{
		{"https://cdn.example.com/podcast/ep-42.mp3", "ep-42.mp3"},
		{"https://example.com/ep.mp3?token=abc", "ep.mp3"},
		{"https://example.com/path/ep.mp3#fragment", "ep.mp3"},
		{"noslash.mp3", "noslash.mp3"},
		{"", ""},
	}
	for _, c := range cases {
		if got := extractFileName(c.url); got != c.want {
			t.Errorf("extractFileName(%q) = %q, want %q", c.url, got, c.want)
		}
	}
}
