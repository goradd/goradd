// +build !release

package config

import (
	"github.com/spekary/goradd/ideas"
	"path/filepath"
)

// The Release constant is used throughout the framework to determine if we are running the development version
// or release version of the product. The development version is designed to make on-going development easier,
// and the release version is designed to run on a deployment server.
// It is off by default, but you can turn it on by building with the -tags "release" flag.
// Combine with the nodebug tag like so: go build -tags "release nodebug"
// You might build a release version that keeps the debug features on if you are building for manual testers
const Release = false

// These directories are available during development, but not for the release build. If you have static files you
// need to locate, you will need to provide a different mechanism to do that. See the main package for how the
// framework does that for assets, by using a pattern in the URL, combined with a flag past in to the application at
// runtime.
var projectDir string
var goraddDir string // filled in by Goradd
var goraddAssets string
var projectAssets string

// GoraddAssets returns the disk location of the github.com/spekary/goradd/assets directory.
func GoraddAssets() string {
	return goraddAssets
}

func ProjectAssets() string {
	return projectAssets
}

func GoraddDir() string {
	return goraddDir
}

func ProjectDir() string {
	return projectDir
}

func init() {
	gopath := ideas.GoPath()
	goraddDir = filepath.Join(gopath, "src","github.com","spekary","goradd")
	projectDir = filepath.Join(gopath, "src","goradd-project")
	goraddAssets = filepath.Join(goraddDir, "assets")
	projectAssets = filepath.Join(projectDir, "assets")
}
