package dialog

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/page/control"

	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
)

const (
	editDlgSaveAction = iota + 10100
	editDlgDeleteAction
)

const (
	SaveButtonID    = "saveBtn"
	CancelButtonnID = "cancelBtn"
	DeleteButtonID  = "deleteBtn"
)

type EditablePanel interface {
	control.PanelI
	Load(ctx context.Context, pk string) error
	Save(ctx context.Context)
	Delete(ctx context.Context)
	DataI() interface{}
}

// EditPanel is a dialog panel that pre-loads Save, Cancel and Delete buttons, and treats its one
// child control as an EditablePanel.
type EditPanel struct {
	DialogPanel
	ObjectName string
}

// GetEditPanel creates a panel that is designed to hold an EditablePanel. It itself will be wrapped with
// the application's default dialog style, and it will automatically get Save, Cancel and Delete buttons.
func GetEditPanel(parent page.ControlI, id string, objectName string) (*EditPanel, bool) {
	if parent.Page().HasControl(id) { // dialog has already been created, but is hidden
		return parent.Page().GetControl(id).(*EditPanel), false
	}

	dlg := NewDialogI(parent.ParentForm(), id+"-dlg")
	dp := new(EditPanel)
	dp.Init(dp, dlg, id, objectName)
	return dp, true
}

func (p *EditPanel) Init(self any, dlg page.ControlI, id string, objectName string) {
	p.DialogPanel.Init(self, dlg, id)
	p.AddButton(p.GT("Delete"),
		DeleteButtonID,
		&ButtonOptions{
			PushLeft:            true,
			ConfirmationMessage: fmt.Sprintf(p.GT("Are you sure you want to delete this %s?"), objectName),
			OnClick:             action.Do(p.ID(), editDlgDeleteAction),
		})

	p.AddCloseButton(p.GT("Cancel"), CancelButtonnID)
	p.AddButton(p.GT("Save"), SaveButtonID, &ButtonOptions{
		Validates: true,
		IsPrimary: true,
		OnClick:   action.Do(p.ID(), editDlgSaveAction),
	})

	p.ObjectName = objectName
}

func (p *EditPanel) Load(ctx context.Context, pk string) (data interface{}, err error) {
	ep := p.EditPanel()
	if ep == nil {
		panic("the child of an EditDialog must be an EditablePanel")
	}

	err = ep.Load(ctx, pk)
	if err != nil {
		return
	}

	if pk == "" {
		// Editing a new item
		p.SetButtonVisible(DeleteButtonID, false)
		p.SetTitle(fmt.Sprintf(p.GT("New %s"), p.ObjectName))
	} else {
		p.SetButtonVisible(DeleteButtonID, true)
		p.SetTitle(fmt.Sprintf(p.GT("Edit %s"), p.ObjectName))
	}
	data = ep.DataI()
	return
}

func (p *EditPanel) DoAction(ctx context.Context, a action.Params) {
	switch a.ID {
	case editDlgSaveAction:
		p.EditPanel().Save(ctx)
		p.Hide()
	case editDlgDeleteAction:
		p.EditPanel().Delete(ctx)
		p.Hide()
	default:
		p.DialogPanel.DoAction(ctx, a)
	}
}

// EditPanel returns the panel that has the edit controls
func (p *EditPanel) EditPanel() EditablePanel {
	children := p.Children()
	if len(children) == 0 {
		return nil
	}
	for i := len(children) - 1; i >= 0; i-- {
		if c, ok := children[i].(EditablePanel); ok {
			return c
		}
	}
	return nil
}

func init() {
	page.RegisterControl(&EditPanel{})
}
