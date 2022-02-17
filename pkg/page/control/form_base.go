package control

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
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


