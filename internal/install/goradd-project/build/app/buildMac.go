package app

//go:generate env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -tags "release nodebug" -ldflags "-s -w" -o ../../deploy/app/macMain -v ../../main
