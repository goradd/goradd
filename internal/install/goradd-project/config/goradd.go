package config

import (
	"github.com/goradd/goradd/pkg/config"
	"path/filepath"
)


func initGoradd() {
	setupPorts()

	setupDateFormats()
	setupTranslator()
	setupCodegen()


	if config.Release {
		//if you want to hardcode the asset directory location, do that here. Otherwise specify it on the command line.
		//config.SetAssetDirectory("serverDirLocation")

		//if you want to hardcode the html directory location, do that here. Otherwise specify it on the command line.
		//config.SetHtmlDirectory("htmlDirLocation")
	} else {
		// This initializes the location of the static html directory for development. You can change it, but be sure to upload the files
		// to the server when you release and then point to them using the htmlDir flag when launching the application in server mode.
		// You can also comment it out if you are not using an html directory.
		config.SetHtmlDirectory(filepath.Join(config.ProjectDir(), "web","html"))
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

func setupCodegen() {
	// Setup codegen customizations here
	//generator.DefaultWrapper = "bootstrap.FormGroup"
}

