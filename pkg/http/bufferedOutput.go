package http

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/goradd"
	grlog "github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/pool"
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
	Len() int
}

type bufferedResponseWriter struct {
	http.ResponseWriter
	buf  *bytes.Buffer
	code int
	disabled bool
	len int
}

func (bw *bufferedResponseWriter) Write(b []byte) (l int, err error) {
	if bw.disabled {
		l, err = bw.ResponseWriter.Write(b)
	} else {
		l,err = bw.buf.Write(b)
	}
	bw.len += l
	return l,err
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

// OutputBuffer returns the current output buffer.
func (bw *bufferedResponseWriter) OutputBuffer() *bytes.Buffer {
	return bw.buf
}

func (bw *bufferedResponseWriter) Len() int {
	if bw.disabled {
		return bw.len
	} else {
		return bw.buf.Len()
	}
}

func DisableOutputBuffering(ctx context.Context) {
	ctx.Value(goradd.BufferContext).(BufferedResponseWriterI).Disable()
}

// ResetOutputBuffer returns the current output buffer and resets the output buffer to nothing.
func ResetOutputBuffer(ctx context.Context) []byte {
	buf := ctx.Value(goradd.BufferContext).(BufferedResponseWriterI).OutputBuffer()
	ret := buf.Bytes()
	buf.Reset()
	return ret
}
// OutputLen returns the number of bytes written to the output, even if
// output buffering is disabled.
func OutputLen(ctx context.Context) int {
	return ctx.Value(goradd.BufferContext).(BufferedResponseWriterI).Len()
}

type defaultBufferedOutputManager struct {
}

// Use injects the buffered output manager into the handler stack.
func (bw *defaultBufferedOutputManager) Use(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		// Setup the output buffer
		outBuf := pool.GetBuffer()
		bwriter := &bufferedResponseWriter{w, outBuf, 0, false, 0}
		ctx := r.Context()
		ctx = context.WithValue(ctx, goradd.BufferContext, bwriter)
		r = r.WithContext(ctx)
		defer pool.PutBuffer(outBuf)
		next.ServeHTTP(bwriter, r)
		if bwriter.code != 0 && bwriter.code != 200 {
			grlog.Error("Buffered write error code ", bwriter.code)
			grlog.Error(r.URL.Path)
			w.WriteHeader(bwriter.code)
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
