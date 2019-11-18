#!/bin/sh

minify -o ../../../deploy/docker/assets/goradd/js/ ../../../deploy/docker/assets/goradd/js/
gzip --best -f -k -r ../../../deploy/docker/assets/goradd/js/*

minify -o ../../../deploy/docker/assets/goradd/css/ ../../../deploy/docker/assets/goradd/css/
gzip --best -f -k -r ../../../deploy/docker/assets/goradd/css/*

minify -o ../../../deploy/docker/assets/bootstrap/js/ ../../../deploy/docker/assets/bootstrap/js/
gzip --best -f -k -r ../../../deploy/docker/assets/bootstrap/js/*

minify -o ../../../deploy/docker/assets/bootstrap/css/ ../../../deploy/docker/assets/bootstrap/css/
gzip --best -f -k -r ../../../deploy/docker/assets/bootstrap/css/*

minify -o ../../../deploy/docker/assets/project/js/ ../../../deploy/docker/assets/project/js/
gzip --best -f -k -r ../../../deploy/docker/assets/project/js/*

minify -o ../../../deploy/docker/assets/project/css/ ../../../deploy/docker/assets/project/css/
gzip --best -f -k -r ../../../deploy/docker/assets/project/css/*

minify -o ../../../deploy/docker/assets/messenger/js/ ../../../deploy/docker/assets/messenger/js/
gzip --best -f -k -r ../../../deploy/docker/assets/messenger/js/*

# Also zip your static html files
gzip --best -f -k -r ../../../deploy/docker/html/*