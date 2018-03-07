package app

import (
	//"flag"
	"context"
	"net/http"
	"os"
	"github.com/spekary/goradd/page"
	"goradd/config"
	"runtime"
	"path"
)


// The application interface. A minimal set of commands that the main routine will ask the application to do.
// The main routine offers a way of creating mock applications, and alternate versions of the application from the default
type ApplicationI interface {
	Init(mode string)
	ProcessRequest(context.Context, http.ResponseWriter, *http.Request)
	ProcessCommand([]string)
	PutContext(context.Context, *http.Request) context.Context
}


// The application base, to be embedded in your application
type Application struct {
}

func (a *Application) Init(mode string) {

	switch mode{
	case "debug":
		config.Mode = config.Debug
	case "rel":
		config.Mode = config.Rel
	case "dev":
		config.Mode = config.Dev
	default:
		panic ("Unknown application mode")
	}
}

// Our application was called from the command line
func (a *Application) ProcessCommand (args []string) {
}


func (a *Application) PutContext(ctx context.Context, r *http.Request) context.Context {
	grctx := &page.Context{}
	grctx.FillFromRequest(os.Args[1:], r)
	return context.WithValue(ctx, "goradd", grctx)
}

func (a *Application) ProcessRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	/* ToDo
	   Problem: If someone duplicates a page, which is easily doable in chrome, two pages will be sharing the same form state. We should detect this
	   somehow and respond by creating a new formstate for each
	   Solution: Return a request counter that increments on each request. If the counter is not in sync, then respond with
	   an error and regenerate a new formstate. The request therefore must
	   state whether it should be a synchronous or asynchronous request so that we ignore the counter on async requests
	   (should be rare, and client is responsible for specifying)
	*/

	//go func() {
	pm := page.GetPageManager()
	if pm == nil {
		panic ("No page manager defined")
	}

		if pm.IsPage(ctx) {
			pm.RunPage(ctx, w) // Use a go routine to make sure our http server is not blocking in general.
		} else if pm.IsAsset(ctx) {
			pm.ServeAsset(ctx, w, r)
		} else {
			// pass this on to ServeFile or ServeContent
			// or return
		}

	//}
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