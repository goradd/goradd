// This application puts the browser tests into an application that will launch a server, run the tests
// and then exit with the results of the tests. It is designed to be used as part of the continuous-integration
// test process, and run in a headless browser when actually run as part fo the ci test.
package main

import (
	"github.com/goradd/goradd/pkg/messageServer"
	"github.com/goradd/goradd/pkg/sys"
	_ "goradd-project/config" // Initialize required variables
	"goradd-project/web/app"
	"log"
	"net/http"
	// Below is where you import packages that register forms
	_ "goradd-project/web/form" // Your  forms.

	_ "github.com/goradd/goradd/pkg/bootstrap/examples" // Bootstrap examples
	_ "github.com/goradd/goradd/test/browser"
	_ "github.com/goradd/goradd/test/page"
	_ "goradd-project/gen" // Code-generated forms
)


func main() {
	var err error

	a := app.MakeApplication()
	messageServer.Start(a.MakeWebsocketMux())
	mux := a.MakeServerMux()

	// This will launch all the registered browser based tests and then exit with the results
	err = sys.LaunchDefaultBrowser("http://localhost:8000/test?all=1")
	if err != nil {
		log.Fatal(err)
	}

	//go func() {
		err = http.ListenAndServe(":8000", mux)
		if err != nil {
			log.Fatal(err)
		}
	//}()

}

