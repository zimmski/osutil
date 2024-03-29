package osutil

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

// HTTPClient defines an HTTP client with sane default settings.
var HTTPClient *http.Client = func() *http.Client {
	c := &http.Client{
		// The timeout defaults of the default client are terrible. See https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/ for details.
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 0, // This timeout includes the whole process of downloading a file. Hence, big files always run into a timeout so we are setting the timeout granularly.
	}
	c.Jar, _ = cookiejar.New(nil)

	return c
}()

// DownloadFile downloads a file from the URL to the file path.
func DownloadFile(url string, filePath string) (err error) {
	resp, err := HTTPClient.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			err = errors.Join(err, e)
		}
	}()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = errors.Join(err, e)
		}
	}()

	_, err = io.Copy(out, resp.Body)

	return err
}

// DownloadFileWithProgress downloads a file from the URL to the file path while printing a progress to STDOUT.
func DownloadFileWithProgress(url string, filePath string) (err error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	response, err := HTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		if e := response.Body.Close(); e != nil {
			if err != nil {
				err = errors.Join(err, e)
			} else {
				err = e
			}
		}
	}()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("downloading file failed with status code %d: %s", response.StatusCode, response.Status)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if e := file.Close(); e != nil {
			if err != nil {
				err = errors.Join(err, e)
			} else {
				err = e
			}
		}
	}()

	pg := ProgressBarBytes(os.Stdout, int(response.ContentLength), "downloading")
	defer func() {
		if e := pg.Close(); e != nil {
			if err != nil {
				err = errors.Join(err, e)
			} else {
				err = e
			}
		}
	}()

	if _, err := io.Copy(io.MultiWriter(file, pg), response.Body); err != nil {
		return err
	}

	return nil
}
