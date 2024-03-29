//** This file was code generated by GoT. DO NOT EDIT. ***

package tutorial

import (
	"context"
	"html"
	"io"
	"io/ioutil"
)

func (ctrl *FilePanel) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

	if _, err = io.WriteString(_w, `<figure>
<figcaption>`); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ctrl.Base); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `</figcaption>
<pre>
<code>
`); err != nil {
		return
	}

	var content string

	if ctrl.File != "" {
		data, err := ioutil.ReadFile(ctrl.File)
		if err == nil {
			content = html.EscapeString((string(data)))
		}
	}

	if _, err = io.WriteString(_w, content); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `
</code>
</pre>
</figure>
`); err != nil {
		return
	}

	return
}
