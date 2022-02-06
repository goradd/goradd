package welcome

import (
	"embed"
	"github.com/goradd/goradd/pkg/http"
)

// This file embeds the static files found here into the application.

//go:embed index.html
var f embed.FS

func init() {
	fs := http.FileSystemServer{Fsys: f, SendModTime: true, MustRespond: false, Hide: []string{".go", ".got", ".tmp"}}
	http.RegisterAppPrefixHandler("/goradd", fs)
}