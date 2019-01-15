// +build !release

package config

import (
	"github.com/goradd/goradd/pkg/config"
	"path/filepath"
	"runtime"
)

// This file sets up config variables that are only available during development.

func init() {
	_, filename, _, _ := runtime.Caller(0)

	// The projectDir points to files in the goradd-project directory. The development version would have all of these
	// files moved to a deployment location, so it is not available in the release version of the app. Doing the setup
	// this way ensures that when we build the release version, we will get a compile time failure if we accidentally try
	// to access the projectDir without making sure we are in the dev version of the app.
	projectDir := filepath.Dir(filepath.Dir(filename))
	config.SetProjectDir(projectDir)
}
