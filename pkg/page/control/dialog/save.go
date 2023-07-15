package dialog

import (
	"context"
	"github.com/goradd/goradd/pkg/page/control"

	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
)

const (
	saveDlgSaveAction = iota + 10200
)

type SaveablePanel interface {
	control.PanelI
	Load(ctx context.Context, pk string) error
	Save(ctx context.Context)
	Data() interface{}
}

// SavePanel is a dialog panel that pre-loads Save, Cancel and Delete buttons, and treats its one
// child control as an EditablePanel.
type SavePanel struct {
	DialogPanel
}

// GetSavePanel creates a panel that is designed to hold an SaveablePanel. It itself will be wrapped with
// the application's default dialog style, and it will automatically get Save, Cancel and Delete buttons.
func GetSavePanel(parent page.ControlI, id string) (*SavePanel, bool) {
	if parent.Page().HasControl(id) { // dialog has already been created, but is hidden
		return parent.Page().GetControl(id).(*SavePanel), false
	}

	dlg := NewDialogI(parent.ParentForm(), id+"-dlg")
	dp := new(SavePanel)
	dp.Init(dp, dlg, id)
	return dp, true
}

func (p *SavePanel) Init(self any, dlg page.ControlI, id string) {
	p.DialogPanel.Init(self, dlg, id)
	p.AddCloseButton(p.GT("Cancel"), CancelButtonnID)
	p.AddButton(p.GT("Save"), SaveButtonID, &ButtonOptions{
		Validates: true,
		IsPrimary: true,
		OnClick:   action.Ajax(p.ID(), saveDlgSaveAction),
	})
}

func (p *SavePanel) Load(ctx context.Context, pk string) (data interface{}, err error) {
	ep := p.SavePanel()
	if ep == nil {
		panic("the child of a SaveDialog must be an SaveablePanel")
	}

	err = ep.Load(ctx, pk)
	if err != nil {
		return
	}

	data = ep.Data()
	return
}

func (p *SavePanel) DoAction(ctx context.Context, a action.Params) {
	switch a.ID {
	case saveDlgSaveAction:
		p.SavePanel().Save(ctx)
		p.Hide()
	}
}

// SavePanel returns the panel that has the controls
func (p *SavePanel) SavePanel() SaveablePanel {
	children := p.Children()
	if len(children) == 0 {
		return nil
	}
	for i := len(children) - 1; i >= 0; i-- {
		if c, ok := children[i].(SaveablePanel); ok {
			return c
		}
	}
	return nil
}

func init() {
	page.RegisterControl(&SavePanel{})
}
