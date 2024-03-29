package panels

import (
	"context"
	. "github.com/goradd/goradd/pkg/bootstrap/control"
	"github.com/goradd/goradd/pkg/bootstrap/examples"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
)

const (
	PopupClick int = iota + 10
)

type ModalPanel struct {
	control.Panel
}

func NewModalPanel(ctx context.Context, parent page.ControlI) {
	p := new(ModalPanel)
	p.Init(p, ctx, parent, "modalPanel")

}

func (p *ModalPanel) Init(self any, ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, id)
	p.Panel.AddControls(ctx,
		ButtonCreator{
			ID:       "popupButton",
			Text:     "Popup Modal",
			OnSubmit: action.Do().ID(PopupClick),
		},
	)

	m := NewModal(p.ParentForm(), "modal")
	m.AddCloseButton("Close Me", "close")
	m.SetTitle("My Modal")

	t := control.NewPanel(m, "modbody")
	t.SetText("What is in the modal?")
}

func (p *ModalPanel) DoAction(ctx context.Context, a action.Params) {
	switch a.ID {
	case PopupClick:
		m := GetModal(p, "modal")
		m.Show()
	}
}

func init() {
	examples.RegisterPanel("modal", "Modal", NewModalPanel, 6)
	page.RegisterControl(&ModalPanel{})
	//browsertest.RegisterTestFunction("Bootstrap Standard Form Do Submit", testForms1AjaxSubmit)
	//browsertest.RegisterTestFunction("Bootstrap Standard Form Server Submit", testForms1ServerSubmit)
}
