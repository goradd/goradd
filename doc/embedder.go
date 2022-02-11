package doc

// This file causes the directory to be served up as a static file directory.
// To make this happen, the main application just needs to import it.

import (
	http2 "github.com/goradd/goradd/pkg/http"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/app"
	"github.com/shurcooL/github_flavored_markdown"
	"io"
	"net/http"
)

func init() {
	// serve up the markdown files in the doc directory
	app.RegisterStaticPath("/goradd/doc", sys.SourceDirectory(), false, nil)
	http2.RegisterFileProcessor(".md", serveMarkdown)
}

// serveMarkdown converts markdown files to html and serves them.
// This would be more efficient if they were preprocessed into html files and served as html,
// but this is an example of how live processing of files can be done.
func serveMarkdown(r io.Reader, w http.ResponseWriter, req *http.Request) error {
	markdown, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	_, err = w.Write(github_flavored_markdown.Markdown(markdown))
	return err
}
