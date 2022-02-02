package http

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteDetector_WriteString(t *testing.T) {
	d := WriteDetector{}

	fn := func(w http.ResponseWriter, r *http.Request) {
		d.ResponseWriter = w
		io.WriteString(&d, "<html><body>Hello World!</body></html>")
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(fn).ServeHTTP(w, req)

	assert.True(t, d.HasWritten)
}

func TestWriteDetector_Write(t *testing.T) {
	d := WriteDetector{}

	fn := func(w http.ResponseWriter, r *http.Request) {
		d.ResponseWriter = w
		d.Write([]byte("<html><body>Hello World!</body></html>"))
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(fn).ServeHTTP(w, req)

	assert.True(t, d.HasWritten)
}

func TestWriteDetector_WriteHeader(t *testing.T) {
	d := WriteDetector{}

	fn := func(w http.ResponseWriter, r *http.Request) {
		d.ResponseWriter = w
		d.WriteHeader(400)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(fn).ServeHTTP(w, req)

	assert.True(t, d.HasWritten)
}

