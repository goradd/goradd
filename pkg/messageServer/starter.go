package messageServer

import (
	"fmt"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/sys"
	"log"
	"net/http"
)

func Start(mux *http.ServeMux) {
	MakeHub()

	if config.WebSocketPort != 0 {
		go func() {
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.WebSocketPort), mux))
		}()
	}

	if config.WebSocketTLSPort != 0 {
		// Here we confirm that the CertFile and KeyFile exist. If they don't, ListenAndServe just exits with an open error
		// and you will not know why.
		if !sys.PathExists(config.WebSocketTLSCertFile) {
			log.Fatalf("WebSocketTLSCertFile does not exist: %s", config.WebSocketTLSCertFile)
		}

		if !sys.PathExists(config.WebSocketTLSKeyFile) {
			log.Fatalf("WebSocketTLSKeyFile does not exist: %s", config.WebSocketTLSKeyFile)
		}

		go func() {
			log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.WebSocketTLSPort), config.WebSocketTLSCertFile, config.WebSocketTLSKeyFile, mux))
		}()
	}
}
