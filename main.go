package main

import (
	"flag"
	"fmt"
	"github.com/spekary/goradd/util"
	"net/http"
	"bytes"
	"path/filepath"
	"strings"
	"os"
	"go/build"
)

var port = flag.Int("p", 8082, "Start the webserver from the given port, example: -p 8082. Default is 8082.")

// Create other flags you might care about here

func main() {
	flag.Parse()

	if !installed() {
		//install()
	}
	err := runWebServer(*port)
	if err != nil {
		fmt.Println(err)
		return
	}
	//launchWebpage()
}

func runWebServer(port int) (err error) {

	mux := http.NewServeMux()

	// A very simple web server to act as an aid to configure and build a goradd app
	mux.Handle("/", serveHome())

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


func installed() bool {
	return util.PathExists(goraddPath())
}

func srcPath() string {
	return filepath.Join(goPath(), "src")
}

func goraddPath() string {
	return filepath.Join(srcPath(), "goradd")
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