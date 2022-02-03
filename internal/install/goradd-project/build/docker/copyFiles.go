package docker

//go:generate gofile remove goradd-project/deploy/docker
//go:generate gofile mkdir goradd-project/deploy/docker

//go:generate gofile copy ../../deploy/app/macMain ../../deploy/docker

//go:generate gofile copy Dockerfile ../../deploy/docker
//go:generate gofile copy docker-compose.yml ../../deploy/docker
//go:generate gofile copy db.cfg ../../deploy/docker
//go:generate gofile copy buildContainer.go ../../deploy/docker

