package embed

//go:generate minify -r --sync -o ../../deploy/embed/ ../../deploy/stage/assets
//go:generate minify -r --sync -o ../../deploy/embed/ ../../deploy/stage/root

//go:generate gofile brotli -v -x *.go ../../deploy/embed/
//go:generate gofile gzip -v -d -x *.go:*.br ../../deploy/embed/





