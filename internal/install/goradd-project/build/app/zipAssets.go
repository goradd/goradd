package app

//go:generate minify -r --sync -o ../../deploy/app/ ../../deploy/stage/assets
//go:generate minify -r --sync -o ../../deploy/app/ ../../deploy/stage/root

//go:generate gofile brotli -v -x *.go ../../deploy/app/
//go:generate gofile gzip -v -d -x *.go:*.br ../../deploy/app/





