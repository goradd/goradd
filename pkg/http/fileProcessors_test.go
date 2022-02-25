package http

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/shurcooL/github_flavored_markdown"
	"github.com/stretchr/testify/assert"
)

// serveMarkdown converts markdown files to html and serves them.
// This would be more efficient if they were preprocessed into html files and served as html,
// but this is an example of how live processing of files can be done.
func markdownProcessor(r io.Reader, w http.ResponseWriter, req *http.Request) error {
	markdown, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	_, err = w.Write(github_flavored_markdown.Markdown(markdown))
	return err
}

func TestRegisterFileProcessor(t *testing.T) {
	fs := os.DirFS("testdata")
	RegisterFileProcessor(".md", markdownProcessor)
	fss := FileSystemServer{Fsys: fs}
	req := httptest.NewRequest("GET", "/test2.md", nil)
	w := httptest.NewRecorder()
	fss.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	assert.True(t, bytes.Contains(body, []byte("Title")))
}

func TestRegisterFileProcessorNotFound(t *testing.T) {
	fs := os.DirFS("testdata")
	RegisterFileProcessor(".md", markdownProcessor)
	fss := FileSystemServer{Fsys: fs}
	req := httptest.NewRequest("GET", "/bad.md", nil)
	w := httptest.NewRecorder()
	fss.ServeHTTP(w, req)
	resp := w.Result()
	assert.Equal(t, 404, resp.StatusCode)
}
