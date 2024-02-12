package config

import (
	"github.com/goradd/goradd/pkg/config"
	"path/filepath"
	"runtime"
)

// This file is for code that changes global variables that are part of the GoRADD framework.

func initGoradd() {
	setupDateFormats()
	setupTranslator()
	//setupBootstrap()
	setupErrorMessage()

	if config.Release {
	} else {
		_, filename, _, _ := runtime.Caller(0)

		// The projectDir points to files in the goradd-project directory. The development version would have all of these
		// files moved to a deployment location, so it is not available in the release version of the app. Doing the setup
		// this way ensures that when we build the release version, we will get a compile time failure if we accidentally try
		// to access the projectDir without making sure we are in the dev version of the app.
		projectDir := filepath.Dir(filepath.Dir(filename))
		config.SetProjectDir(projectDir)
	}

	//config.ApiPrefix = "/myapi" // Uncomment this to change the prefix for api calls (aka REST calls). The default is "/api".
	//config.WebsocketMessengerPrefix = "/ws/" // Sets the websocket messenger prefix. Set to blank to turn off the websocket messenger.
	// Otherwise, set to whatever prefix will not conflict with the rest of your app.

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
	dialog.SetNewDialogFunction(func(form page.FormI, id string) dialog.DialogI {
		return control2.NewModal(form, id)
	})
}

*/

// setupErrorMessage allows you to customize the error message that will appear to users if the code panics.
func setupErrorMessage() {
	/*
			page.HtmlErrorMessage = `<h1 id="err-title">Error</h1>
		<p>
		An unexpected error has occurred and your request could not be processed. The error has been logged and we will
		attempt to fix the problem as soon as possible.
		</p>`*/
}
