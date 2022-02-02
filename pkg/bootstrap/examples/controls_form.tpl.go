//** This file was code generated by GoT. DO NOT EDIT. ***

package examples

import (
	"context"
	"io"
)

func (ctrl *ControlsForm) AddHeadTags() {
	ctrl.FormBase.AddHeadTags()
	if "Goradd Bootstrap Examples" != "" {
		ctrl.Page().SetTitle("Goradd Bootstrap Examples")
	}

	// deal with body attributes too
	ctrl.Page().BodyAttributes = `
`
}

func (ctrl *ControlsForm) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}

		if err = ctrl.Page().GetControl("nav").Draw(ctx, _w); err != nil {
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

		if err = ctrl.Page().GetControl("nav").ProcessAttributeString(``).Draw(ctx, _w); err != nil {
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

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}

		if err = ctrl.Page().GetControl("detailPanel").Draw(ctx, _w); err != nil {
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

		if err = ctrl.Page().GetControl("detailPanel").ProcessAttributeString(``).Draw(ctx, _w); err != nil {
			return
		}

		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	return
}
