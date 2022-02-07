package docker

// Note that this file gets copied to the deploy directory so it gets run from there

//go:generate docker build -t grapp .
//go:generate docker image prune --force

// Export the container so that it can be copied to another computer
// and imported there with docker load command.

//go:generate docker save -o ./grapp.tar grapp
