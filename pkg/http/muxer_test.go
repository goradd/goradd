package http

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func clearGlobals() {
	patternHandlers = make(handlerMap)
	appHandlers = make(handlerMap)
	PatternMuxer  = nil
	AppMuxer = nil
}
func TestRegisterAppHandler(t *testing.T) {
	clearGlobals()
	type args struct {
		pattern string
		handler http.Handler
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestRegisterAppPathHandler(t *testing.T) {
	clearGlobals()
	type args struct {
		prefix  string
		handler http.Handler
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	clearGlobals()
	type args struct {
		pattern string
		handler http.Handler
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestRegisterPathHandler(t *testing.T) {
	clearGlobals()
	type args struct {
		prefix  string
		handler http.Handler
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestUseAppMuxer(t *testing.T) {
	clearGlobals()
	type args struct {
		mux  Muxer
		next http.Handler
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UseAppMuxer(tt.args.mux, tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UseAppMuxer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsePatternMuxer(t *testing.T) {
	clearGlobals()
	type args struct {
		mux  Muxer
		next http.Handler
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UsePatternMuxer(tt.args.mux, tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UsePatternMuxer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_registerHandler(t *testing.T) {
	clearGlobals()
	type args struct {
		pattern string
		handler http.Handler
		m       handlerMap
		mux     Muxer
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_registerPrefixHandler(t *testing.T) {
	clearGlobals()
	type args struct {
		path string
	}
	type wantT struct {
		code int
		result string
	}
	tests := []struct {
		name string
		path string
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
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
	RegisterPathHandler("/test/test2/", mux2)
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
	RegisterAppPathHandler("/test/test2/", mux2)
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

	RegisterPathHandler("", http.HandlerFunc(fnFoundRoot))
	assert.Panics(t, func() {
		RegisterPathHandler("/",  http.HandlerFunc(fnFoundRoot)) //
	},
	"Blank path and root path should be equal",
	)
	RegisterPathHandler("test", http.HandlerFunc(fnFound))

	mux2 := http.NewServeMux()
	mux2.Handle("/test3", http.HandlerFunc(fnFound))
	RegisterPathHandler("/test/test2/", mux2)
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