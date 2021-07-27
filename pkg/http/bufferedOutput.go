package http

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/goradd"
	grlog "github.com/goradd/goradd/pkg/log"
	buf2 "github.com/goradd/goradd/pkg/pool"
	"net/http"
)

var bufferedOutputManager User

// SetBufferedOutputManager injects the given manager as the global buffered output manager
func SetBufferedOutputManager(u User) {
	bufferedOutputManager = u
}

func BufferedOutputManager() User {
	return bufferedOutputManager
}

type BufferedResponseWriterI interface {
	http.ResponseWriter
	Disable()
	OutputBuffer() *bytes.Buffer
}

type bufferedResponseWriter struct {
	http.ResponseWriter
	buf  *bytes.Buffer
	code int
	disabled bool
}

func (bw *bufferedResponseWriter) Write(b []byte) (int, error) {
	if bw.disabled {
		return bw.ResponseWriter.Write(b)
	}
	return bw.buf.Write(b)
}

func (bw *bufferedResponseWriter) WriteHeader(code int) {
	if bw.disabled {
		bw.ResponseWriter.WriteHeader(code)
		return
	}
	bw.code = code
}

// Disable will turn off the buffered response and send responses directly to
// the response writer above it in the response writer chain.
//
// Some parts of goradd depend on buffered responses. In particular, the session manager
// relies on buffered responses. If you call disable, you will not be able to change
// session data, but you should still be able to read session data.
func (bw *bufferedResponseWriter) Disable()  {
	bw.disabled = true
}

// OutputBuffer returns the output buffer
func (bw *bufferedResponseWriter) OutputBuffer() *bytes.Buffer {
	return bw.buf
}

func DisableOutputBuffering(ctx context.Context) {
	ctx.Value(goradd.BufferContext).(BufferedResponseWriterI).Disable()
}

func OutputBuffer(ctx context.Context) *bytes.Buffer {
	return ctx.Value(goradd.BufferContext).(BufferedResponseWriterI).OutputBuffer()
}

type defaultBufferedOutputManager struct {
}

// Use injects the buffered output manager into the handler stack.
func (bw *defaultBufferedOutputManager) Use(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		// Setup the output buffer
		outBuf := buf2.GetBuffer()
		bw := &bufferedResponseWriter{w, outBuf, 0, false}
		ctx := r.Context()
		ctx = context.WithValue(ctx, goradd.BufferContext, bw)
		r = r.WithContext(ctx)
		defer buf2.PutBuffer(outBuf)
		next.ServeHTTP(bw, r)
		if bw.code != 0 && bw.code != 200 {
			grlog.Error("Buffered write error code ", bw.code)
			w.WriteHeader(bw.code)
		}
		_, e := w.Write(outBuf.Bytes())
		if e != nil {
			grlog.Error("Buffered write error ", e.Error())
		}
		//log.Printf("Buffered write %d bytes %v %s", i, w.Header(), outBuf.String())
	}
	return http.HandlerFunc(fn)
}



func init() {
	// Initialize the default buffered output handler
	SetBufferedOutputManager(new(defaultBufferedOutputManager))
}
