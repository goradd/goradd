// +build release

package config

import (
	"github.com/goradd/goradd/pkg/config"
	"path"
)

func BootstrapAssets() string {
	return path.Join(config.AssetPrefix, "bootstrap")
}

