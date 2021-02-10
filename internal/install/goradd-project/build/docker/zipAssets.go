package docker

//go:generate minify -r -o ../../../deploy/docker/assets/goradd/js/ ../../../deploy/docker/assets/goradd/js/
//go:generate gzip --best -f -k -r ../../../deploy/docker/assets/goradd/js/

//go:generate minify -r -o ../../../deploy/docker/assets/goradd/css/ ../../../deploy/docker/assets/goradd/css/
//go:generate gzip --best -f -k -r ../../../deploy/docker/assets/goradd/css/

//go:generate minify -r -o ../../../deploy/docker/assets/bootstrap/js/ ../../../deploy/docker/assets/bootstrap/js/
//go:generate gzip --best -f -k -r ../../../deploy/docker/assets/bootstrap/js/

//go:generate minify -r -o ../../../deploy/docker/assets/bootstrap/css/ ../../../deploy/docker/assets/bootstrap/css/
//go:generate gzip --best -f -k -r ../../../deploy/docker/assets/bootstrap/css/

//go:generate minify -r -o ../../../deploy/docker/assets/project/js/ ../../../deploy/docker/assets/project/js/
//go:generate gzip --best -f -k -r ../../../deploy/docker/assets/project/js/

//go:generate minify -r -o ../../../deploy/docker/assets/project/css/ ../../../deploy/docker/assets/project/css/
//go:generate gzip --best -f -k -r ../../../deploy/docker/assets/project/css/


