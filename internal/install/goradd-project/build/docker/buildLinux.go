package docker

//go:generate env GOOS=linux GOARCH=amd64 go build -tags "release nodebug" -ldflags "-s -w" -o ../../../deploy/docker/grapp -v ../../main

