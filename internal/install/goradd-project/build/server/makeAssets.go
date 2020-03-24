package server

//go:generate gofile remove goradd-project/../deploy/server/assets/*
//go:generate gofile mkdir goradd-project/../deploy/server/assets/goradd goradd-project/../deploy/server/assets/project goradd-project/../deploy/server/assets/bootstrap

//go:generate gofile copy -x scss:less github.com/goradd/goradd/web/assets/* goradd-project/../deploy/server/assets/goradd
//go:generate gofile copy -x scss:less goradd-project/web/assets/* goradd-project/../deploy/server/assets/project
//go:generate gofile copy -x scss:less github.com/goradd/goradd/pkg/bootstrap/assets/* goradd-project/../deploy/server/assets/bootstrap

// Javascript associated with the messenger service. Change this to copy the support files for the messenger you choose.
//go:generate gofile mkdir goradd-project/../deploy/server/assets/messenger
//go:generate gofile copy -x scss:less github.com/goradd/goradd/pkg/messageServer/ws/assets/* goradd-project/../deploy/server/assets/messenger

// Copy your static files
//go:generate gofile mkdir goradd-project/../deploy/server/html
//go:generate gofile copy -x scss:less goradd-project/web/html/* goradd-project/../deploy/server/html/

