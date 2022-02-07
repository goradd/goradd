package docker

//go:generate gofile remove goradd-project/deploy/docker
//go:generate gofile mkdir goradd-project/deploy/docker

//go:generate gofile copy ../../deploy/app/*Main ../../deploy/docker

//go:generate gofile copy Dockerfile ../../deploy/docker
//go:generate gofile copy buildContainer.go ../../deploy/docker

