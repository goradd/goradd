//** This file was code generated by GoT. DO NOT EDIT. ***

package panels

import (
	"context"
	"io"

	"github.com/goradd/goradd/pkg/page/control"
)

// DrawTemplate draws the content of the matching control's template file.
func (ctrl *Forms1Panel) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

	if _, err = io.WriteString(_w, `<h2>Standard Form Layout</h2>
<p>
This is an example of a very generic form layout in Bootstrap.
</p>
`); err != nil {
		return
	}

	if ctrl.Page().HasControl("nameText-ff") {
		ctrl.Page().GetControl("nameText-ff").(control.LabelAttributer).LabelAttributes().MergeString(`class="form-label"`)
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("nameText-ff").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("nameText-ff").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if ctrl.Page().HasControl("childrenText-ff") {
		ctrl.Page().GetControl("childrenText-ff").(control.LabelAttributer).LabelAttributes().MergeString(`class="form-label"`)
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("childrenText-ff").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("childrenText-ff").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `<div>

`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("status").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("status").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, ` `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `You are: `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ` `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("radioResult").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("radioResult").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("dogCheck-ff").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("dogCheck-ff").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `

`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("ajaxButton").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("ajaxButton").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("serverButton").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("serverButton").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	return
}
