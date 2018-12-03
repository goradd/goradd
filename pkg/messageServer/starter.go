package messageServer

import (
	"fmt"
	"github.com/spekary/goradd/pkg/config"
	"log"
	"net/http"
)

func Start(mux *http.ServeMux) {
	MakeHub()

	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.WebSocketPort), mux))
	}()

	if config.WebSocketTLSPort != 0 {
		go func() {
			log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.WebSocketTLSPort), config.WebSocketTLSCertFile, config.WebSocketTLSKeyFile, mux))
		}()
	}
}

