package assets

// This file embeds the static files found here into the application.

import (
	"embed"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/http"
	"path"
)

//go:embed js
var a embed.FS

func init() {
	http.RegisterAssetDirectory(path.Join(config.AssetPrefix, "/test"), a)
}
