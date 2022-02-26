package http

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/goradd/goradd/pkg/goradd"
	grlog "github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/pool"
)

var bufferedOutputManager User = &defaultBufferedOutputManager{}

// SetBufferedOutputManager injects the given manager as the global buffered output manager.
// Call this to change it to your own version.
//
// Your custom version should create a BufferedResponseWriterI object and attach it to the
// goradd.BufferContext value in the context. See defaultBufferedOutputManager.Use for
// an example.
func SetBufferedOutputManager(u User) {
	bufferedOutputManager = u
}

// BufferedOutputManager returns the buffered output manager.
func BufferedOutputManager() User {
	return bufferedOutputManager
}

// BufferedResponseWriterI is the interface for a BufferedResponseWriter.
type BufferedResponseWriterI interface {
	http.ResponseWriter
	Disable()
	OutputBuffer() *bytes.Buffer
	Len() int
}

type bufferedResponseWriter struct {
	http.ResponseWriter
	buf      *bytes.Buffer
	code     int
	disabled bool
	len      int
}

// Write writes to the BufferedResponseWriter.
// The length of what was written, and any errors, are returned.
func (bw *bufferedResponseWriter) Write(b []byte) (l int, err error) {
	// TODO: Set max buffer size
	if bw.disabled {
		l, err = bw.ResponseWriter.Write(b)
	} else {
		l, err = bw.buf.Write(b)
	}
	bw.len += l
	return l, err
}

// WriteHeader writes an error code to the BufferedResponseWriter.
//
// If output buffering is enabled, you can call this function multiple times, and
// as opposed to the standard ResponseWriter, calling WriteHeader does not actually
// send the code to the output, but stores it until the buffer is written.
func (bw *bufferedResponseWriter) WriteHeader(code int) {
	if bw.disabled {
		bw.ResponseWriter.WriteHeader(code)
		return
	}
	bw.code = code
}

// WriteString satisfies the StringWriter interface for WriteString optimization.
func (bw *bufferedResponseWriter) WriteString(s string) (l int, err error) {
	if bw.disabled {
		l, err = io.WriteString(bw.ResponseWriter, s)
	} else {
		l, err = bw.buf.WriteString(s)
	}
	bw.len += l
	return l, err
}

// Disable should turn off the buffered response and send responses directly to
// the response writer above it in the response writer chain.
func (bw *bufferedResponseWriter) Disable() {
	bw.disabled = true
}

// OutputBuffer returns the current output buffer.
func (bw *bufferedResponseWriter) OutputBuffer() *bytes.Buffer {
	return bw.buf
}

// Len returns the total length of the data that has been written to the output buffer.
func (bw *bufferedResponseWriter) Len() int {
	if bw.disabled {
		return bw.len
	} else {
		return bw.buf.Len()
	}
}

// DisableOutputBuffering turns off output buffering.
//
// This function is useful if you want to stream large quantities of data to the response writer,
// and you want to avoid the memory allocation required for buffering.
// The default output buffer manager will not output what has been buffered so far.
// You should not try to re-enable output buffering.
func DisableOutputBuffering(ctx context.Context) {
	ctx.Value(goradd.BufferContext).(BufferedResponseWriterI).Disable()
}

// ResetOutputBuffer returns the current output buffer and resets the output buffer to nothing.
func ResetOutputBuffer(ctx context.Context) []byte {

	if i := ctx.Value(goradd.BufferContext); i == nil {
		return nil
	} else if bw := i.(BufferedResponseWriterI); bw == nil {
		return nil
	} else if buf := bw.OutputBuffer(); buf == nil {
		return nil
	} else {
		ret := buf.Bytes()
		buf.Reset()
		return ret
	}

}

// OutputLen returns the number of bytes written to the output.
//
// If output buffering is disabled, this will be the number of bytes actually written to
// the response writer. If output buffering is enabled, this is the number of bytes in the
// buffer waiting to be sent.
func OutputLen(ctx context.Context) int {
	return ctx.Value(goradd.BufferContext).(BufferedResponseWriterI).Len()
}

type defaultBufferedOutputManager struct {
}

// Use injects the buffered output manager into the handler stack.
func (bw *defaultBufferedOutputManager) Use(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		// Set up the output buffer
		outBuf := pool.GetBuffer()
		bwriter := &bufferedResponseWriter{w, outBuf, 0, false, 0}
		ctx := r.Context()
		ctx = context.WithValue(ctx, goradd.BufferContext, bwriter)
		r = r.WithContext(ctx)
		defer pool.PutBuffer(outBuf)
		next.ServeHTTP(bwriter, r)
		if bwriter.code != 0 && bwriter.code != 200 {
			w.WriteHeader(bwriter.code)
		}
		_, e := w.Write(outBuf.Bytes())
		if e != nil {
			grlog.Error("Buffered write error ", e.Error())
		}
	}
	return http.HandlerFunc(fn)
}
