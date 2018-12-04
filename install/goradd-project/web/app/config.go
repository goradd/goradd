package app

import (
	"github.com/spekary/goradd/pkg/config"
	"goradd-project/config/dbconfig"
	"path/filepath"
	"runtime"
)

func configure() {
	dbconfig.InitDatabases()

	setupPorts()

	setupJQuery()
	setupDateFormats()
	setupTranslator()

	if config.Release {
		//if you want to hardcode the asset directory location, do that here. Otherwise specify it on the command line.
		//config.SetAssetDirectory("serverDirLocation")

		//if you want to hardcode the html directory location, do that here. Otherwise specify it on the command line.
		//config.SetHtmlDirectory("htmlDirLocation")
	} else {
		// You can also use a different location for your static files in development mode. Just be sure to upload them
		// to the server when you release and then point to them using the htmlDir flag when launching the application in server mode.
		//config.SetHtmlDirectory("htmlDirLocation")
	}

}

// setupPorts gives you an opportunity to hardcode the port values and certificate locations in your app.
// you can also put these on the command line
func setupPorts() {
	/*
	//config.UseFCGI = true

	config.Port = 8000
	config.TLSPort = 8001 // This will require ssl certificates.

	// You will need to put in the path to your certfile and keyfile below.
	// The default implementation only uses these for the release build.
	config.TLSCertFile = ""
	config.TLSKeyFile = ""

	config.WebSocketPort = 8101
	config.WebSocketTLSPort = 8102 // This will require ssl certificates.

	// You will need to put in the path to your certfile and keyfile below.
	// The default implementation only uses these for the release build.
	// You can use the same ones that you use for normal SSL communication over http.
	config.WebSocketTLSCertFile = config.TLSCertFile
	config.WebSocketTLSKeyFile = config.TLSKeyFile
*/
}

func setupJQuery() {
	// These are optional setups to control how jQuery and jQueryUI are retrieved from the web

	/*
	if config.Release {
		// This will get jquery from jquery's own code base
		config.SetJqueryPath("https://code.jquery.com/jquery-3.3.1.min.js", map[string]string{"integrity": "sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=", "crossorigin": "anonymous"})
	} else {
		config.SetJqueryPath (path.Join(config.GoraddAssets(),"/js/jquery3.js"), nil)
	}
	*/

	/*
	if config.Release {
		config.SetJqueryUIPath("https://code.jquery.com/ui/1.12.1/jquery-ui.min.js", map[string]string{"integrity": "sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=", "crossorigin": "anonymous"})
	} else {
		config.SetJqueryUIPath(path.Join(config.GoraddAssets(),"/js/jquery-ui.js"), nil)
	}
	*/
}

func setupDateFormats() {
	// uncomment below to change the default formats used to display dates
	// these may get deprecated in favor of using something that is localized to the browser's locale.
	/*
	config.DefaultDateFormat = "January 2, 2006"
	config.DefaultTimeFormat = "3:04 am"
	config.DefaultDateTimeFormat = "January 2, 2006 3:04am"
	*/
}

func setupTranslator() {
	// Here is where you would insert your translator as the global translation service.
	// i18n.SetTranslator(myTranslator)
}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	config.SetProjectDir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
}
