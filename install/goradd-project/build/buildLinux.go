package build

//go:generate go generate ./makeAssets.go

// This builds for a linux box. Change it depending on your deployment server
//go:generate env GOOS=linux GOARCH=386 go build -tags "release" -ldflags "-s -w" -o ../../../deploy/ -v ../main

//go:generate ./putApp.sh

