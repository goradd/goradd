//** This file was code generated by GoT. DO NOT EDIT. ***

package welcome

import (
	"context"
	"io"

	"github.com/goradd/goradd/pkg/http"
)

func init() {
	http.RegisterDrawFunc("/goradd/index.html",
		func(ctx context.Context, _w io.Writer) (err error) {

			if _, err = io.WriteString(_w, `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8"/>
<title>Welcome to Goradd</title>
<link href="`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, http.MakeLocalPath("/assets/goradd/welcome/css/welcome.css")); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `" rel="stylesheet">
</head>
<body>
<h1>Welcome to Goradd!</h1>
<p>
Congratulations, you have correctly installed Goradd! Here are some additional things you can
do to learn about goradd and begin building your application.
</p>

<ul>
<li><a href="`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, http.MakeLocalPath("/goradd/configure.html")); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `">Configure your database</a></li>
<li><a href="`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, http.MakeLocalPath("/goradd/tutorial.g")); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `">Walk through the goradd tutorial</a></li>
<li><a href="`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, http.MakeLocalPath("/goradd/examples/controls.g")); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `">View the standard control examples</a></li>
<li><a href="`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, http.MakeLocalPath("/goradd/examples/bootstrap.g")); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `">View the Bootstrap control examples</a></li>
<li><a href="`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, http.MakeLocalPath("/goradd/build.g?cmd=codegen")); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `">Run the code generator</a></li>
<li><a href="`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, http.MakeLocalPath("/goradd/forms/")); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `">View the generated forms</a></li>
<li><a href="`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, http.MakeLocalPath("/goradd/doc/contents.md")); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `">View the documentation</a></li>
<li><a href="`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, http.MakeLocalPath("/goradd/Test.g")); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `
`); err != nil {
				return
			}

			if _, err = io.WriteString(_w, `">Run the browser-based tests</a></li>
</ul>
</body>
</html>
`); err != nil {
				return
			}

			return

		})
}
