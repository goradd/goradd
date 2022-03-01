package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/goradd/gofile/pkg/sys"
)

// This file runs the tests found in the test directory. It is set up so that code coverage can be checked as well.
// Doing this is a little tricky, since got generates go code that then gets compiled and run again. Each part of the
// process may generate errors. We test the process from end to end, but to do code coverage, we must directly access
// the main file as part of the test.
func TestGoradd(t *testing.T) {
	go main()

	time.Sleep(time.Second * 5)

	var cmd string

	os := runtime.GOOS
	switch os {
	case "windows":
		cmd = "Chrome --headless --remote-debugging-port=9222 http://localhost:8000/goradd/Test.g?all=1"
	case "darwin":
		cmd = `"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --headless --remote-debugging-port=9222 http://localhost:8000/goradd/Test.g?all=1`
	case "linux":
		cmd = "google-chrome-stable --headless --remote-debugging-port=9222 http://localhost:8000/goradd/Test.g?all=1"

	if _, err := sys.ExecuteShellCommand(cmd); err != nil {
		if e, ok := err.(*exec.Error); ok {
			_, _ = fmt.Fprintln(os.Stderr, "error running browser test :"+e.Error())
			os.Exit(1)
		} else if err2, ok2 := err.(*exec.ExitError); ok2 {
			_, _ = fmt.Fprintln(os.Stderr, string(err2.Stderr))
			os.Exit(1)
		}
	}
}

