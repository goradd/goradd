package app

// Recreate the stage directory.
//go:generate gofile remove goradd-project/deploy/app
//go:generate gofile remove goradd-project/deploy/stage
//go:generate gofile mkdir goradd-project/deploy/app
//go:generate gofile mkdir goradd-project/deploy/stage/root
//go:generate gofile mkdir goradd-project/deploy/stage/assets/goradd goradd-project/deploy/stage/assets/project goradd-project/deploy/stage/assets/bootstrap

// Copy static files
//go:generate gofile copy goradd-project/web/root/* goradd-project/deploy/stage/root

// Copy assets.
//go:generate gofile copy -x scss:less:*.map:README.txt:*.go github.com/goradd/goradd/web/assets/* goradd-project/deploy/stage/assets/goradd
//go:generate gofile copy -x scss:less:*.map:README.txt:*.go goradd-project/web/assets/* goradd-project/deploy/stage/assets/project

// Javascript associated with the messenger service.
//go:generate gofile mkdir goradd-project/deploy/stage/assets/messenger
//go:generate gofile copy -x scss:less:*.map:README.txt:*.go github.com/goradd/goradd/pkg/messageServer/ws/assets/* goradd-project/deploy/stage/assets/messenger

