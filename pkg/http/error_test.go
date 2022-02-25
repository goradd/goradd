package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorReporter_Use(t *testing.T) {
	tests := []struct {
		name     string
		f        func()
		want     []byte
		wantCode int
	}{
		{"std panic", func() { panic("test") }, []byte{}, 500},
		{"int panic", func() { panic(300) }, []byte{}, 300},
		{"send message", func() { SendErrorMessage("test", 300) }, []byte("test"), 300},
		{"redirect", func() { Redirect("loc", 300) }, []byte{}, 300},
		{"unauthorized", func() { SendUnauthorized() }, []byte{}, http.StatusUnauthorized},
		{"send error", func() { SendErrorCode(501) }, []byte{}, 501},
		{"forbidden", func() { SendForbidden() }, []byte{}, http.StatusForbidden},
		{"method not allowed", func() { SendMethodNotAllowed("GET") }, []byte{}, http.StatusMethodNotAllowed},
		{"not found", func() { SendNotFound() }, []byte{}, http.StatusNotFound},
		{"not found message", func() { SendNotFoundMessage("test") }, []byte("test"), http.StatusNotFound},
		{"bad request", func() { SendBadRequest() }, []byte{}, http.StatusBadRequest},
		{"bad request message", func() { SendBadRequestMessage("test") }, []byte("test"), http.StatusBadRequest},
		{"server error", func() { panic(NewServerError("msg", "mode", nil, 2, "HtmlErrorMessage")) }, []byte("HtmlErrorMessage"), 500},
		{"std error", func() { panic(fmt.Errorf("test")) }, []byte{}, 500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := func(w http.ResponseWriter, r *http.Request) {
				tt.f()
			}
			e := ErrorReporter{}
			h := e.Use(http.HandlerFunc(fn))
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			resp := w.Result()

			r, _ := io.ReadAll(resp.Body)
			assert.Equal(t, tt.want, r)
			assert.Equal(t, tt.wantCode, resp.StatusCode)
		})
	}
}
