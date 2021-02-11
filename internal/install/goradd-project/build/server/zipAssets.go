package server

//go:generate minify -r -o ../../../deploy/server/assets/goradd/js/ ../../../deploy/server/assets/goradd/js/
//go:generate gzip --best -f -k -r ../../../deploy/server/assets/goradd/js/

//go:generate minify -r -o ../../../deploy/server/assets/goradd/css/ ../../../deploy/server/assets/goradd/css/
//go:generate gzip --best -f -k -r ../../../deploy/server/assets/goradd/css/

//go:generate minify -r -o ../../../deploy/server/assets/bootstrap/js/ ../../../deploy/server/assets/bootstrap/js/
//go:generate gzip --best -f -k -r ../../../deploy/server/assets/bootstrap/js/

//go:generate minify -r -o ../../../deploy/server/assets/bootstrap/css/ ../../../deploy/server/assets/bootstrap/css/
//go:generate gzip --best -f -k -r ../../../deploy/server/assets/bootstrap/css/

//go:generate minify -r -o ../../../deploy/server/assets/project/js/ ../../../deploy/server/assets/project/js/
//go:generate gzip --best -f -k -r ../../../deploy/server/assets/project/js/

//go:generate minify -r -o ../../../deploy/server/assets/project/css/ ../../../deploy/server/assets/project/css/
//go:generate gzip --best -f -k -r ../../../deploy/server/assets/project/css/

//go:generate gzip --best -f -k -r ../../../deploy/server/html/

