//go:build !release
// +build !release

package welcome

import (
	"embed"
	"github.com/goradd/goradd/pkg/http"
)

// This file embeds the static files found here into the application during development.
//
// For deployment, these files should be copied to the deployment directory, compressed
// and embedded from there. See goradd-project/build and goradd-project/deploy.

//go:embed index.html
var f embed.FS

func init() {
	fs := http.FileSystemServer{Fsys: f, SendModTime: true, MustRespond: false, Hide: []string{".go", ".got", ".tmp"}}
	http.RegisterPathHandler("/goradd", fs)
}