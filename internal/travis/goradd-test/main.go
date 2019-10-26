// This application puts the browser tests into an application that will launch a server, run the tests
// and then exit with the results of the tests. It is designed to be used as part of the continuous-integration
// test process, and run in a headless browser when actually run as part of the ci test.
package main

import (
	log2 "github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/messageServer/ws"
	"github.com/goradd/goradd/pkg/page"
	_ "goradd-project/config" // Initialize required variables
	"goradd-project/web/app"
	"log"
	"net/http"
	// Below is where you import packages that register forms
	_ "goradd-project/web/form" // Your  forms.

	_ "github.com/goradd/goradd/pkg/bootstrap/examples" // Bootstrap examples
	_ "github.com/goradd/goradd/test/browsertest"
	_ "github.com/goradd/goradd/web/examples/controls/panels"
	_ "goradd-project/gen" // Code-generated forms
)

func main() {
	var err error

	a := app.MakeApplication()
	log2.SetLogger(log2.FrameworkDebugLog, nil) // get rid of framework log for now
	ws.Start(a.MakeWebsocketMux())
	mux := a.MakeServerMux()

	// Make sure we always test with serialization turned on so that serialization is tested during the travis tests
	page.SetPageCache(page.NewSerializedPageCache(100, 60*60*24))

	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}

	// Now launch a browser with this address:  http://localhost:8000/goradd/Test.g?all=1

}
