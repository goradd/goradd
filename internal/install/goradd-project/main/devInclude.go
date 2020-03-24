// +build !release

package main

import (
	_ "github.com/goradd/goradd/pkg/bootstrap/examples" // Bootstrap examples
	"github.com/goradd/goradd/pkg/config"
	_ "github.com/goradd/goradd/test/browsertest"
	"github.com/goradd/goradd/web/app"
	_ "github.com/goradd/goradd/web/examples"
	_ "github.com/goradd/goradd/web/welcome"
	"github.com/shurcooL/github_flavored_markdown"
	_ "goradd-project/gen" // Code-generated forms
	"io/ioutil"
	"net/http"
	"path/filepath"
)

// This file conditionally builds examples, generated forms, and other things
// useful during development, but that you definitely do not want in a release of your program.

func init() {
	// serve up the markdown files in the doc directory
	app.RegisterStaticPath("/goradd/doc", filepath.Join(config.GoraddDir(), "/doc"))
	app.RegisterStaticFileProcessor(".md", serveMarkdown)
}

func serveMarkdown(file string, w http.ResponseWriter, r *http.Request) {
	markdown, _ := ioutil.ReadFile(file)
	_, _ = w.Write(github_flavored_markdown.Markdown(markdown))
}
