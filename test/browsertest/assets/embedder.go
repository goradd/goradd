//go:build !release
// +build !release

package assets

// This file embeds the static files found here into the application during development.

import (
	"embed"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/http"
	"path"
)

//go:embed js
var a embed.FS

func init() {
	// The path below is the same path the assets should be copied to for deployment.
	http.RegisterAssetDirectory(path.Join(config.AssetPrefix, "/test"), a)
}
