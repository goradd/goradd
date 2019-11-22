#!/usr/bin/env bash

minify -o ../../../deploy/server//assets/goradd/js/ ../../../deploy/server//assets/goradd/js/
gzip --best -f -k -r ../../../deploy/server//assets/goradd/js/*

minify -o ../../../deploy/server//assets/goradd/css/ ../../../deploy/server//assets/goradd/css/
gzip --best -f -k -r ../../../deploy/server//assets/goradd/css/*

minify -o ../../../deploy/server//assets/bootstrap/js/ ../../../deploy/server//assets/bootstrap/js/
gzip --best -f -k -r ../../../deploy/server//assets/bootstrap/js/*

minify -o ../../../deploy/server//assets/bootstrap/css/ ../../../deploy/server//assets/bootstrap/css/
gzip --best -f -k -r ../../../deploy/server//assets/bootstrap/css/*

minify -o ../../../deploy/server//assets/project/js/ ../../../deploy/server//assets/project/js/
gzip --best -f -k -r ../../../deploy/server//assets/project/js/*

minify -o ../../../deploy/server//assets/project/css/ ../../../deploy/server//assets/project/css/
gzip --best -f -k -r ../../../deploy/server//assets/project/css/*

# Also zip your static html files
gzip --best -f -k -r ../../../deploy/server//html/*