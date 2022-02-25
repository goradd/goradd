package http

import (
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileSystemServer_ServeHTTP_Brotli(t *testing.T) {
	fs := os.DirFS("testdata")
	fss := FileSystemServer{Fsys: fs}
	req := httptest.NewRequest("GET", "/test1.txt", nil)
	req.Header.Set("Accept-Encoding", "br")
	w := httptest.NewRecorder()
	fss.ServeHTTP(w, req)
	resp := w.Result()
	assert.EqualValues(t, "br", resp.Header.Get("Content-Encoding"))
	body, _ := io.ReadAll(resp.Body)
	// test for some kind of encoding
	assert.Less(t, 1, len(body))
	assert.NotEqualValues(t, 1, "text")
}

func TestFileSystemServer_ServeHTTP_Gzip(t *testing.T) {
	fs := os.DirFS("testdata")
	fss := FileSystemServer{Fsys: fs}
	req := httptest.NewRequest("GET", "/test1.txt", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	fss.ServeHTTP(w, req)
	resp := w.Result()
	assert.EqualValues(t, "gzip", resp.Header.Get("Content-Encoding"))
	body, _ := io.ReadAll(resp.Body)
	// test for some kind of encoding
	assert.Less(t, 1, len(body))
	assert.NotEqualValues(t, 1, "text")
}

func TestFileSystemServer_ServeHTTP(t *testing.T) {
	fs := os.DirFS("testdata")

	tests := []struct {
		name           string
		sendModTime    bool
		useCacheBuster bool
		hide           []string
		path           string
		encoding       string
		wantContent    string
		wantCode       int
	}{
		{"std", false, false, nil, "/test1.txt", "", "test", 200},
		{"plain", false, false, nil, "/plain/test1.txt", "gzip", "test", 200},
		{"gzip", false, false, nil, "/gzip/test1.txt", "", "test", 200},
		{"brotli", false, false, nil, "/brotli/test1.txt", "", "test", 200},
		{"not found", false, false, nil, "/abc", "", "", 404},
		{"index1", false, false, nil, "/plain", "", "index", 200},
		{"index2", false, false, nil, "/plain/", "", "index", 200},
		{"index3", false, false, nil, "/", "", "index", 200},
		{"cacheBust", false, true, nil, "/plain/gr.abc/test1.txt", "", "test", 200},
		{"hide", false, false, []string{".txt"}, "/plain/test1.txt", "", "", 404},
		{"modTime", true, false, nil, "/plain/test1.txt", "", "test", 200},
		{"badPath", false, false, nil, "/plain/../../test1.txt", "", "", 404},
		{"badGzip", false, false, nil, "/gzip/bad.txt", "", "", 404},
		{"badBrotli", false, false, nil, "/brotli/bad.txt", "", "", 404},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fss := FileSystemServer{
				Fsys:           fs,
				SendModTime:    tt.sendModTime,
				UseCacheBuster: tt.useCacheBuster,
				Hide:           tt.hide,
			}
			req := httptest.NewRequest("GET", tt.path, nil)
			if tt.encoding != "" {
				req.Header.Set("Accept-Encoding", tt.encoding)
			}
			w := httptest.NewRecorder()
			fss.ServeHTTP(w, req)
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if tt.wantContent != "" {
				assert.EqualValues(t, tt.wantContent, body)
			}
			assert.EqualValues(t, tt.wantCode, resp.StatusCode)
			mt := resp.Header.Get("Last-Modified")
			assert.Equal(t, tt.sendModTime, mt != "")
		})
	}
}
