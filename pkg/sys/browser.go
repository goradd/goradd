package sys

import (
	"fmt"
	"github.com/goradd/gofile/pkg/sys"
	"runtime"
)

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
// One nice feature of google chrome is that you can launch it, give it a URL, and then the browser will listen
// for the URL and load it once the server on the other end becomes active.
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

