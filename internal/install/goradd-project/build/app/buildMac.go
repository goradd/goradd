package app

//go:generate env GOOS=darwin GOARCH=amd64 go build -tags "release nodebug" -ldflags "-s -w" -o ../../deploy/app/macMain -v ../../main

