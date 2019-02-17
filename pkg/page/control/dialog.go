package control

import (
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
)

// event codes
const (
	ButtonClick = iota + 3000
	DialogClose
)

const DialogButtonEvent = "gr-dlgbtn"

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


// DialogI defines the publicly consumable api that the goradd framework uses to interact with a dialog.
//
// More and more CSS and javascript frameworks are coming out with their own forms of dialog, which is usually a
// combination of html, css and a javascript widget. goradd has many ways of potentially interacting with
// dialogs, but to be able to inject a dialog into the framework, we need a consistent interface for all to use.
//
// This particular interface has been implemented in a simple default dialog and Bootstrap dialogs.
// As more needs arise, we can modify the interface to accommodate as many frameworks as possible.
//
// Dialog implementations should descend from the Panel control.
// Dialogs should be able to be a member of a form or control object
// and appear with an Open call, but they should also be able to be instantiated on the fly.
// The framework has hooks for both, and if you are creating a dialog implementation,
// see the current Bootstrap implementation for more direction.
// Feel free to implement more than just the functions listed. These are the minimal set to allow goradd to
// use a dialog implementation. When possible, implementations should use the same function signatures found
// here to do the same work. For example, SetHasCloseBox is defined here, and in the Bootstrap Modal implementation
// with the same function signature, and other implementations should attempt to do the same,
// but it is not enforced by an interface.
type DialogI interface {
	PanelI

	SetTitle(string)
	SetDialogStyle(state DialogStyle)
	Open()
	Close()
}

// Our own implementation of a dialog. This works cooperatively with javascript in goradd.js to create a minimal
// implementation of the dialog interface.
type Dialog struct {
	Panel
	buttonBar   *Panel
	titleBar    *Panel
	closeBox    *Button
	isOpen      bool
	dialogStyle DialogStyle
	title       string
	//validators map[string]bool
}

// DialogButtonOptions are optional additional items you can add to a dialog button.
type DialogButtonOptions struct {
	// Validates indicates that this button will validate the dialog
	Validates bool
	// IsPrimary indicates that this is a submit button so the user can press enter to activate it
	IsPrimary bool
	// The ConfirmationMessage string will appear with a yes/no box making sure the user wants the action. 
	// This is usually used when the action could be destructive, like a Delete button.
	ConfirmationMessage string
	// PushLeft pushes this button to the left side of the dialog. Buttons are typically aligned right. 
	// This is helpful to separate particular buttons from the main grouping of buttons.
	PushLeft bool
	// IsClose will set the button up to automatically close the dialog. Detect closes with the DialogCloseEvent if needed.
	// The button will not send a DialogButton event.
	IsClose bool
	// Options are additional options specific to the dialog implementation you are using.
	Options map[string]interface{}
}

// NewDialog creates a new dialog.
func NewDialog(parent page.ControlI, id string) *Dialog {
	d := &Dialog{}

	d.Init(d, parent, id) // parent is always the form
	return d
}

// Init is called by subclasses of the dialog.
func (d *Dialog) Init(self DialogI, parent page.ControlI, id string) {
	// We add the dialog to the form. The form acts as a dialog controller/container too.
	overlay := parent.Page().GetControl("groverlay")

	if overlay == nil {
		overlay = NewPanel(parent.ParentForm(), "groverlay")
		overlay.SetShouldAutoRender(true)
	} else {
		overlay.SetVisible(true)
	}

	d.Panel.Init(self, overlay, id)
	d.Tag = "div"

	d.titleBar = NewPanel(d, d.ID() + "_title")
	d.titleBar.AddClass("gr-dialog-title")

	d.buttonBar = NewPanel(d, d.ID() + "_buttons")
	d.buttonBar.AddClass("gr-dialog-buttons")
	d.SetValidationType(page.ValidateChildrenOnly) // allows sub items to validate and have validation stop here
	d.On(event.DialogClosed(), action.Ajax(d.ID(), DialogClose), action.PrivateAction{})

	//d.FormBase().AddStyleSheetFile(config.GORADD_FONT_AWESOME_CSS, nil)
}

// SetTitle sets the title of the dialog
func (d *Dialog) SetTitle(t string) {
	d.titleBar.SetText(t)
}

// Title returns the title of the dialog
func (d *Dialog) Title() string {
	return d.titleBar.Text()
}

// ΩDrawingAttributes is called by the framework to set temporary attributes just before drawing.
func (d *Dialog) ΩDrawingAttributes() *html.Attributes {
	a := d.Panel.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "dialog")
	return a
}

// AddButton adds the given button to the dialog.
func (d *Dialog) AddButton(
	label string,
	id string,
	options *DialogButtonOptions,
) {
	if label == "" {
		id = label
	}
	btn := NewButton(d.buttonBar, id)
	btn.SetLabel(label)

	if options != nil {
		if options.IsPrimary {
			btn.SetIsPrimary(true)
		}

		if options.Validates {
			//d.validators[id] = true
			btn.SetValidationType(page.ValidateContainer)
		}

		if options.PushLeft {
			btn.AddClass("push-left")
		}

		if options.ConfirmationMessage == "" {
			btn.On(event.Click(), action.Trigger(d.ID(), DialogButtonEvent, id))
		} else {
			btn.On(event.Click(),
				action.Confirm(options.ConfirmationMessage),
				action.Trigger(d.ID(), DialogButtonEvent, id),
			)
		}
	}

	d.Refresh()
}

