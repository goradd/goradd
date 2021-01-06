#!/bin/sh
pwd
cp Dockerfile ../../../deploy/docker
cp docker-compose.yml ../../../deploy/docker
cp db.cfg ../../../deploy/docker
cd ../../../deploy/docker || exit
docker build -t grapp .
docker image prune --force

# export the container so that it can be copied to another computer
# and imported there with docker load command
docker save -o ./grapp.tar grapp