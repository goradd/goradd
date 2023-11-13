//go:build !release

package assets

// This file embeds the static files found here into the application during development.
//
// For deployment, these files should be copied to the deployment directory, compressed
// and embedded from there. See goradd-project/build and goradd-project/deploy.

import (
	"embed"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/http"
	"path"
)

//go:embed js
var a embed.FS

func init() {
	http.RegisterAssetDirectory(path.Join(config.AssetPrefix, "messenger"), a)
}
