package docker

//go:generate gofile remove goradd-project/../deploy/docker/assets
//go:generate gofile mkdir goradd-project/../deploy/docker/assets/goradd goradd-project/../deploy/docker/assets/project goradd-project/../deploy/docker/assets/bootstrap

//go:generate gofile copy -x scss:less github.com/goradd/goradd/web/assets/* goradd-project/../deploy/docker/assets/goradd
//go:generate gofile copy -x scss:less goradd-project/web/assets/* goradd-project/../deploy/docker/assets/project
//go:generate gofile copy -x scss:less github.com/goradd/goradd/pkg/bootstrap/assets/* goradd-project/../deploy/docker/assets/bootstrap

// Javascript associated with the messenger service. Change this to copy the support files for the messenger you choose.
//go:generate gofile mkdir goradd-project/../deploy/docker/assets/messenger
//go:generate gofile copy -x scss:less github.com/goradd/goradd/pkg/messageServer/ws/assets/* goradd-project/../deploy/docker/assets/messenger

// Copy your static files
//go:generate gofile remove goradd-project/../deploy/docker/html
//go:generate gofile copy goradd-project/web/html goradd-project/../deploy/docker/


