//** This file was code generated by GoT. DO NOT EDIT. ***

package tutorial

import (
	"context"
	"io"

	"github.com/goradd/goradd/pkg/http"
	"github.com/goradd/goradd/pkg/page"
)

func (ctrl *IndexForm) AddHeadTags() {
	ctrl.FormBase.AddHeadTags()
	if "Tutorial" != "" {
		ctrl.Page().SetTitle("Tutorial")
	}

	// deal with body attributes too
	ctrl.Page().BodyAttributes = ``
}

func (ctrl *IndexForm) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

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
`); err != nil {
		return
	}
	path := page.GetContext(ctx).HttpContext.URL.Path
	if _, err = io.WriteString(_w, `<div id="sidebar" class="open">
    <a href="javascript:void(0)" id="togglebtn" onclick="toggleSidebar();"><span id="isopen">&larrb;</span><span id="isclosed">&rarrb;</span></a>
    <div id="sidebar_content">
        <h2><a href="`); err != nil {
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

	if _, err = io.WriteString(_w, `">Home</a></h2>
        <h2>ORM</h2>
          <ul>
        `); err != nil {
		return
	}

	for _, pr := range pages["orm"] {

		if _, err = io.WriteString(_w, `
            <li><a href="`); err != nil {
			return
		}

		if _, err = io.WriteString(_w, http.MakeLocalPath(path+`?pageID=orm-`+pr.id)); err != nil {
			return
		}

		if _, err = io.WriteString(_w, `">`); err != nil {
			return
		}

		if _, err = io.WriteString(_w, pr.title); err != nil {
			return
		}

		if _, err = io.WriteString(_w, `</a></li>
        `); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `
          </ul>
          <h2>Controls</h2>
           <ul>
             `); err != nil {
		return
	}

	for _, pr := range pages["controls"] {

		if _, err = io.WriteString(_w, `
                 <li><a href="`); err != nil {
			return
		}

		if _, err = io.WriteString(_w, http.MakeLocalPath(path+`?pageID=controls-`+pr.id)); err != nil {
			return
		}

		if _, err = io.WriteString(_w, `">`); err != nil {
			return
		}

		if _, err = io.WriteString(_w, pr.title); err != nil {
			return
		}

		if _, err = io.WriteString(_w, `</a></li>
             `); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `
           </ul>

  </div>
</div>
<div id="content" class="open">
<h1>Tutorial</h1>
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("viewSourceButton").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("viewSourceButton").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `<div id="detail_container">
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

	if _, err = io.WriteString(_w, `</div>
</div>

`); err != nil {
		return
	}

	return
}
