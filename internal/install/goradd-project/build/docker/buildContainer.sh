#!/bin/sh
pwd
cp Dockerfile ../../../deploy/docker
cp docker-compose.yml ../../../deploy/docker
cd ../../../deploy/docker || exit
docker build -t grapp .
