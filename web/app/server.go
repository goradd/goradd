package app

import (
	"context"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/sys"
	"log"
	"net/http"
)

var server *http.Server

// ListenAndServeTLSWithTimeouts starts a secure web server with timeouts. The default http server does
// not have timeouts by default, which leaves the server open to certain attacks that would start
// a connection, but then very slowly read or write. Timeout values are taken from global variables
// defined in config, which you can set at init time.
func ListenAndServeTLSWithTimeouts(addr, certFile, keyFile string, handler http.Handler) error {
	// Here we confirm that the CertFile and KeyFile exist. If they don't, ListenAndServeTLS just exit with an open error
	// and you will not know why.
	if !sys.PathExists(certFile) {
		log.Fatalf("TLSCertFile does not exist: %s", config.TLSCertFile)
	}

	if !sys.PathExists(keyFile) {
		log.Fatalf("TLSKeyFile does not exist: %s", config.TLSKeyFile)
	}

	// TODO: https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/ recommends keeping track
	// of open connections using the ConnState hook for debugging purposes.

	server = &http.Server{
		Addr: addr,
		ReadTimeout:  config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
		Handler:      handler,
	}
	return server.ListenAndServeTLS(certFile, keyFile)
}

// ListenAndServeWithTimeouts starts a web server with timeouts. The default http server does
// not have timeouts, which leaves the server open to certain attacks that would start
// a connection, but then very slowly read or write. Timeout values are taken from global variables
// defined in config, which you can set at init time. This non-secure version is appropriate
// if you are serving behind another server, like apache or nginx.
func ListenAndServeWithTimeouts(addr string, handler http.Handler) error {

	// TODO: https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/ recommends keeping track
	// of open connections using the ConnState hook for debugging purposes.

	server = &http.Server{
		Addr: addr,
		ReadTimeout:  config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
		Handler:      handler,
	}
	return server.ListenAndServe()
}

// Shutdown performs a graceful shutdown of the server, returning any errors found.
func Shutdown(ctx context.Context) error {
	return server.Shutdown(ctx)
}
