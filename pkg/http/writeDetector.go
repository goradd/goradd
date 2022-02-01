package http

import (
	"net/http"
)

type ResponseRewinder interface {
	http.ResponseWriter
	Rewind(n int)
}
// WriteDetector is a utility for Handlers to detect whether a sub-handler has responded to an http request.
type WriteDetector struct {
	ResponseRewinder
	HasWritten bool
}

func (d *WriteDetector) Write(b []byte) (l int, err error) {
	d.HasWritten = true
	return d.ResponseRewinder.Write(b)
}

func (d *WriteDetector) WriteHeader(code int) {
	d.HasWritten = true
	d.ResponseRewinder.WriteHeader(code)
}

