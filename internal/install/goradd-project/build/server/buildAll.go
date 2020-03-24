package server

//go:generate gofile mkdir goradd-project/../deploy
//go:generate gofile remove goradd-project/../deploy/server
//go:generate gofile mkdir goradd-project/../deploy/server
//go:generate go generate ./makeAssets.go
//go:generate go generate ./buildLinux.go
//go:generate go generate ./zipAssets.go

