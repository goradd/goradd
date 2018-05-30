package examples

import (
	"context"
	"github.com/spekary/goradd/page"
	localPage "goradd/page"
	grctl "github.com/spekary/goradd/page/control"

	bootstrap "github.com/spekary/goradd/bootstrap/control"
)


const COMPONENTS_PATH = "/bootstrap/components"
const COMPONENTS_ID = "ComponentsForm"

const (
)

type ComponentsForm struct {
	localPage.FormBase
	ButtonPanel	*grctl.Fieldset
}

func NewComponentsForm(ctx context.Context) page.FormI {
	f := &ComponentsForm{}
	f.Init(ctx, f, COMPONENTS_PATH, COMPONENTS_ID)
	return f
}

func (f *ComponentsForm) CreateControls(ctx context.Context) {
	f.ButtonPanel = grctl.NewFieldset(f)
	bootstrap.NewButton(f.ButtonPanel).SetText("Button1")
}

func (f *ComponentsForm) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	}
}


func init() {
	page.RegisterPage(COMPONENTS_PATH, NewComponentsForm, COMPONENTS_ID)
}

