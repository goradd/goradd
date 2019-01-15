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
