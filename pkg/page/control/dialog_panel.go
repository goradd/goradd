package control

import (
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
)

// DialogButtonOptions are optional additional items you can add to a dialog button.
type DialogButtonOptions struct {
	// Validates indicates that this button will validate the dialog
	Validates bool
	// The ConfirmationMessage string will appear with a yes/no box making sure the user wants the action.
	// This is usually used when the action could be destructive, like a Delete button.
	ConfirmationMessage string
	// PushLeft pushes this button to the left side of the dialog. Buttons are typically aligned right.
	// This is helpful to separate particular buttons from the main grouping of buttons. Be sure to insert
	// all the PushLeft buttons before the other buttons.
	PushLeft bool
	// IsClose will set the button up to automatically close the dialog. Detect closes with the DialogCloseEvent if needed.
	// The button will not send a DialogButton event.
	IsClose bool
	// IsPrimary styles the button as a primary button and makes it the default when a return is pressed
	IsPrimary bool
	// OnClick is the action that will happen when the button is clicked. If you provide an action, the DialogButton event
	// will not be sent to the dialog. If you do not provide an action, the dialog will receive a DialogButton event instead.
	OnClick action.ActionI
	// Options are additional options specific to the dialog implementation you are using.
	Options map[string]interface{}
}

// DialogStyle represents the style of the dialog, whether its a plain dialog (the default),
// or whether it should display additional indicators showing that its indicating an error, warning,
// information, or success. Not all css frameworks support all of these styles.
type DialogStyle int

const (
	DialogStyleDefault DialogStyle = iota
	DialogStyleError
	DialogStyleWarning
	DialogStyleInfo
	DialogStyleSuccess
)

type DialogPanel struct {
	Panel
}

func GetDialogPanel(parent page.ControlI, id string) (dialogPanel *DialogPanel, isNew bool) {
	if parent.Page().HasControl(id)  { // dialog has already been created, but is hidden
		return parent.Page().GetControl(id).(*DialogPanel), false
	}

	dlg := NewDialogI(parent.ParentForm(), id + "-dlg")
	dialogPanel = new(DialogPanel)
	dialogPanel.Self = dialogPanel
	dialogPanel.Init(dlg, id)
	return dialogPanel, true
}

func (p *DialogPanel) Init(parent page.ControlI, id string)  {
	p.Panel.Init(parent, id)
	p.AddClass("gr-dlg-pnl") // Give the ability to provide a consistent style across the app to all panels in a dialog
}


func (p *DialogPanel) GetDialog() DialogI {
	return p.Page().GetControl(p.ID() + "-dlg").(DialogI)
}

func (p *DialogPanel) OnClose(a action.ActionI) {
	p.GetDialog().On(event.DialogClosed().Validate(page.ValidateNone), a)
}

func (p *DialogPanel) OnButton(a action.ActionI) {
	p.GetDialog().On(event.DialogButton(), a)
}

func (p *DialogPanel) Show() {
	p.GetDialog().Show()
}

func (p *DialogPanel) Hide() {
	p.GetDialog().Hide()
}

func (p *DialogPanel) SetTitle(t string) {
	p.GetDialog().SetTitle(t)
}

func (p *DialogPanel) SetDialogStyle(s DialogStyle) {
	p.GetDialog().SetDialogStyle(s)
	p.Refresh()
}

func (p *DialogPanel) SetHasCloseBox(h bool) {
	p.GetDialog().SetHasCloseBox(h)
}

func (p *DialogPanel) AddButton(label string, id string, options *DialogButtonOptions) {
	p.GetDialog().AddButton(label, id, options)
}

func (p *DialogPanel) AddCloseButton(label string, id string) {
	p.GetDialog().AddCloseButton(label, id)
}

func (p *DialogPanel) SetButtonVisible(id string, visible bool) {
	p.GetDialog().SetButtonVisible(id, visible)
}
func (p *DialogPanel) SetButtonStyle(id string, a html.Style) {
	p.GetDialog().SetButtonStyle(id, a)
}

// RemoveButton removes the given button from the dialog
func (p *DialogPanel) RemoveButton(id string) {
	p.GetDialog().RemoveButton(id)
}

// RemoveAllButtons removes all the buttons from the dialog
func (p *DialogPanel) RemoveAllButtons() {
	p.GetDialog().RemoveAllButtons()
}

// Alert is used by the framework to create an alert type message dialog.
//
// If you specify no buttons, a close box in the corner will be created that will just close the dialog. If you
// specify just a string in buttons, or just one string as a slice of strings,
// one button will be shown that will just close the message.
//
// If you specify more than one button, the first button will be the default button (the one pressed if the
// user presses the return key). In this case, you will need to detect the button by calling
// OnButton(action) on the dialog panel returned.
// You will also be responsible for calling Hide() on the dialog panel after detecting a button in this case.
// You can detect a close button by calling OnClose(action).
// Call SetDialogStyle on the result to control the look of the alery.
func Alert(parent page.ControlI, message string, buttons interface{}) *DialogPanel {
	dialogPanel,_ := GetDialogPanel(parent, "gr-alert")
	dialogPanel.SetText(message)
	dialogPanel.RemoveAllButtons()
	if buttons != nil {
		dialogPanel.SetHasCloseBox(false)
		switch b := buttons.(type) {
		case string:
			dialogPanel.AddCloseButton(b, "")
		case []string:
			if len(b) == 1 {
				dialogPanel.AddCloseButton(b[0], "")
			} else {
				for _, l := range b {
					dialogPanel.AddButton(l, "", nil)
				}
			}
		}
	} else {
		dialogPanel.SetHasCloseBox(true)
	}
	dialogPanel.Show()
	return dialogPanel
}

func YesNo(parent page.ControlI, message string, resultAction action.ActionI) *DialogPanel {
	p := Alert(parent, message, []string{parent.GT("Yes"), parent.GT("No")})
	p.OnButton(resultAction)
	return p
}
