package control

import (
	"context"
	"fmt"

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
	PanelI
	Load(ctx context.Context, pk string) error
	Save(ctx context.Context)
	Delete(ctx context.Context)
	DataI() interface{}
}

// DialogEditPanel is a dialog panel that pre-loads Save, Cancel and Delete buttons, and treats its one
// child control as an EditablePanel.
type DialogEditPanel struct {
	DialogPanel
	ObjectName string
}

// GetDialogEditPanel creates a panel that is designed to hold an EditablePanel. It itself will be wrapped with
// the application's default dialog style, and it will automatically get Save, Cancel and Delete buttons.
func GetDialogEditPanel(parent page.ControlI, id string, objectName string) (*DialogEditPanel, bool) {
	if parent.Page().HasControl(id) { // dialog has already been created, but is hidden
		return parent.Page().GetControl(id).(*DialogEditPanel), false
	}

	dlg := NewDialogI(parent.ParentForm(), id+"-dlg")
	dp := new(DialogEditPanel)
	dp.Self = dp
	dp.Init(dlg, id, objectName)
	return dp, true
}

func (p *DialogEditPanel) Init(dlg page.ControlI, id string, objectName string) {
	p.DialogPanel.Init(dlg, id)
	p.AddButton(p.GT("Delete"),
		DeleteButtonID,
		&DialogButtonOptions{
			PushLeft:            true,
			ConfirmationMessage: fmt.Sprintf(p.GT("Are you sure you want to delete this %s?"), objectName),
			OnClick:             action.Ajax(p.ID(), editDlgDeleteAction),
		})

	p.AddCloseButton(p.GT("Cancel"), CancelButtonnID)
	p.AddButton(p.GT("Save"), SaveButtonID, &DialogButtonOptions{
		Validates: true,
		IsPrimary: true,
		OnClick:   action.Ajax(p.ID(), editDlgSaveAction),
	})

	p.ObjectName = objectName
}

func (p *DialogEditPanel) Load(ctx context.Context, pk string) (data interface{}, err error) {
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

func (p *DialogEditPanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case editDlgSaveAction:
		p.EditPanel().Save(ctx)
		p.Hide()
	case editDlgDeleteAction:
		p.EditPanel().Delete(ctx)
		p.Hide()
	}
}

// EditPanel returns the panel that has the edit controls
func (p *DialogEditPanel) EditPanel() EditablePanel {
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
	page.RegisterControl(&DialogEditPanel{})
}
