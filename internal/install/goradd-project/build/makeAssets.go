package build

//go:generate gofile remove goradd-project/../deploy/assets/*
//go:generate gofile mkdir goradd-project/../deploy/assets/goradd goradd-project/../deploy/assets/project goradd-project/../deploy/assets/bootstrap

//go:generate gofile copy -x scss:less github.com/goradd/goradd/web/assets/* goradd-project/../deploy/assets/goradd
//go:generate gofile copy -x scss:less goradd-project/web/assets/* goradd-project/../deploy/assets/project
//go:generate gofile copy -x scss:less github.com/goradd/goradd/pkg/bootstrap/assets/* goradd-project/../deploy/assets/bootstrap

// Copy your static files
//go:generate gofile copy -x scss:less goradd-project/web/html/* goradd-project/../deploy/html
