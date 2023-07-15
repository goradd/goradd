package web

// This file embeds the static files into the application as static files during
// development.
//
// For deployment, these files would be copied to the deployment directory, compressed
// and embedded from there.
//
// The root directory represents the root of the website, or the "/" directory.
// You can modify this by using the ProxyPath setting if your application is served
// from a logical subdirectory of a bigger website.
//
// Simply layout all of your static files here and they will be served as if they were part
// of a file system. If the user navigates to a file that does not exist in the file system here,
// the Goradd muxer will be invoked and the rest of your application will serve the resource requested.

import (
	"embed"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/http"
	"io/fs"
	"path"
)

//go:embed root/*
var root embed.FS

//go:embed assets/*
var a embed.FS

func init() {
	// This server is designed to serve HTML type files that can be bookmarked. It servese them out of the
	// root path. Feel free to change the path as needed, or delete this section all together if you are
	// not serving any html file.
	sub, _ := fs.Sub(root, "root")
	serv := http.FileSystemServer{Fsys: sub, SendModTime: true}
	http.RegisterAppPrefixHandler("/", serv)

	// This server serves assets that are not usually bookmarked. It uses a method called
	// cache-busting to make sure that when you deploy new versions of these files, the client
	// will not use a previous cached version, but if the file did not change, the client
	// can still use a cached version.
	// If you have no custom assets, like css, javascript, fonts or images in the project's assets
	// directory, you can delete this section.
	sub, _ = fs.Sub(a, "assets")
	http.RegisterAssetDirectory(path.Join(config.AssetPrefix, "project"), sub)
}
