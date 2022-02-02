package http

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)


func Test_BufferedOutput(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_,_ = io.WriteString(w, "Hey")
		b := outputBuffer(r.Context())
		assert.EqualValues(t,"Hey", b.String(), "Output was not buffered")
		assert.Equal(t, 3, OutputLen(r.Context()), "Len was not recorded")
	}
	h := BufferedOutputManager().Use(http.HandlerFunc(fn))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assert.EqualValues(t, "Hey", body, "Output was not forwarded to writer")
}

func Test_UnbufferedOutput(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		DisableOutputBuffering(r.Context())
		_,_ = io.WriteString(w, "Hey")
		b := outputBuffer(r.Context())
		assert.NotEqualValues(t,"Hey", b.String(), "Output was buffered but should not be")
		assert.Equal(t, 3, OutputLen(r.Context()), "Len was not recorded")
	}
	h := BufferedOutputManager().Use(http.HandlerFunc(fn))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assert.EqualValues(t, "Hey", body, "Output was not forwarded to writer")
}

func Test_BufferedOutputCode(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_,_ = io.WriteString(w, "Hey")
		w.WriteHeader(300)

		b := outputBuffer(r.Context())
		assert.EqualValues(t,"Hey", b.String(), "Output was not buffered")
		assert.Equal(t, 3, OutputLen(r.Context()), "Len was not recorded")
	}
	h := BufferedOutputManager().Use(http.HandlerFunc(fn))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assert.EqualValues(t, "Hey", body, "Output was not forwarded to writer even when header was prior written")
	assert.Equal(t, 300, resp.StatusCode, "Code was not recorded when sent after first response.")
}

func Test_UnbufferedOutputCode(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		DisableOutputBuffering(r.Context())
		_,_ = io.WriteString(w, "Hey")
		w.WriteHeader(300)

		b := outputBuffer(r.Context())
		assert.NotEqualValues(t,"Hey", b.String(), "Output was not buffered")
		assert.Equal(t, 3, OutputLen(r.Context()), "Len was not recorded")
	}
	h := BufferedOutputManager().Use(http.HandlerFunc(fn))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assert.EqualValues(t, "Hey", body, "Output was not forwarded to writer even when header was prior written")
	assert.Equal(t, 200, resp.StatusCode, "Code was recorded even when sent after first response.")

}

func outputBuffer(ctx context.Context) *bytes.Buffer {
	return ctx.Value(goradd.BufferContext).(BufferedResponseWriterI).OutputBuffer()
}
