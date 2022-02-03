//go:build release
// +build release

package web

// This file embeds the static files into the application as static files for
// the release build.
//
import (
	"embed"
	"github.com/goradd/goradd/pkg/http"
	"io/fs"
	"github.com/goradd/goradd/pkg/config"
)

//go:embed embed/root/*
var root embed.FS

//go:embed embed/assets/*
var a embed.FS

func init() {
	// This server is designed to serve HTML type files that can be bookmarked.
	sub, _ := fs.Sub(root, "embed")
	sub, _ = fs.Sub(sub, "root")

	serv := http.FileSystemServer{Fsys: sub, SendModTime: true, MustRespond: false}
	http.RegisterPrefixHandler("/", serv)

	// This server serves assets that are not usually bookmarked. It uses a method called
	// cache-busting to make sure that when you deploy new versions of these files, the client
	// will not use a previous cached version, but if the file did not change, the client
	// can still use a cached version.
	sub,_ = fs.Sub(a, "embed")
	sub,_ = fs.Sub(sub, "assets")
	http.RegisterAssetDirectory(config.AssetPrefix, sub)
}
