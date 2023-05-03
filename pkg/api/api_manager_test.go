package api

import (
	http2 "github.com/goradd/goradd/pkg/http"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var fnFound = func(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Found")
}
var fnNotFound = func(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Not Found")
}

func TestRegisterAppPattern(t *testing.T) {

	type args struct {
		pattern string
		handler http.HandlerFunc
	}

	h := http2.UseAppMuxer(http2.NewMux(), http.HandlerFunc(fnNotFound))

	RegisterAppPattern("/a/", fnFound)

	req := httptest.NewRequest("GET", "/api/a", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp := w.Result()
	b, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Found", string(b))

	req = httptest.NewRequest("GET", "/api/a/", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp = w.Result()
	b, _ = io.ReadAll(resp.Body)
	assert.Equal(t, "Found", string(b))

}

func TestRegisterPattern(t *testing.T) {
	type args struct {
		pattern string
		handler http.HandlerFunc
	}

	h := http2.UsePatternMuxer(http2.NewMux(), http.HandlerFunc(fnNotFound))

	RegisterPattern("/a/", fnFound)

	req := httptest.NewRequest("GET", "/api/a", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp := w.Result()
	b, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Found", string(b))

	req = httptest.NewRequest("GET", "/api/a/", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp = w.Result()
	b, _ = io.ReadAll(resp.Body)
	assert.Equal(t, "Found", string(b))
}
