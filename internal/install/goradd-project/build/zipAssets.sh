#!/usr/bin/env bash

minify -o ../../deploy/assets/goradd/js/ ../../deploy/assets/goradd/js/
gzip --best -f -k -r ../../deploy/assets/goradd/js/*

minify -o ../../deploy/assets/goradd/css/ ../../deploy/assets/goradd/css/
gzip --best -f -k -r ../../deploy/assets/goradd/css/*

minify -o ../../deploy/assets/bootstrap/js/ ../../deploy/assets/bootstrap/js/
gzip --best -f -k -r ../../deploy/assets/bootstrap/js/*

minify -o ../../deploy/assets/bootstrap/css/ ../../deploy/assets/bootstrap/css/
gzip --best -f -k -r ../../deploy/assets/bootstrap/css/*

minify -o ../../deploy/assets/project/js/ ../../deploy/assets/project/js/
gzip --best -f -k -r ../../deploy/assets/project/js/*

minify -o ../../deploy/assets/project/css/ ../../deploy/assets/project/css/
gzip --best -f -k -r ../../deploy/assets/project/css/*

# Also zip your static html files
gzip --best -f -k -r ../../deploy/html/*