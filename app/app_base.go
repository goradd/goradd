package app

import (
	//"flag"
	"github.com/spekary/goradd/page"
	"goradd/config"
	"net/http"
	"os"
	"path"
	"runtime"
)

// The application interface. A minimal set of commands that the main routine will ask the application to do.
// The main routine offers a way of creating mock applications, and alternate versions of the application from the default
type ApplicationI interface {
	Init(mode string)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	ProcessCommand([]string)
	PutContext(*http.Request) *http.Request
}

// The application base, to be embedded in your application
type Application struct {
}

func (a *Application) Init(mode string) {

	switch mode {
	case "debug":
		config.Mode = config.AppModeDebug
	case "rel":
		config.Mode = config.AppModeRelease
	case "dev":
		config.Mode = config.AppModeDevelopment
	default:
		panic("Unknown application mode")
	}
}

// Our application was called from the command line
func (a *Application) ProcessCommand(args []string) {
}

func (a *Application) PutContext(r *http.Request) *http.Request {
	return page.PutContext(r, os.Args[1:])
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	pm := page.GetPageManager()
	if pm == nil {
		panic("No page manager defined")
	}

	ctx := r.Context()
	buf := page.GetBuffer()
	defer page.PutBuffer(buf)
	if pm.IsPage(ctx) {
		headers, errCode := pm.RunPage(ctx, buf)
		if headers != nil {
			for k, v := range headers {
				// Multi-value headers can simply be separated with commas I believe
				w.Header().Set(k, v)
			}
		}
		if errCode != 0 {
			w.WriteHeader(errCode)
		}
		w.Write(buf.Bytes())
	}
}

func init() {
	var filename string
	var ok bool

	_, filename, _, ok = runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	config.GoraddDir = path.Dir(path.Dir(filename))
}
