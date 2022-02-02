package http

import (
	"io"
	"net/http"
)


// WriteDetector is a utility for Handlers to detect whether a sub-handler has responded to an HTTP request.
type WriteDetector struct {
	http.ResponseWriter
	HasWritten bool
}

func (d *WriteDetector) Write(b []byte) (l int, err error) {
	d.HasWritten = true
	return d.ResponseWriter.Write(b)
}

func (d *WriteDetector) WriteHeader(code int) {
	d.HasWritten = true
	d.ResponseWriter.WriteHeader(code)
}

// WriteString will cause WriteString to be used on the enclosed ResponseWriter if the ResponseWriter
// supports it.
func (d *WriteDetector) WriteString(s string) (l int, err error)  {
	d.HasWritten = true
	return io.WriteString(d.ResponseWriter, s)
}

