package embed

//go:generate env GOOS=darwin GOARCH=amd64 go build -tags "release nodebug" -ldflags "-s -w" -o ../../deploy/embed -v ../../main

