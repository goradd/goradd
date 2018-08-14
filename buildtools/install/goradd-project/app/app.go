package app

import (
	"github.com/spekary/goradd/app"
	"github.com/spekary/goradd/html"
	grlog "github.com/spekary/goradd/log"
	"github.com/spekary/goradd/page"
	"goradd-project/config"
	"log"
	"net/http"
	"os"
)

type Application struct {
	app.Application

	//dbs [1]db.DB

	// Your own vars, methods and overrides
}

func (a *Application) initDatabases() {
	//db := db.NewMysql5("username/password")

	//a.dbs[0] = db
}

func (a *Application) GetDb(i int) {
	//return a.dbs[i]
}

func (a *Application) Test() {
	//var p model.Person
	//var p2 model.People

	//p.Load(1)
	//p2.Load()
}

func (a *Application) Init(mode string) {
	a.Application.Init(mode)
	a.initDatabases()

	// Replace this if you would like a different error display
	page.ErrorPageFunc = page.DefaultErrorPageTmpl

	// Control how pages are cached. This will vary depending on whether you are using multiple machines to run your app, and whether you are in development mode, etc.
	page.SetPageCache(page.NewFastPageCache())

	// Control how pages are serialized if a serialization cache is being used
	page.SetPageEncoder(page.GobPageEncoder{})

	grlog.Loggers[grlog.InfoLog] = log.New(os.Stdout, "Info: ", log.Ldate|log.Ltime)
	grlog.Loggers[grlog.FrameworkDebugLog] = log.New(os.Stdout, "Framework: ", log.Ldate|log.Ltime|log.Lshortfile)
	grlog.Loggers[grlog.DebugLog] = log.New(os.Stdout, "AppModeDebug: ", log.Ldate|log.Ltime|log.Lshortfile)
	grlog.Loggers[grlog.WarningLog] = log.New(os.Stderr, "Warning: ", log.Ldate|log.Ltime|log.Llongfile)
	grlog.Loggers[grlog.ErrorLog] = log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime|log.Llongfile)

	// Control whether to try to minify output. The default minifies if we are in release mode.
	if config.Mode == config.AppModeRelease {
		config.Minify = true
	}

	page.DefaultCheckboxLabelDrawingMode = html.LABEL_AFTER

}

// PutContext allocates a blank context object for our application specific context data. It can be populated later.
func (a *Application) PutContext(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = PutContext(ctx)
	r = r.WithContext(ctx)
	// be sure to call the superclass version so the goradd framework can operate
	return a.Application.PutContext(r)
}

