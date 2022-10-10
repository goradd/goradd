package welcome

import (
	"embed"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/http"
	"io/fs"
	"path"
)

// This file embeds the static files found here into the application.

var f embed.FS

//go:embed assets/css assets/font
var a embed.FS

func init() {
	fsys := http.FileSystemServer{Fsys: f, SendModTime: true, Hide: []string{".go", ".got", ".tmp"}}
	http.RegisterAppPrefixHandler("/goradd", fsys)

	sub, _ := fs.Sub(a, "assets")
	http.RegisterAssetDirectory(path.Join(config.AssetPrefix, "goradd", "welcome"), sub)

}
