package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

// This file runs the tests found in the test directory. It is set up so that code coverage can be checked as well.
// Doing this is a little tricky, since got generates go code that then gets compiled and run again. Each part of the
// process may generate errors. We test the process from end to end, but to do code coverage, we must directly access
// the main file as part of the test.
func TestGoradd(t *testing.T) {
	done := make(chan bool)
	go testMain(done)

	time.Sleep(time.Second * 3)

	var appName string

	currentOs := runtime.GOOS
	switch currentOs {
	case "windows":
		// old versions had it in the x86 directory
		appName = `C:\Program Files\Google\Chrome\Application\chrome.exe`
	case "darwin":
		appName = `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`
	case "linux":
		appName = "google-chrome-stable"
	}

	cmd := exec.Command(appName, "--headless", "--remote-debugging-port=9222", "http://localhost:8000/goradd/Test.g?all=1")
	err := cmd.Start()

	if err != nil {
		if e, ok := err.(*exec.Error); ok {
			_, _ = fmt.Fprintln(os.Stderr, "error running browser test :"+e.Error())
			os.Exit(1)
		} else if err2, ok2 := err.(*exec.ExitError); ok2 {
			_, _ = fmt.Fprintln(os.Stderr, string(err2.Stderr))
			os.Exit(1)
		} else {
			panic(err)
		}
	}

	_ = <- done // wait until main is done
	_ = cmd.Process.Kill() // stop the browser
}

