package build

//go:generate gofile remove GOPATH/deploy/assets/*
//go:generate gofile mkdir GOPATH/deploy/assets/goradd GOPATH/deploy/assets/project GOPATH/deploy/assets/bootstrap

//go:generate gofile copy GOPATH/src/github.com/spekary/goradd/assets/* GOPATH/deploy/assets/goradd
//go:generate gofile copy GOPATH/src/goradd-project/assets/* GOPATH/deploy/assets/project
//go:generate gofile copy GOPATH/src/github.com/spekary/goradd/bootstrap/assets/* GOPATH/deploy/assets/bootstrap

