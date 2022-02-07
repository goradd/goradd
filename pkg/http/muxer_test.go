package http

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func clearGlobals() {
	patternHandlers = make(handlerMap)
	appHandlers = make(handlerMap)
	PatternMuxer  = nil
	AppMuxer = nil
}

func Test_serveMuxer(t *testing.T) {
	clearGlobals()
	fnFound := func(w http.ResponseWriter, r *http.Request) {
		_,_ = io.WriteString(w, "Found")
	}
	fnNotFound := func(w http.ResponseWriter, r *http.Request) {
		_,_ = io.WriteString(w, "Not Found")
	}

	type testT struct {
		path  string
		code  int
		result string
	}
	tests := []testT {
		testT{"/test", 200, "Found"},
		testT{"/test2", 200, "Not Found"},
		testT{"/test/test2/", 200, "Not Found"},
		testT{"/test/test2/test3", 200, "Found"},
		testT{"/test/test2/test4", 200, "Not Found"},
	}

	mux := http.NewServeMux()
	mux.Handle("/test", http.HandlerFunc(fnFound))

	mux2 := http.NewServeMux()
	mux2.Handle("/test3", http.HandlerFunc(fnFound))
	mux.Handle("/test/test2/", http.StripPrefix("/test/test2", UseMuxer(mux2, http.HandlerFunc(fnNotFound))))

	h := UseMuxer(mux, http.HandlerFunc(fnNotFound))

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			s := string(body)

			assert.Equal(t, tt.code, resp.StatusCode)
			assert.Equal(t, tt.result, s)
		})
	}

}

func Test_PatternRegistrations(t *testing.T) {
	clearGlobals()
	fnFound := func(w http.ResponseWriter, r *http.Request) {
		_,_ = io.WriteString(w, "Found")
	}
	fnNotFound := func(w http.ResponseWriter, r *http.Request) {
		_,_ = io.WriteString(w, "Not Found")
	}

	type testT struct {
		path  string
		code  int
		result string
	}

	RegisterHandler("/test", http.HandlerFunc(fnFound))
	mux2 := http.NewServeMux()
	mux2.Handle("/test3", http.HandlerFunc(fnFound))
	RegisterPrefixHandler("/test/test2/", mux2)
	mux := http.NewServeMux(); // test using mux in the middle of registration process
	h := UsePatternMuxer(mux,http.HandlerFunc(fnNotFound ))
	RegisterHandler("/test4", http.HandlerFunc(fnFound))
	assert.Panics(t, func() {
		h = UsePatternMuxer(mux,http.HandlerFunc(fnNotFound ))
	}, "Can only use the muxer once")

	tests := []testT {
		testT{"/test", 200, "Found"},
		testT{"/test4", 200, "Found"},
		testT{"/test/test2/", 404, "404 page not found\n"},
		testT{"/test/test2/test3", 200, "Found"},
		testT{"/test/test2/test4", 404, "404 page not found\n"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			s := string(body)

			assert.Equal(t, tt.code, resp.StatusCode)
			assert.Equal(t, tt.result, s)
		})
	}

}

func Test_AppRegistrations(t *testing.T) {
	clearGlobals()
	fnFound := func(w http.ResponseWriter, r *http.Request) {
		_,_ = io.WriteString(w, "Found")
	}
	fnNotFound := func(w http.ResponseWriter, r *http.Request) {
		_,_ = io.WriteString(w, "Not Found")
	}

	type testT struct {
		path  string
		code  int
		result string
	}

	RegisterAppHandler("/test", http.HandlerFunc(fnFound))
	mux2 := http.NewServeMux()
	mux2.Handle("/test3", http.HandlerFunc(fnFound))
	RegisterAppPrefixHandler("/test/test2/", mux2)
	mux := http.NewServeMux(); // test using mux in the middle of registration process
	h := UseAppMuxer(mux,http.HandlerFunc(fnNotFound ))
	RegisterAppHandler("/test4", http.HandlerFunc(fnFound))
	assert.Panics(t, func() {
		h = UseAppMuxer(mux,http.HandlerFunc(fnNotFound ))
	}, "Can only use the muxer once")

	tests := []testT {
		testT{"/test", 200, "Found"},
		testT{"/test4", 200, "Found"},
		testT{"/test/test2/", 404, "404 page not found\n"},
		testT{"/test/test2/test3", 200, "Found"},
		testT{"/test/test2/test4", 404, "404 page not found\n"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			s := string(body)

			assert.Equal(t, tt.code, resp.StatusCode)
			assert.Equal(t, tt.result, s)
		})
	}

}

func Test_PathRegistrations(t *testing.T) {
	clearGlobals()
	fnFound := func(w http.ResponseWriter, r *http.Request) {
		_,_ = io.WriteString(w, "Found")
	}
	fnFoundRoot := func(w http.ResponseWriter, r *http.Request) {
		_,_ = io.WriteString(w, "Found Root")
	}
	fnNotFound := func(w http.ResponseWriter, r *http.Request) {
		_,_ = io.WriteString(w, "Not Found")
	}

	type testT struct {
		path  string
		code  int
		result string
	}

	RegisterPrefixHandler("", http.HandlerFunc(fnFoundRoot))
	assert.Panics(t, func() {
		RegisterPrefixHandler("/",  http.HandlerFunc(fnFoundRoot)) //
	},
	"Blank path and root path should be equal",
	)
	RegisterPrefixHandler("test", http.HandlerFunc(fnFound))

	mux2 := http.NewServeMux()
	mux2.Handle("/test3", http.HandlerFunc(fnFound))
	RegisterPrefixHandler("/test/test2/", mux2)
	mux := http.NewServeMux(); // test using mux in the middle of registration process
	h := UsePatternMuxer(mux,http.HandlerFunc(fnNotFound ))

	tests := []testT {
		testT{"/test4", 200, "Found Root"},
		testT{"/test/", 200, "Found"},
		testT{"/test/test2/", 404, "404 page not found\n"},
		testT{"/test/test2/test3", 200, "Found"},
		testT{"/test/test2/test4", 404, "404 page not found\n"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			s := string(body)

			assert.Equal(t, tt.code, resp.StatusCode)
			assert.Equal(t, tt.result, s)
		})
	}

}