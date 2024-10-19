package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
)

func silentLogger() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(io.Discard)
	return l
}

func TestDownloadFileSuccess(t *testing.T) {
	const body = "fake-mp3-bytes"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	dst := filepath.Join(t.TempDir(), "track.mp3")
	var progressed, succeeded, failed bool

	err := DownloadFile(dst, srv.URL,
		func(string) { progressed = true },
		func() { succeeded = true },
		func() { failed = true },
		silentLogger(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !succeeded || failed {
		t.Fatalf("expected success callback, got success=%v failure=%v", succeeded, failed)
	}
	if !progressed {
		t.Error("expected progress callback to be called")
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("reading downloaded file: %v", err)
	}
	if string(got) != body {
		t.Errorf("expected %q, got %q", body, string(got))
	}
}

func TestDownloadFileBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	dst := filepath.Join(t.TempDir(), "track.mp3")
	var failed bool

	err := DownloadFile(dst, srv.URL,
		func(string) {},
		func() { t.Error("success callback must not be called on 404") },
		func() { failed = true },
		silentLogger(),
	)
	if err == nil {
		t.Fatal("expected an error on 404")
	}
	if !failed {
		t.Error("expected failure callback on 404")
	}
}

func TestWriteCounterAccumulates(t *testing.T) {
	var last string
	wc := &WriteCounter{onProgress: func(s string) { last = s }, Logger: silentLogger()}

	n, err := wc.Write(make([]byte, 500))
	if err != nil || n != 500 {
		t.Fatalf("write returned n=%d err=%v", n, err)
	}
	if _, err := wc.Write(make([]byte, 524)); err != nil {
		t.Fatalf("second write err: %v", err)
	}
	if wc.Total != 1024 {
		t.Fatalf("expected total 1024, got %d", wc.Total)
	}
	if last == "" {
		t.Error("expected progress string to be set")
	}
}
