package http

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkDefaultMux(b *testing.B) { benchmarkMuxMatch(b, http.NewServeMux()) }
func BenchmarkGoraddMux(b *testing.B)  { benchmarkMuxMatch(b, NewMux()) }
func benchmarkMuxMatch(b *testing.B, mux Muxer) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	}
	mux.Handle("/", http.HandlerFunc(fn))
	mux.Handle("/index", http.HandlerFunc(fn))
	mux.Handle("/home", http.HandlerFunc(fn))
	mux.Handle("/about", http.HandlerFunc(fn))
	mux.Handle("/contact", http.HandlerFunc(fn))
	mux.Handle("/robots.txt", http.HandlerFunc(fn))
	mux.Handle("/products/", http.HandlerFunc(fn))
	mux.Handle("/products/1", http.HandlerFunc(fn))
	mux.Handle("/products/2", http.HandlerFunc(fn))
	mux.Handle("/products/3", http.HandlerFunc(fn))
	mux.Handle("/products/3/image.jpg", http.HandlerFunc(fn))
	mux.Handle("/admin", http.HandlerFunc(fn))
	mux.Handle("/admin/products/create", http.HandlerFunc(fn))
	mux.Handle("/admin/products/update", http.HandlerFunc(fn))
	mux.Handle("/admin/products/delete", http.HandlerFunc(fn))
	mux.Handle("/admin/products/", http.HandlerFunc(fn))
	mux.Handle("/a/b/", http.HandlerFunc(fn))
	mux.Handle("/a/b", http.HandlerFunc(fn))
	mux.Handle("/a/b/c", http.HandlerFunc(fn))
	mux.Handle("/a/b/c/", http.HandlerFunc(fn))

	paths := []string{"/", "/notfound", "/admin/", "/admin/foo", "/contact", "/products",
		"/products/", "/products/3/image.jpg"}
	var requests []*http.Request
	for _, p := range paths {
		requests = append(requests, httptest.NewRequest("GET", p, nil))
	}
	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if h, p := mux.Handler(requests[i%len(paths)]); h != nil && p == "" {
			b.Error("impossible")
		}
	}
	b.StopTimer()
}

func Test_Mux(t *testing.T) {
	clearGlobals()
	fnFound := func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "Found")
	}

	mux := NewMux()
	mux.Handle("/test", http.HandlerFunc(fnFound))
	mux.Handle("/a", http.HandlerFunc(fnFound))
	mux.HandleFunc("/a/b/", fnFound)

	assert.Panics(t, func() {
		mux.HandleFunc("/a/b/", fnFound)
	}, "Should only register the pattern once")

	assert.Panics(t, func() {
		mux.HandleFunc("", fnFound)
	}, "Should not allow empty pattern")

	assert.Panics(t, func() {
		mux.Handle("/abcd", nil)
	}, "Should not allow nil handler")

	assert.Panics(t, func() {
		mux.HandleFunc("/abcde", nil)
	}, "Should not allow nil handler func")

	h := UseMuxer(mux, http.NotFoundHandler())

	type testT struct {
		path string
		code int
	}
	tests := []testT{
		{"/", 404},
		{"/test", 200},
		{"/test2", 404},
		{"/test/", 404},
		{"/test/test2/", 404},
		{"/test/test2", 404},
		{"/a/b", 301},
		{"/a/b/", 200},
		{"/a/b/c", 200},
		{"//a/b/c", 301},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			resp := w.Result()
			assert.Equal(t, tt.code, resp.StatusCode)
		})
	}

	assert.Panics(t, func() {
		mux.HandleFunc("/a/b/c/d", fnFound)
	}, "Should not allow adding patterns after muxer has been used")

}

func Test_multiLevelMux(t *testing.T) {
	clearGlobals()
	fnFound := func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "Found")
	}
	fnNotFound := func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "Not Found")
	}

	type testT struct {
		path   string
		code   int
		result string
	}
	tests := []testT{
		{"/test", 200, "Found"},
		{"/test2", 200, "Not Found"},
		{"/test/test2/", 200, "Not Found"},
		{"/test/test2/test3", 200, "Found"},
		{"/test/test2/test4", 200, "Not Found"},
	}

	mux := NewMux()
	mux.Handle("/test", http.HandlerFunc(fnFound))

	mux2 := NewMux()
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

func Test_MuxHandlerPattern(t *testing.T) {
	clearGlobals()
	fnFound := func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "Found")
	}

	mux := NewMux()
	mux.Handle("/test", http.HandlerFunc(fnFound))
	mux.Handle("/a", http.HandlerFunc(fnFound))
	mux.Handle("/a/b/", http.HandlerFunc(fnFound))

	UseMuxer(mux, http.NotFoundHandler())

	type testT struct {
		path    string
		pattern string
	}
	tests := []testT{
		{"/", ""},
		{"/test", "/test"},
		{"/test2", ""},
		{"/test/", ""},
		{"/a/b", "/a/b/"},
		{"/a/b/", "/a/b/"},
		{"/a/b/c", "/a/b/"},
		{"//a/b/c", "/a/b/"},
		{"/a/c/../b", "/a/b/"},
		{"/a/c/../b/", "/a/b/"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			_, p := mux.Handler(req)
			assert.Equal(t, tt.pattern, p)
		})
	}

}
