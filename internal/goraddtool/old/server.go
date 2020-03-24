package old

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/goradd/goradd/pkg/sys"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var port = flag.Int("p", 8082, "Start the webserver from the given port, example: -p 8082. Default is 8082.")
var dev = flag.Bool("dev", false, "WARNING! Automatically erases the installation directories. Only use this if you are working on the installer itself. Not for general use.")
var installed bool
var results string // result of whatever the current operation is
var stop bool
var cmd string
var cwd string // path we were launched from. This will be the installation directory.
var modules map[string]string

func Launch() {
	var err error
	if cwd, err = os.Getwd(); err != nil {
		log.Fatal(fmt.Errorf("could not get the working directory: %s", err.Error()))
	}

	if len(os.Args) < 2 {
		log.Fatal("")
	}

	flag.Parse()

	if *dev {
		// Erase the installation directory so that it appears that we are in a new install.
		// Use this only as an aid to developing the installer.
		os.RemoveAll(projectPath())
	}
	installed = isInstalled()

	launchWebpage() // hopefully by the time the web page gets to this address, the server will have started
	err = runWebServer(*port)
	if err != nil {
		log.Fatal(fmt.Errorf("there was a problem running the web server: %s", err.Error()))
	}
	return
}

func runWebServer(port int) (err error) {

	mux := http.NewServeMux()

	// A very simple web server to act as an aid to configure and build a goradd app
	mux.Handle("/", serveHome())
	mux.Handle("/installer", serveInstall())
	mux.Handle("/builder", serveBuilder())

	// The two "Serve" functions below will launch go routines for each request, so that multiple requests can be
	// processed in parallel. This may mean multiple requests for the same override, depending on the structure of the override.
	addr := fmt.Sprintf(":%d", port)
	err = http.ListenAndServe(addr, mux)

	return err
}

// serveHome serves the main page
func serveHome() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		drawHome(buf)
		w.Write(buf.Bytes())
	}
	return http.HandlerFunc(fn)
}

func isInstalled() bool {
	return sys.PathExists(projectPath())
}

func srcPath() string {
	return filepath.Dir(cwd)
}

func projectPath() string {
	return filepath.Join(srcPath(), "goradd-project")
}

func goraddPath() string {
	var err error
	if modules == nil {
		if modules, err = sys.ModulePaths(); err != nil {
			return ""
		}
	}

	if v, ok := modules["github.com/goradd/goradd"]; ok {
		return v
	}
	if v, ok := modules["github.com"]; ok {
		return filepath.Join(srcPath(), v, "goradd", "goradd")
	}
	return ""
}

func launchWebpage() {

	switch runtime.GOOS {
	case `darwin`:
		_, _, err := executeCmd(`open`, fmt.Sprintf("http://localhost:%d", *port))
		if err == nil {
			return
		}
	case `windows`:
		_, _, err := executeCmd(`start`, fmt.Sprintf("http://localhost:%d", *port))
		if err == nil {
			return
		}
	}
	fmt.Printf("The goradd server is running. Go to http://localhost:%d in a browser", *port)
}

func executeCmd(command string, args ...string) (stdOutText string, stdErrText string, err2 error) {
	cmd := exec.Command(command, args...)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		err2 = err
		return
	}

	defer stdOut.Close()

	scanner := bufio.NewScanner(stdOut)
	go func() {
		for scanner.Scan() {
			stdOutText += scanner.Text() + "\n"
		}
	}()

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		err2 = err
		return
	}

	defer stdErr.Close()

	stdErrScanner := bufio.NewScanner(stdErr)
	go func() {
		for stdErrScanner.Scan() {

			stdErrText += stdErrScanner.Text() + "\n"
		}
	}()

	err = cmd.Start()
	if err != nil {
		err2 = err
		return
	}

	err = cmd.Wait()

	if err != nil {
		err2 = err
	}
	return
}
