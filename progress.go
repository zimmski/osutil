package osutil

import (
	"fmt"
	"io"
	"time"

	"github.com/schollz/progressbar/v3"
)

var defaultProgressBarOptions = []progressbar.Option{
	progressbar.OptionFullWidth(),
	progressbar.OptionSetRenderBlankState(false), // Print the progress bar when the first data arrives, not on initialization. This allows the user to print log entries until that first data, i.e. only one progress bar show until then.
	progressbar.OptionSetTheme(progressbar.Theme{
		Saucer:        "=",
		SaucerHead:    ">",
		SaucerPadding: " ",
		BarStart:      "[",
		BarEnd:        "]",
	}),
	progressbar.OptionSetWidth(80),
	progressbar.OptionShowCount(),
	progressbar.OptionSpinnerType(14),
	progressbar.OptionThrottle(65 * time.Millisecond),
}

// ProgressBar returns a progress bar for counting items with sane defaults that prints its updates to the given writer.
func ProgressBar(stream io.Writer, max int, description ...string) (progress *progressbar.ProgressBar) {
	var d string
	if len(description) > 0 {
		d = description[0]
	}

	os := []progressbar.Option{
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(stream, "\n")
		}),
		progressbar.OptionSetDescription(d),
		progressbar.OptionShowIts(),
		progressbar.OptionSetWriter(stream),
	}
	os = append(os, defaultProgressBarOptions...)

	return progressbar.NewOptions(max, os...)
}

// ProgressBarBytes returns a progress bar for counting bytes with sane defaults that prints its updates to the given writer.
func ProgressBarBytes(stream io.Writer, length int, description ...string) (progress *progressbar.ProgressBar) {
	var d string
	if len(description) > 0 {
		d = description[0]
	}

	os := []progressbar.Option{
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(stream, "\n")
		}),
		progressbar.OptionSetDescription(d),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWriter(stream),
	}
	os = append(os, defaultProgressBarOptions...)

	return progressbar.NewOptions(length, os...)
}

// ActivityIndicator prints a spinning activity indicator to the given stream until the indicator is stopped.
func ActivityIndicator(stream io.Writer, description ...string) (stopIndicator func()) {
	var d string
	if len(description) > 0 {
		d = description[0]
	}

	os := []progressbar.Option{
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(stream, "\n")
		}),
		progressbar.OptionSetDescription(d),
		progressbar.OptionSetWriter(stream),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionThrottle(65 * time.Millisecond),
	}

	const unknownLength = -1
	indicator := progressbar.NewOptions(unknownLength, os...)
	runChannel := make(chan struct{})
	go func() {
		for {
			select {
			case <-runChannel:
				return
			default:
				indicator.Add(1)
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	return func() {
		close(runChannel)
		_ = indicator.Finish()
	}
}
