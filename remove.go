package osutil

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/avast/retry-go"
)

// RemoveTemporaryDirectory removes the given temporary directory path from disk with special handling for Windows.
// The reason we need special handling is because Windows seems to be colossally stupid when it comes to handling the open-ess of files and directories. https://$INTERNAL/symflower/symflower/-/merge_requests/2399#note_293837.
func RemoveTemporaryDirectory(directoryPath string) {
	isRetryableError := func(err error) bool {
		return IsWindows() && (strings.Contains(err.Error(), "Access is denied") || strings.Contains(err.Error(), "The process cannot access the file because it is being used by another process."))
	}

	if err := retry.Do(
		func() error {
			return os.RemoveAll(directoryPath)
		},
		retry.Attempts(3),
		retry.Delay(2*time.Second),
		retry.LastErrorOnly(true),
		retry.RetryIf(func(err error) bool {
			// On Windows we sometimes receive an access denied error which happens because a process has exited but is not yet cleaned up by the operating system and still has a handler to a file we want to delete. In this case we wait a while and try again to remove the directory.
			return isRetryableError(err)
		}),
	); err != nil {
		if !isRetryableError(err) {
			if err != nil {
				panic(err)
			}
		}

		// At this point we have given our best to delete the directory. The only chance we now have is to end this processes we are currently in and then delete the directory. See https://$INTERNAL/symflower/symflower/-/merge_requests/2399#note_293837 for details. The only thing we can now do is to log the error, so the processes higher up in the process tree can know about this tragedy and deal with it accordingly.
		fmt.Fprintf(os.Stderr, "cannot remove temporary directory: %s\n", err)
	}
}
