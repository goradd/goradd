package sys

import (
	"fmt"
	"github.com/goradd/gofile/pkg/sys"
	"runtime"
)

// LaunchDefaultBrowser launches the system's default browser and opens the browser to the given url.
// Not that most Linux systems do not have a default.
func LaunchDefaultBrowser(url string) (err error) {
	switch runtime.GOOS {
	case `darwin`:
		_, err = sys.ExecuteShellCommand(fmt.Sprintf("open %s", url))
	case `windows`:
		_, err = sys.ExecuteShellCommand(fmt.Sprintf("start %s", url))
	}

	// Ubuntu has a way to get the preferred browser. Not sure about other flavors of Unix, but we don't really have
	// a way of detecting a Unix flavor
	return
}

// LaunchChromeHeadlessBrowser will launch google chrome with the given url.
func LaunchChrome(url string) (err error) {
	go func() {
		switch runtime.GOOS {
		case `darwin`:
			_, err = sys.ExecuteShellCommand(fmt.Sprintf(`"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" %s`, url))
		case `windows`:
			_, err = sys.ExecuteShellCommand(fmt.Sprintf("start chrome %s", url))
		case `linux`:
			_, err = sys.ExecuteShellCommand(fmt.Sprintf("google-chrome-stable %s &", url))
		}
	}()

	return
}
