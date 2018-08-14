package buildtools


import (
	"net/http"
	"fmt"
	"bytes"
	"github.com/spekary/goradd/util"
	"path/filepath"
	"strings"
	"os"
	"go/build"
	"flag"
	"runtime"
	"bufio"
	"os/exec"
)

var port = flag.Int("p", 8082, "Start the webserver from the given port, example: -p 8082. Default is 8082.")
var installed bool
var results string // result of whatever the current operation is
var stop bool

func Launch() {
	flag.Parse()

	os.RemoveAll(projectPath()) // TODO: Delete this line. Just for testing.
	installed = isInstalled()

	launchWebpage() // hopefully by the time the web page gets to this address, the server will have started
	err := runWebServer(*port)
	if err != nil {
		fmt.Println("There was a problem running the web server.")
		fmt.Println(err)
		return
	}

}


func runWebServer(port int) (err error) {

	mux := http.NewServeMux()

	// A very simple web server to act as an aid to configure and build a goradd app
	mux.Handle("/", serveHome())
	mux.Handle("/install", serveInstall())

	// The two "Serve" functions below will launch go routines for each request, so that multiple requests can be
	// processed in parallel. This may mean multiple requests for the same override, depending on the structure of the override.
	addr := fmt.Sprintf(":%d", port)
	err = http.ListenAndServe(addr, mux)

	return err
}


// serveHome serves the main page
func serveHome() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		buf := new (bytes.Buffer)
		drawHome(buf)
		w.Write(buf.Bytes())
	}
	return http.HandlerFunc(fn)
}


func isInstalled() bool {
	return util.PathExists(projectPath())
}

func srcPath() string {
	return filepath.Join(goPath(), "src")
}

func projectPath() string {
	return filepath.Join(srcPath(), "goradd-project")
}

func goraddPath() string {
	return filepath.Join(srcPath(), "github.com", "spekary", "goradd")
}

func goPath() string {
	goPaths := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))
	if len(goPaths) == 0 {
		return build.Default.GOPATH
	} else if goPaths[0] == "" {
		return build.Default.GOPATH
	}
	return goPaths[0]
}

func launchWebpage() {

	switch runtime.GOOS {
	case `darwin`:
		_,_,err := executeCmd(`open`,  fmt.Sprintf("http://localhost:%d", *port))
		if err == nil {
			return
		}
	case `windows`:
		_,_,err := executeCmd(`start`,  fmt.Sprintf("http://localhost:%d", *port))
		if err == nil {
			return
		}
	}
	fmt.Sprintln("The goradd server is running. Go to http://localhost:%d in a browser", *port)
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
