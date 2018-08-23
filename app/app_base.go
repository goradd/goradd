package app

import (
	//"flag"
	"github.com/spekary/goradd/page"
	"net/http"
	"os"
)

// The application interface. A minimal set of commands that the main routine will ask the application to do.
// The main routine offers a way of creating mock applications, and alternate versions of the application from the default
type ApplicationI interface {
	Init()
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	ProcessCommand([]string)
	PutContext(*http.Request) *http.Request
}

// The application base, to be embedded in your application
type Application struct {
}

func (a *Application) Init() {

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
		panic("No override manager defined")
	}

	ctx := r.Context()
	grctx := page.GetContext(ctx)
	buf := grctx.AppContext.OutBuf
	if pm.IsPage(grctx) {
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
	}
}
