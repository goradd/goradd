// +build !release

package config

import (
	"github.com/goradd/goradd/pkg/sys"
	"path/filepath"
)

func BootstrapAssets() string {
	src := sys.SourceDirectory()
	assetDir := filepath.Join(filepath.Dir(src), "assets")
	return assetDir
}

