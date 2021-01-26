package config

import (
	"github.com/goradd/goradd/pkg/config"
	"path/filepath"
	"runtime"
)

func initGoradd() {
	setupPorts()

	setupDateFormats()
	setupTranslator()
	//setupBootstrap()

	if config.Release {
		//if you want to hardcode the asset directory location, do that here. Otherwise specify it on the command line.
		//config.SetAssetDirectory("serverDirLocation")

		//if you want to hardcode the html directory location, do that here. Otherwise specify it on the command line.
		//config.SetHtmlDirectory("htmlDirLocation")
	} else {
		_, filename, _, _ := runtime.Caller(0)

		// The projectDir points to files in the goradd-project directory. The development version would have all of these
		// files moved to a deployment location, so it is not available in the release version of the app. Doing the setup
		// this way ensures that when we build the release version, we will get a compile time failure if we accidentally try
		// to access the projectDir without making sure we are in the dev version of the app.
		projectDir := filepath.Dir(filepath.Dir(filename))
		config.SetProjectDir(projectDir)

		// This initializes the location of the static html directory for development. You can change it, but be sure to upload the files
		// to the server when you release and then point to them using the htmlDir flag when launching the application in server mode.
		// You can also comment it out if you are not using an html directory.
		config.SetHtmlDirectory(filepath.Join(config.ProjectDir(), "web", "html"))
	}

	//config.ApiPrefix = "/myapi" // Uncomment this to change the prefix for api calls (aka REST calls). The default is "/api".
	//config.WebsocketMessengerPrefix = "/ws/" // Sets the websocket messenger prefix. Set to blank to turn off the websocket messenger.
												// Otherwise, set to whatever prefix will not conflict with the rest of your app.

}

// setupPorts gives you an opportunity to hardcode the port values and certificate locations in your app.
// you can also put these on the command line
func setupPorts() {
	/*
		config.Port = 8000
		config.TLSPort = 8001 // This will require ssl certificates.

		// You will need to put in the path to your certfile and keyfile below.
		// The default implementation only uses these for the release build.
		config.TLSCertFile = ""
		config.TLSKeyFile = ""
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

/*
func setupBootstrap() {
	control.SetNewDialogFunction(func(form page.FormI, id string) control.DialogI {
		return control2.NewModal(form, id)
	})
}
*/