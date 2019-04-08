package controls

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/web/examples/controls/panels"
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

type createFunction func(ctx context.Context, parent page.ControlI)
var controls = []struct {
		key string
		name string
		f createFunction
	}{
	{"", "Home", panels.NewDefaultPanel},
	{"textbox", "Textboxes", panels.NewTextboxPanel},
	{"checkbox", "Checkboxes and Radio Buttons", panels.NewCheckboxPanel},
	{"selectlist", "Selection Lists", panels.NewSelectListPanel},
	{"table", "Tables", panels.NewTablePanel},
	{"tabledb", "Tables - Checkbox Column", panels.NewTableCheckboxPanel},
	{"tablecheckbox", "Tables - Database Columns", panels.NewTableDbPanel},
	{"hlist", "Nested Lists", panels.NewHListPanel},
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
		case "tableproxy":
			panels.NewTableProxyPanel(ctx, f.detail)
		case "hlist":
			panels.NewHListPanel(ctx, f.detail)

			// TODO: TableSelect, TableSort, TableStyler
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