// RemoveButton removes the given button from the dialog
func (d *Dialog) RemoveButton(id string) {
	d.buttonBar.RemoveChild(id)
	d.Refresh()
	//delete(d.validators, id)

}

// RemoveAllButtons removes all the buttons from the dialog
func (d *Dialog) RemoveAllButtons() {
	d.buttonBar.RemoveChildren()
	d.buttonBar.Refresh()
	//delete(d.validators, id)
}

// SetButtonVisible sets the visible state of the button. Hidden buttons are still rendered, but are 
// styled so that they are not shown.
func (d *Dialog) SetButtonVisible(id string, visible bool) {
	if ctrl := d.buttonBar.Child(id); ctrl != nil {
		ctrl.SetVisible(false)
	}
}

// SetButtonStyle sets css styles on a button that is already in the dialog
func (d *Dialog) SetButtonStyles(id string, a *html.Style) {
	if ctrl := d.buttonBar.Child(id); ctrl != nil {
		ctrl.SetStyles(a)
	}
}

// SetHasCloseBox adds a close box so that the dialog can be closed in a way that is independent of buttons.
// Often this is an X button in the upper right corner of the dialog.
func (d *Dialog) SetHasCloseBox(h bool) {
	if h && d.closeBox == nil {
		d.addCloseBox()
	} else if !h && d.closeBox != nil {
		d.closeBox.Remove()
		d.closeBox = nil
	}
}

func (d *Dialog) addCloseBox() {
	d.closeBox = NewButton(d.titleBar, d.ID() + "_closebox")
	d.closeBox.AddClass("gr-dialog-close")
	d.closeBox.SetText(`<i class="fa fa-times"></i>`)
	d.closeBox.SetEscapeText(false)
	d.closeBox.On(event.Click(), action.Ajax(d.ID(), DialogClose))
}

// AddCloseButton adds a button to the list of buttons with the given label, but this button will trigger the DialogCloseEvent
// instead of the DialogButtonEvent. The button will also close the dialog.
func (d *Dialog) AddCloseButton(label string, id string) {
	btn := NewButton(d.buttonBar, id)
	btn.SetLabel(label)
	btn.On(event.Click(), action.Trigger(d.ID(), event.DialogClosedEvent, nil))
	// Note: We will also do the public doAction with a DialogCloseEvent
}

// Action is called by the framework and will respond to the DialogClose action sent by any close buttons on the 
// page to close the dialog. You do not normally need to call this.
func (d *Dialog) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case DialogClose:
		d.Close()
	}
}

// Open will show the dialog.
func (d *Dialog) Open() {
	d.SetVisible(true)
	d.isOpen = true
}

// Close will hide the dialog.
func (d *Dialog) Close() {
	d.SetVisible(false)
	d.isOpen = false
	parent := d.Parent()
	if len(parent.Children()) == 1 {
		parent.SetVisible(false)
	}
	d.Remove()
}

// SetDialogStyle sets the style of the dialog.
func (d *Dialog) SetDialogStyle(s DialogStyle) {
	d.dialogStyle = s
	d.Refresh()
}

// Alert is used by the framework to create an alert type message dialog.
//
// If you specify no buttons, a close box in the corner will be created that will just close the dialog. If you
// specify just a string in buttons, or just one string as a slice of strings,
// one button will be shown that will just close the message.
//
// If you specify more than one button, the first button will be the default button (the one pressed if the
// user presses the return key). In this case, you will need to detect the button by adding a
// On(event.DialogButton(), action) to the dialog returned.
// You will also be responsible for calling "Close()" on the dialog after detecting a button in this case.
//
// Call SetAlertFunction to register a different alert function for the framework to use.
func Alert(form page.FormI, message string, buttons interface{}) DialogI {
	return alertFunc(form, message, buttons)
}

func defaultAlert(form page.FormI, message string, buttons interface{}) DialogI {
	dlg := NewDialog(form, "")
	dlg.SetText(message)
	if buttons != nil {
		switch b := buttons.(type) {
		case string:
			dlg.AddCloseButton(b,"")
		case []string:
			if len(b) == 1 {
				dlg.AddCloseButton(b[0],"")
			} else {
				dlg.AddButton(b[0], "", &DialogButtonOptions{IsPrimary: true})
				for _, l := range b[1:] {
					dlg.AddButton(l, "", nil)
				}
			}
		}
	} else {
		dlg.SetHasCloseBox(true)
	}
	dlg.Open()
	return dlg
}

type AlertFuncType func(form page.FormI, message string, buttons interface{}) DialogI

var alertFunc AlertFuncType = defaultAlert // default to our built in one

// SetAlertFunction will set the entire framework's alert function to this function. The alert function
// is called whenever the framework needs to display an alert. Currently, this is done only from the code
// generated forms. Css/js frameworks that want to work with goradd should call this from an init()
// function to enable goradd to use it to display its alerts.
func SetAlertFunction(f AlertFuncType) {
	alertFunc = f
}