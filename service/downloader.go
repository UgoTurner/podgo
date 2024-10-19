package service

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
)

// WriteCounter tracks download progress and reports it.
type WriteCounter struct {
	Total      uint64
	onProgress func(string)
	Logger     *logrus.Logger
}

// Write updates the total number of bytes written and triggers the progress callback.
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.onProgress(humanize.Bytes(wc.Total))
	return n, nil
}

// Callbacks for handling progress, success, and failure events.
type ProgressCallback func(string)
type SuccessCallback func()
type FailureCallback func()

// DownloadFile downloads a file from the provided URL, writing it to the specified filepath.
// It accepts callbacks for progress updates, success, and failure handling.
func DownloadFile(filepath string, url string, progress ProgressCallback, success SuccessCallback, failure FailureCallback, logger *logrus.Logger) error {
	// Create the output file
	out, err := os.Create(filepath)
	if err != nil {
		logger.Errorf("Failed to create file '%s': %v", filepath, err)     // Log the error
		failure()                                                          // Call the failure callback
		return fmt.Errorf("failed to create file '%s': %w", filepath, err) // Wrap and return the error
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil {
			logger.Warnf("Failed to close file '%s': %v", filepath, closeErr) // Log any close error
		}
	}()

	// Perform the GET request
	resp, err := http.Get(url)
	if err != nil {
		logger.Errorf("Failed to download from URL '%s': %v", url, err)     // Log the error
		failure()                                                           // Call the failure callback
		return fmt.Errorf("failed to download from URL '%s': %w", url, err) // Wrap and return the error
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.Warnf("Failed to close response body for URL '%s': %v", url, closeErr) // Log any close error
		}
	}()

	// Ensure a valid response code
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("Bad status for URL '%s': %s", url, resp.Status) // Log the bad status
		failure()                                                      // Call the failure callback
		return fmt.Errorf("bad status: %s", resp.Status)               // Return an error
	}

	// Track the download progress
	counter := &WriteCounter{onProgress: progress, Logger: logger} // Pass logger to WriteCounter

	// Copy the response body to the file and track progress
	if _, err := io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		logger.Errorf("Failed to copy content to '%s': %v", filepath, err)     // Log the error
		failure()                                                              // Call the failure callback
		return fmt.Errorf("failed to copy content to '%s': %w", filepath, err) // Wrap and return the error
	}

	// Call success callback if download completes
	success()
	return nil
}
