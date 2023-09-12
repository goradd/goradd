package dialog

import (
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
)

// Alert creates an alert type message Dialog.
//
// Alert returns a DialogPanel, which you can use to do additional modifications on the dialog.
// The dialog will immediately be shown.
//
// The parent parameter is a control that will act as the parent of the dialog.
//
// Use title and message to set the title and content of the alert.
//
// hasClose specifies whether the alert will have a close box in the upper corner that automatically closes the alert.
//
// If you specify no buttons, either set hasClose to true, or add buttons on the returned dialog that will close the dialog.
// You can detect the close action by calling OnClose(action) on the returned dialog, and then responding
// to the action in a DoAction handler.
//
// If you specify one button, clicking it will close the dialog and be the same as the close button.
// You can detect the close action by calling OnClose(action) on the returned dialog, and then responding
// to the action in a DoAction handler.
//
// If you specify more buttons, the first button will be the default button (the one pressed if the
// user presses the return key). You will need to detect the button by calling
// OnButton(action) on the dialog panel returned, and then responding to the button clicked in a DoAction handler.
// You will also be responsible for calling Hide() on the dialog panel in that same DoAction handler.
//
// Call SetDialogStyle on the returned DialogPanel to control the look of the alert.
//
// Example:
//
//	  const dlgAction = 1000
//
//	  func (f *MyForm) popupAlert() {
//	      dlg := dialog.Alert(form, "Hello", "Just saying hi", false, "Great", "Thanks", "See ya")
//	      dlg.OnButton(action.Ajax(f.ID(), dlgAction))
//	  }
//
//	  func (f *MyForm) DoAction(ctx context.Context, a action.Params) {
//	      switch a.ID {
//		     case dlgAction:
//	          btnText := a.EventValueString()
//	          print(btnText) // output button to log
//	          dlg := dialog.GetDialog(f, a.ControlId)
//	          dlg.Hide()
//	      }
//	  }
func Alert(parent page.ControlI, title string, message string, hasClose bool, buttons ...string) *DialogPanel {
	dialogPanel, _ := GetDialogPanel(parent, "gr-alert")
	dialogPanel.SetText(message)
	dialogPanel.RemoveAllButtons()
	dialogPanel.SetHasCloseBox(hasClose)
	dialogPanel.SetTitle(title)

	switch len(buttons) {
	case 0:
	case 1:
		dialogPanel.AddCloseButton(buttons[0], "")
	default:
		for _, l := range buttons {
			dialogPanel.AddButton(l, "", nil)
		}
	}

	dialogPanel.Show()
	return dialogPanel
}

// YesNo is an Alert that has just two buttons, a Yes and a No button.
//
// title and message are the title of the alert and the message displayed.
//
// Your DoAction handler must handle hiding the modal.
//
// Example:
//
//	  const dlgAction = 1000
//
//	  func (f *MyForm) popupAlert() {
//	      dialog.YesNo(form, "Hello", "Just saying hi", false, action.Ajax(f.ID(), dlgAction))
//	  }
//
//	  func (f *MyForm) DoAction(ctx context.Context, a action.Params) {
//	      switch a.ID {
//		     case dlgAction:
//	          btnText := a.EventValueString()
//	          print(btnText) // output button to log
//	          dlg := dialog.GetDialog(f, a.ControlId)
//	          dlg.Hide()
//	      }
//	  }
func YesNo(parent page.ControlI, title string, message string, resultAction action.ActionI) *DialogPanel {
	p := Alert(parent, title, message, false, parent.GT("Yes"), parent.GT("No"))
	p.OnButton(resultAction)
	return p
}
