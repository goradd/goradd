package examples

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"

	bootstrap "github.com/goradd/goradd/pkg/bootstrap/control"
)

const COMPONENTS_PATH = "/bootstrap/components"
const COMPONENTS_ID = "ComponentsForm"

const ()

type ComponentsForm struct {
	control.FormBase
	ButtonPanel *control.Fieldset
}

func NewComponentsForm(ctx context.Context) page.FormI {
	f := &ComponentsForm{}
	f.Init(ctx, f, COMPONENTS_PATH, COMPONENTS_ID)
	f.CreateControls(ctx)
	return f
}

func (f *ComponentsForm) CreateControls(ctx context.Context) {
	f.ButtonPanel = control.NewFieldset(f, "")
	bootstrap.NewButton(f.ButtonPanel, "").SetText("Button1")
}

func (f *ComponentsForm) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	}
}

func init() {
	page.RegisterPage(COMPONENTS_PATH, NewComponentsForm, COMPONENTS_ID)
}
