package control

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
	"strings"
)

const (
	databaseProfileAction = iota + 10000
)

// The FormBase is the control that all Form objects should include, and is the master container for all other goradd controls.
type FormBase struct {
	page.FormBase
}

type MockForm struct {
	FormBase
}

func init() {
	page.RegisterControl(&MockForm{})
}

// NewMockForm creates a form that should be used as a parent of a control when unit testing the control.
func NewMockForm() *MockForm {
	f := &MockForm{}
	f.Self = f
	f.FormBase.Init(nil, "MockFormId")
	return f
}

func (f *MockForm) AddRelatedFiles() {
}

// Init initializes the FormBase. Call this before adding other controls.
func (f *FormBase) Init(ctx context.Context, id string) {
	// Most of the FormBase code is in page.FormBase. The code below specifically adds popup windows and controls
	// to all standard forms, mostly for debug and development purposes.

	f.FormBase.Init(ctx, id)

	if db.IsProfiling(ctx) {
		btn := NewButton(f, "grProfileButton")
		btn.SetText("SQL Profile <i class='fas fa-arrow-circle-down' ></i>")
		btn.SetTextIsHtml(true)
		btn.On(event.Click(), action.Ajax(f.ID(), databaseProfileAction))
		btn.SetShouldAutoRender(true)

		panel := NewPanel(f, "grProfilePanel")
		panel.SetShouldAutoRender(true)
		panel.SetTextIsHtml(true)
		panel.SetVisible(false)
	}

	/*	TODO: Add a dialog and designer click if in design mode
			if (defined('QCUBED_DESIGN_MODE') && QCUBED_DESIGN_MODE == 1) {
			// Attach custom event to dialog to handle right click menu items sent by form

			$dlg = new Q\ModelConnector\EditDlg ($objClass, 'qconnectoreditdlg');

			$dlg->addAction(
				new Q\Event\On('qdesignerclick'),
				new Q\Action\Ajax ('ctlDesigner_Click', null, null, 'ui')
			);
		}

	*/

}

func (f *FormBase) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case databaseProfileAction:
		if f.Page().HasControl("grProfilePanel") {
			c := f.Page().GetControl("grProfilePanel")
			if c.IsVisible() {
				c.SetVisible(false)
			} else {
				c.SetVisible(true)
				var s string
				if profiles := db.GetProfiles(ctx); profiles != nil {
					for _, profile := range profiles {
						dif := profile.EndTime.Sub(profile.BeginTime)
						sql := strings.Replace(profile.Sql, "\n", "<br />", -1)
						s += fmt.Sprintf(`<p class="profile"><div>Time: %s Begin: %s End: %s</div><div>%s</div></p>`,
							dif.String(), profile.BeginTime.Format("3:04:05.000"), profile.EndTime.Format("3:04:05.000"), sql)
					}
				}
				c.SetText(s)
			}
		}
	}
}


