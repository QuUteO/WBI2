package downloader

import (
	"fmt"
	"net/http"
	"time"
)

type Downloader struct {
	client *http.Client
}

func NewDownloader(timeout time.Duration) *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (d *Downloader) Download(targetURL string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "mini-wget/1.0")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status %s for %s", resp.Status, targetURL)
	}

	return resp, nil
}
