//** This file was code generated by got. ***

package page

import (
	"bytes"
	"context"

	"github.com/spekary/goradd/html"
)

func LabelTmpl(ctx context.Context, w LabelWrapperType, ctrl ControlI, h string, buf *bytes.Buffer) {
	labelAttr := w.LabelAttributes().String()

	buf.WriteString(`<div id="`)

	buf.WriteString(ctrl.ID())

	buf.WriteString(`_ctl" `)

	buf.WriteString(ctrl.WrapperAttributes().String())

	buf.WriteString(` >
`)
	if ctrl.Label() != "" {
		buf.WriteString(`    `)
		if ctrl.TextIsLabel() {
			buf.WriteString(`  <span id="`)

			buf.WriteString(ctrl.ID())

			buf.WriteString(`_lbl" class="goradd-lbl" `)

			buf.WriteString(labelAttr)

			buf.WriteString(`>`)

			buf.WriteString(ctrl.Label())

			buf.WriteString(`</span>
    `)
		} else {

			buf.WriteString(`  <label id="`)

			buf.WriteString(ctrl.ID())

			buf.WriteString(`_lbl" class="goradd-lbl"`)
			if ctrl.HasFor() {
				buf.WriteString(` for="`)

				buf.WriteString(ctrl.ID())

				buf.WriteString(`" `)
			}

			buf.WriteString(` `)

			buf.WriteString(labelAttr)

			buf.WriteString(`>`)

			buf.WriteString(ctrl.Label())

			buf.WriteString(`</label>
    `)
		}
	}

	buf.WriteString(html.Indent(h))

	buf.WriteString(`
  <div id="`)

	buf.WriteString(ctrl.ID())

	buf.WriteString(`_err" class="goradd-error">`)

	buf.WriteString(ctrl.ValidationMessage())

	buf.WriteString(`</div>
`)
	if ctrl.Instructions() != "" {
		buf.WriteString(`  <div id="`)

		buf.WriteString(ctrl.ID())

		buf.WriteString(`_inst" class="goradd-instructions" >`)

		buf.WriteString(ctrl.Instructions())

		buf.WriteString(`</div>
`)
	}

	buf.WriteString(`</div>
`)

	return

}
