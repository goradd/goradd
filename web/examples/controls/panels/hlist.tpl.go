//** This file was code generated by GoT. DO NOT EDIT. ***

package panels

import (
	"context"
	"io"
)

func (ctrl *HListPanel) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

	if _, err = io.WriteString(_w, `<h2>Hierarchical Lists</h2>
<p>
Hierarchical lists dynamically create standard html `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `&lt;ol&gt; and &lt;ul&gt; `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ` type of lists. These lists normally are not
interactive, but many javascript/css libraries use these structures to create interactive dropdown menus,
hierarchical checklists, and more.
</p>
<h3>Ordered List</h3>

`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}

		if err = ctrl.Page().GetControl("orderedList").Draw(ctx, _w); err != nil {
			return
		}

		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}

		if err = ctrl.Page().GetControl("orderedList").ProcessAttributeString(``).Draw(ctx, _w); err != nil {
			return
		}

		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `
<h3>Unordered List</h3>
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}

		if err = ctrl.Page().GetControl("unorderedList").Draw(ctx, _w); err != nil {
			return
		}

		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}

		if err = ctrl.Page().GetControl("unorderedList").ProcessAttributeString(``).Draw(ctx, _w); err != nil {
			return
		}

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
