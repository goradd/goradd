package travis

import (
	"path/filepath"
	"runtime"
)

// This code is here simply to help the goradd tool find the test directory.
// In a modules aware world, at the very start of using goradd, you have nothing. Possibly multiple versions of
// goradd are installed in the GOPATH directory, and there is no initial frame of reference for how to get to
// the goradd-project directory. Soooo, we need to read the location that is included in the compiled version
// of the application.

var TestFolderLocation string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	TestFolderLocation = filepath.Dir(filename)
}
