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
			panels.NewTextboxPanel(f.detail, "")
		default:
			panels.NewDefaultPanel(f.detail, "")
		}
	} else {
		panels.NewDefaultPanel(f.detail, "")
	}
}

func init() {
	page.RegisterPage(ControlsFormPath, NewControlsForm, ControlsFormId)
}

