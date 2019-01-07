package service

import (
	"io"
	"net/http"
	"os"

	humanize "github.com/dustin/go-humanize"
)

type WriteCounter struct {
	Total      uint64
	onProgress func(string)
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.onProgress(humanize.Bytes(wc.Total))
	return n, nil
}

type onProgress func(string)
type onSuccess func()
type onFailure func()

func DownloadFile(filepath string, url string, progress onProgress, success onSuccess, fail onFailure) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	counter := &WriteCounter{onProgress: progress}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		fail()
		return err
	}
	success()

	return nil
}
