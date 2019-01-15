package build

import (
	"net/http"
	"bytes"
	"strings"
	"html"
	"fmt"
	"github.com/goradd/goradd/ideas"
	"path/filepath"
)

var dependencies = []string {
	"github.com/goradd/got",				// Template processor
	"golang.org/x/tools/cmd/goimports",		// For auto-fixing of import declarations
	"github.com/alexedwards/scs",			// session management
	"github.com/patrickmn/go-cache",		// dependency of scs
	"github.com/gedex/inflector",			// Pluralizing
	"github.com/knq/snaker",				// Snake-case to CamelCase and back
	"github.com/go-sql-driver/mysql",		// Mysql driver
	"github.com/microcosm-cc/bluemonday",	// Filters text input for possible XSS attacks
}

func serveInstall() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		buf := new (bytes.Buffer)
		cmd := r.FormValue("cmd")
		if cmd == "start" {
			go startInstaller()
		}
		drawInstaller(buf)
		w.Write(buf.Bytes())
	}
	return http.HandlerFunc(fn)

}

func startInstaller() {
	results = ""
	stop = false

	// copy
	if isInstalled() {
		results += fmt.Sprintf("Error: %s exists. If you really want to completely reinstall, delete that directory. Otherwise, consider updating.\n", projectPath() )
		stop = true
		return
	} else {
		err := goGet(dependencies...)
		if err != nil {
			stop = true
			return
		}
		results += "Copying goradd-project directory\n"
		err = ideas.DirectoryCopy(filepath.Join(goraddPath(), "buildtools", "install", "goradd-project"), srcPath())
		if err != nil {
			results += fmt.Sprintf("Error copying goradd-project directory: %s", err.Error())
			stop = true
			return
		}

		results += "Copying goradd-tmp directory\n"
		err = ideas.DirectoryCopy(filepath.Join(goraddPath(), "buildtools", "install", "goradd-tmp"), srcPath())
		if err != nil {
			results += fmt.Sprintf("Error copying goradd-tmp directory: %s", err.Error())
			stop = true
			return
		}

		results += "Installing gofile\n"
		cmdResult, errStr, err := executeCmd("go", "install", filepath.Join(srcPath(), "github.com", "goradd", "goradd", "buildtools", "gofile"))
		if err != nil {
			results += errStr
			stop = true
			return
		} else {
			results += cmdResult
		}

	}

	results += "Success!"
	stop = true
}

func goGet(items ...string) (err error) {
	var errStr string
	var cmdResult string
	for _, item := range items {
		results += "go get " + item + "\n"
		cmdResult, errStr, err = executeCmd("go", "get", item)
		if err != nil {
			results += errStr
			return
		} else {
			results += cmdResult
		}
	}
	return
}

func textToHtml(in string) (out string) {
	in = html.EscapeString(in)
	in = strings.Replace(in, "\n\n", "<p>", -1)
	out = strings.Replace(in, "\n", "<br />", -1)
	return
}
