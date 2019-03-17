package controls

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/examples/controls/panels"
)

const ControlsFormPath = "/goradd/examples/controls.g"
const ControlsFormId = "ControlsForm"

const (
	TestButtonAction = iota + 1
)

type ControlsForm struct {
	FormBase
	detail 		  *Panel
}

func NewControlsForm(ctx context.Context) page.FormI {
	f := &ControlsForm{}
	f.Init(ctx, f, ControlsFormPath, ControlsFormId)
	f.AddRelatedFiles()

	f.detail = NewPanel(f, "detailPanel")

	return f
}

func (f *ControlsForm) LoadControls(ctx context.Context) {
	if id, ok := page.GetContext(ctx).FormValue("control"); ok {
		switch id {
		case "textbox":
			panels.NewTextboxPanel(ctx, f.detail)
		case "checkbox":
			panels.NewCheckboxPanel(ctx, f.detail)
		case "selectlist":
			panels.NewSelectListPanel(ctx, f.detail)
		case "table":
			panels.NewTablePanel(ctx, f.detail)
		case "tabledb":
			panels.NewTableDbPanel(ctx, f.detail)
		case "tablecheckbox":
			panels.NewTableCheckboxPanel(ctx, f.detail)
		case "hlist":
			panels.NewHListPanel(ctx, f.detail)
		default:
			panels.NewDefaultPanel(ctx, f.detail)
		}
	} else {
		panels.NewDefaultPanel(ctx, f.detail)
	}
}

func init() {
	page.RegisterPage(ControlsFormPath, NewControlsForm, ControlsFormId)
}

