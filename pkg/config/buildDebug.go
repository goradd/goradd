//go:build !nodebug

package config

// The Debug constant is used throughout the framework to turn on or off various debugging features. It is on by
// default. To turn it off, build with the -tags "nodebug" flag
const Debug = true
