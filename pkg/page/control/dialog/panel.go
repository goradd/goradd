package dialog

import (
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/html5tag"
)

// ButtonOptions are optional additional items you can add to a dialog button.
type ButtonOptions struct {
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

// Style represents the style of the dialog, whether its a plain dialog (the default),
// or whether it should display additional indicators showing that its indicating an error, warning,
// information, or success. Not all css frameworks support all of these styles.
type Style int

const (
	DefaultStyle Style = iota
	ErrorStyle
	WarningStyle
	InfoStyle
	SuccessStyle
)

// A DialogPanel is the interface between the default dialog style, and a panel.
// To put a dialog on the screen, call GetDialogPanel()
// and then add child controls to that panel, call AddButton() to add buttons to the dialog, and then call Show().
type DialogPanel struct {
	control.Panel
}

// GetDialogPanel returns a new dialog panel if the given dialog panel does not already exist on the form,
// or it returns the dialog panel with the given id that already exists. isNew will indicate whether it
// created a new dialog, or is returning an existing one.
func GetDialogPanel(parent page.ControlI, id string) (dialogPanel *DialogPanel, isNew bool) {
	if parent.Page().HasControl(id) { // dialog has already been created, but is hidden
		return parent.Page().GetControl(id).(*DialogPanel), false
	}

	dlg := NewDialogI(parent.ParentForm(), id+"-dlg")
	dialogPanel = new(DialogPanel)
	dialogPanel.Init(dialogPanel, dlg, id)
	return dialogPanel, true
}

// Init is called by the framework to initialize the
func (p *DialogPanel) Init(self any, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, id)
	p.AddClass("gr-dlg-pnl") // Give the ability to provide a consistent style across the app to all panels in a dialog
}

func (p *DialogPanel) getDialog() DialogI {
	return p.Page().GetControl(p.ID() + "-dlg").(DialogI)
}

// OnClose attaches an action that will happen when the dialog closes.
func (p *DialogPanel) OnClose(a action.ActionI) {
	p.getDialog().On(event.DialogClosed().Validate(event.ValidateNone), a)
}

// OnButton attaches an action handler that responds to button presses. The id of the pressed button will
// be in the Event value of the action.
func (p *DialogPanel) OnButton(a action.ActionI) {
	p.getDialog().On(event.DialogButton().Action(a))
}

// Show will bring a hidden dialog up on the screen.
func (p *DialogPanel) Show() {
	p.getDialog().Show()
}

// Hide will make a dialog invisible. The dialog will still be part of the form object.
func (p *DialogPanel) Hide() {
	p.getDialog().Hide()
}

// SetTitle sets the title of the dialog
func (p *DialogPanel) SetTitle(t string) {
	p.getDialog().SetTitle(t)
}

// SetDialogStyle sets the style of the dialog.
func (p *DialogPanel) SetDialogStyle(s Style) {
	p.getDialog().SetDialogStyle(s)
	p.Refresh()
}

// SetHasCloseBox will put a close box in the upper right corner of the dialog
func (p *DialogPanel) SetHasCloseBox(h bool) {
	p.getDialog().SetHasCloseBox(h)
}

// AddButton adds a button to the dialog.
func (p *DialogPanel) AddButton(label string, id string, options *ButtonOptions) {
	p.getDialog().AddButton(label, id, options)
}

// AddCloseButton will add a button to the dialog that just closes the dialog.
func (p *DialogPanel) AddCloseButton(label string, id string) {
	p.getDialog().AddCloseButton(label, id)
}

// SetButtonVisible will show or hide a specific button that has already been added to the dialog.
func (p *DialogPanel) SetButtonVisible(id string, visible bool) {
	p.getDialog().SetButtonVisible(id, visible)
}

// SetButtonStyle sets the style of the given button
func (p *DialogPanel) SetButtonStyle(id string, a html5tag.Style) {
	p.getDialog().SetButtonStyle(id, a)
}

// SetButtonText sets the text of the given button
func (p *DialogPanel) SetButtonText(id string, text string) {
	p.getDialog().SetButtonText(id, text)
}

// MergeButtonAttributes merges the give attributes into the button's current attributes.
func (p *DialogPanel) MergeButtonAttributes(id string, a html5tag.Attributes) {
	p.getDialog().MergeButtonAttributes(id, a)
}

// RemoveButton removes the given button from the dialog
func (p *DialogPanel) RemoveButton(id string) {
	p.getDialog().RemoveButton(id)
}

// RemoveAllButtons removes all the buttons from the dialog
func (p *DialogPanel) RemoveAllButtons() {
	p.getDialog().RemoveAllButtons()
}

func init() {
	page.RegisterControl(&DialogPanel{})
}
