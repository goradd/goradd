// This application puts the browser tests into an application that will launch a server, run the tests
// and then exit with the results of the tests. It is designed to be used as part of the continuous-integration
// test process, and run in a headless browser when actually run as part of the ci test.
package main

import (
	log2 "github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/messageServer"
	_ "goradd-project/config" // Initialize required variables
	"goradd-project/web/app"
	"log"
	"net/http"
	// Below is where you import packages that register forms
	_ "goradd-project/web/form" // Your  forms.

	_ "github.com/goradd/goradd/pkg/bootstrap/examples" // Bootstrap examples
	_ "github.com/goradd/goradd/test/browsertest"
	_ "github.com/goradd/goradd/web/examples/controls"
	_ "goradd-project/gen" // Code-generated forms
)

func main() {
	var err error

	a := app.MakeApplication()
	log2.SetLogger(log2.FrameworkDebugLog, nil) // get rid of framework log for now
	messageServer.Start(a.MakeWebsocketMux())
	mux := a.MakeServerMux()

	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}

	// Now launch a browser with this address:  http://localhost:8000/test?all=1

}
