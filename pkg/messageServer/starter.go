package messageServer

import (
	"fmt"
	"goradd-project/config"
	"log"
	"net/http"
)

func Start(mux *http.ServeMux) {
	MakeHub()

	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.GoraddWebSocketPort), mux))
	}()

	if config.Release {
		go func() {
			log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.GoraddWebSocketTLSPort), config.GoraddWebSocketTLSCertFile, config.GoraddWebSocketTLSKeyFile, mux))
		}()
	}
}

