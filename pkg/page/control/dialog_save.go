package control

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
)

const (
	saveDlgSaveAction = iota + 10200
)

type SaveablePanel interface {
	PanelI
	Load(ctx context.Context, pk string) error
	Save(ctx context.Context)
	Data() interface{}
}

// DialogSavePanel is a dialog panel that pre-loads Save, Cancel and Delete buttons, and treats its one
// child control as an EditablePanel.
type DialogSavePanel struct {
	DialogPanel
}

// GetDialogSavePanel creates a panel that is designed to hold an SaveablePanel. It itself will be wrapped with
// the application's default dialog style, and it will automatically get Save, Cancel and Delete buttons.
func GetDialogSavePanel(parent page.ControlI, id string) (*DialogSavePanel, bool) {
	if parent.Page().HasControl(id) { // dialog has already been created, but is hidden
		return parent.Page().GetControl(id).(*DialogSavePanel), false
	}

	dlg := NewDialogI(parent.ParentForm(), id+"-dlg")
	dp := new(DialogSavePanel)
	dp.Self = dp
	dp.Init(dlg, id)
	return dp, true
}

func (p *DialogSavePanel) Init(dlg page.ControlI, id string) {
	p.DialogPanel.Init(dlg, id)
	p.AddCloseButton(p.GT("Cancel"), "cancel")
	p.AddButton(p.GT("Save"), "saveBtn", &DialogButtonOptions{
		Validates:true,
		IsPrimary:true,
		OnClick:action.Ajax(p.ID(), saveDlgSaveAction),
	})
}

func (p *DialogSavePanel) Load(ctx context.Context, pk string) (data interface{}, err error) {
	ep := p.SavePanel()
	if ep == nil {
		panic ("the child of a SaveDialog must be an SaveablePanel")
	}

	err = ep.Load(ctx, pk)
	if err != nil {
		return
	}

	data = ep.Data()
	return
}

func (p *DialogSavePanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case editDlgSaveAction:
		p.SavePanel().Save(ctx)
		p.Hide()
	}
}

// EditPanel returns the panel that has the edit controls
func (p *DialogSavePanel) SavePanel() SaveablePanel {
	children := p.Children()
	if len(children) == 0 {
		return nil
	}
	for i := len(children) - 1; i >= 0; i-- {
		if c,ok := children[i].(SaveablePanel); ok {
			return c
		}
	}
	return nil
}

func init() {
	page.RegisterControl(&DialogSavePanel{})
}
