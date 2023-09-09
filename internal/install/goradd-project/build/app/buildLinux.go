package app

//go:generate env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags "release nodebug" -ldflags "-s -w" -o ../../deploy/app/unixMain -v ../../main
