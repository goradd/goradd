package app

//go:generate env GOOS=windows GOARCH=amd64 go build -tags "release nodebug" -ldflags "-s -w" -o ../../deploy/app/winMain -v ../../main

