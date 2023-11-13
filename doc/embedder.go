// Package doc contains documentation as Markdown files.
package doc

import (
	"embed"
	http2 "github.com/goradd/goradd/pkg/http"
	"github.com/russross/blackfriday/v2"
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
	output := blackfriday.Run(markdown)
	_, err = w.Write(output)
	return err
}
