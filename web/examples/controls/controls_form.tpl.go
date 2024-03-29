//** This file was code generated by GoT. DO NOT EDIT. ***

package controls

import (
	"context"
	"io"
)

// AddHeadTags adds items that will appear in the head tag of the html page.
func (ctrl *ControlsForm) AddHeadTags() {
	ctrl.FormBase.AddHeadTags()
	if "Control Examples" != "" {
		ctrl.Page().SetTitle("Control Examples")
	}

	// deal with body attributes too
	ctrl.Page().BodyAttributes = ``
}

// DrawTemplate draws the content of the matching form's template file.
func (ctrl *ControlsForm) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

	if _, err = io.WriteString(_w, `<script>
function toggleSidebar() {
    g$('sidebar').toggleClass('open');
    g$('content').toggleClass('open');
}
</script>
`); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `
<div id="sidebar" class="open">
    <a href="javascript:void(0)" id="togglebtn" onclick="toggleSidebar();"><span id="isopen">&larrb;</span><span id="isclosed">&rarrb;</span></a>
    <div id="sidebar_content">
        <div class="controlList_scroll">
            `); err != nil {
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
		ctrl.Page().GetControl("listPanel").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("listPanel").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `        </div>
    </div>
</div>
<div id="content" class="open">
    <h1>Control Examples</h1>

    <div class="detail_container">
        `); err != nil {
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
		ctrl.Page().GetControl("detailPanel").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("detailPanel").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `    </div>
</div>
`); err != nil {
		return
	}

	return
}
