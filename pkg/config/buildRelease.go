//go:build release
// +build release

package config

import (
	"path"
)

// The Release constant is used throughout the framework to determine if we are running the development version
// or release version of the product. The development version is designed to make on-going development easier,
// and the release version is designed to run on a deployment server.
// It is off by default, but you can turn it on by building with the -tags "release" flag.
// Combine with the nodebug tag like so: go build -tags "release nodebug"
// You might build a release version that keeps the debug features on if you are building for manual testers
const Release = true

// This is the asset directory used as a central repository for the assets. The assets must be copied here as
// part of the deployment process. The variable must be set up as part of the initialization process.

func GoraddAssets() string {
	return path.Join(AssetPrefix, "goradd")
}

func ProjectAssets() string {
	return path.Join(AssetPrefix, "project")
}

func SetProjectDir(path string) {
	panic("do not call SetProjectDir in the Release build")
}

// This is here just to allow things to build, but should not be called
func ProjectDir() string {
	panic("do not call ProjectDir in the Release build")
	return ""
}

// This is here just to allow things to build, but should not be called
func GoraddDir() string {
	panic("do not call GoraddDir in Release build")
	return ""
}
