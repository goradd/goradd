//** This file was code generated by GoT. DO NOT EDIT. ***

package panels

import (
	"context"
	"io"
)

// DrawTemplate draws the content of the matching control's template file.
func (ctrl *FileSelectPanel) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

	if _, err = io.WriteString(_w, `<h2>File Select Button</h2>
<p>
File Upload Buttons let the user select and upload files.
</p>
<p>
File uploading can present security risks to your website. You should always check your files for integrity to
prevent malicious activity.
</p>
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("uploadButton").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("uploadButton").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `
`); err != nil {
		return
	}

	return
}
