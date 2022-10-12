// Package doc contains documentation as Markdown files for the
package doc

import (
	"embed"
	http2 "github.com/goradd/goradd/pkg/http"
	"github.com/shurcooL/github_flavored_markdown"
	"io"
	"net/http"
)

//go:embed *
var docFS embed.FS

func init() {
	fs := http2.FileSystemServer{Fsys: docFS, SendModTime: true}
	http2.RegisterAppPrefixHandler("/goradd/doc", fs)
	http2.RegisterFileProcessor(".md", serveMarkdown)
}

// serveMarkdown converts markdown files to html and serves them.
// This would be more efficient if they were preprocessed into html files and served as html,
// but this is an example of how to process text files.
func serveMarkdown(r io.Reader, w http.ResponseWriter, req *http.Request) error {
	markdown, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	_, err = w.Write(github_flavored_markdown.Markdown(markdown))
	return err
}
