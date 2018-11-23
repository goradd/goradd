package app

import (
	"github.com/spekary/goradd/pkg/html"
	grlog "github.com/spekary/goradd/pkg/log"
	"github.com/spekary/goradd/pkg/page"
	"github.com/spekary/goradd/web/app"
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


func (a *Application) GetDb(i int) {
	//return a.dbs[i]
}

func (a *Application) Test() {
	//var p model.Person
	//var p2 model.People

	//p.Load(1)
	//p2.Load()
}

func (a *Application) Init() {
	a.Application.Init()

	// Replace this if you would like a different error display
	page.ErrorPageFunc = page.DefaultErrorPageTmpl

	// Control how pages are cached. This will vary depending on whether you are using multiple machines to run your app, and whether you are in development mode, etc.
	page.SetPageCache(page.NewFastPageCache())

	// Control how pages are serialized if a serialization cache is being used
	page.SetPageEncoder(page.GobPageEncoder{})

	// Framework error log
	grlog.Loggers[grlog.FrameworkDebugLog] = log.New(os.Stdout, "Framework: ", log.Ldate|log.Ltime|log.Lshortfile)

	// App specific loggers
	grlog.Loggers[grlog.InfoLog] = log.New(os.Stdout, "Info: ", log.Ldate|log.Ltime)
	grlog.Loggers[grlog.WarningLog] = log.New(os.Stderr, "Warning: ", log.Ldate|log.Ltime|log.Llongfile)
	grlog.Loggers[grlog.ErrorLog] = log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime|log.Llongfile)
	grlog.Loggers[grlog.DebugLog] = log.New(os.Stdout, "Debug: ", log.Ldate|log.Ltime|log.Lshortfile)

	page.DefaultCheckboxLabelDrawingMode = html.LabelAfter

	page.RegisterAssetDirectory(config.GoraddAssets(), config.AssetPrefix + "goradd")
	page.RegisterAssetDirectory(config.ProjectAssets(), config.AssetPrefix + "project")
}

// PutContext allocates a blank context object for our application specific context data. It can be populated later.
func (a *Application) PutContext(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = PutContext(ctx)
	r = r.WithContext(ctx)
	// be sure to call the superclass version so the goradd framework can operate
	return a.Application.PutContext(r)
}

