package control

import (
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/action"
	"github.com/spekary/goradd/page/event"
	"github.com/spekary/goradd/page/control/control_base"
)

// event codes
const (
	ButtonClick = iota + 3000
	DialogClose
)

const DialogButtonEvent = "gr-dlgbtn"

const (
	DialogStateError = iota + 1
	DialogStateWarning
)

/*
DialogI defines the publicly coconsumable api that the QCubed framework uses to interact with a dialog.

More and more CSS and javascript frameworks are coming out with their own forms of dialog, which is usually a
combination of html tag(s), css and javascript widget. QCubed has many ways of potentially interacting with
* dialogs, but to be able to inject a dialog into the framework, we need a consistent interface for all to use.
*
* This particular interface has been implemented in both JQuery UI dialogs and Bootstrap dialogs. As more needs arise,
* we can modify the interface to accomodate as many frammeworks as possible.
*
* Dialogs should descend from the Panel control. Dialogs should be able to be a member of a form or control object
* and appear with an Open call, but they should also be able to be instantiated on the fly. The framework has hooks for
* both, and if you are creating a dialog implementation, see the current JQuery UI and Bootstrap implementations for more
* direction.
*
* Feel free to implement more than just the function listed. These are the minimal set to allow your dialog to be used
* by the default QCubed framework.
*/
type DialogI interface {
	page.ControlI
}

// Our own implementation of a dialog. This works cooperatively with javascript in qcubed.js to create a minimal
// implementation of the dialog interface.
type Dialog struct {
	control_base.Panel
	buttonBar   *Panel
	titleBar    *Panel
	closeBox	*Button
	isOpen      bool
	dialogState int
	title		string
	//validators map[string]bool
}

// DialogButtonOptions are optional additional items you can add to a dialog button.
type DialogButtonOptions struct {
	// Set Validates to true to indicate that this button will validate the dialog
	Validates bool
	// Set IsPrimary to true to make this a submit button so the user can press enter to activate it
	IsPrimary bool
	// ConfirmationMessage will appear with a yes/no box making sure the user wants the action. This is usually used
	// when the action could be destructive, like a Delete button.
	ConfirmationMessage string
	// PushLeft pushes this button to the left side of the dialog. Buttons are typically aligned right. This is helpful to separate particular
	// buttons from the main grouping of buttons.
	PushLeft bool
	// Options are additional options specific to the dialog implementation you are using.
	Options map[string]interface{}
}

func NewDialog(parent page.ControlI) *Dialog {
	d := &Dialog{}
	d.Tag = "div"

	// We add the dialog to the overlay. The overlay acts as a dialog controller/container too.
	overlay := parent.Page().GetControl("groverlay")

	if overlay == nil {
		overlay = NewPanel(parent.Form())
		overlay.SetID("groverlay")
		overlay.SetShouldAutoRender(true)
	} else {
		overlay.SetVisible(true)
	}

	d.Init(overlay) // parent is always the overlay
	return d
}

func (c *Dialog) Init(parent page.ControlI) {
	c.Panel.Init(c, parent)
	c.titleBar = NewPanel(c)
	c.titleBar.SetID(c.ID() + "_title")
	c.titleBar.AddClass("gr-dialog-title")

	c.buttonBar = NewPanel(c)
	c.buttonBar.SetID(c.ID() + "_buttons")
	c.buttonBar.AddClass("gr-dialog-buttons")
	c.SetValidationType(page.ValidateChildrenOnly) // allows sub items to validate and have validation stop here
	c.On(event.DialogClose(), action.Ajax(c.ID(), DialogClose), action.PrivateAction{})

	//c.Form().AddStyleSheetFile(config.GORADD_FONT_AWESOME_CSS, nil)
}

func (c *Dialog) SetTitle(t string) *Dialog {
	c.titleBar.SetText(t)
	return c
}

func (c *Dialog) Title() string {
	return c.titleBar.Text()
}


func (c *Dialog) DrawingAttributes() *html.Attributes {
	a := c.Panel.DrawingAttributes()
	a.SetDataAttribute("grctl", "dialog")
	return a
}

func (c *Dialog) AddButton(
	label string,
	id string,
	options *DialogButtonOptions,
) page.ControlI {
	btn := NewButton(c.buttonBar)
	btn.SetLabel(label)
	if label == "" {
		id = label
	}
	btn.SetID(id)

	if options != nil {
		if options.IsPrimary {
			btn.SetIsPrimary(true)
		}

		if options.Validates {
			//c.validators[id] = true
			btn.SetValidationType(page.ValidateContainer)
		}

		if options.PushLeft {
			btn.AddClass("push-left")
		}

		if options.ConfirmationMessage == "" {
			btn.OnClick(action.Trigger(c.ID(), DialogButtonEvent, id))
		} else {
			btn.OnClick(
				action.Confirm(options.ConfirmationMessage),
				action.Trigger(c.ID(), DialogButtonEvent, id),
			)
		}
	}

	c.Refresh()
	return btn
}

func (c *Dialog) RemoveButton(id string) {
	c.RemoveChild(id)
	c.Refresh()
	//delete(c.validators, id)

}

func (c *Dialog) RemoveAllButtons() {
	c.buttonBar.RemoveChildren()
	c.Refresh()
	//delete(c.validators, id)
}

func (c *Dialog) SetButtonVisible(id string, visible bool) {
	if ctrl := c.buttonBar.Child(id); ctrl != nil {
		ctrl.SetVisible(false)
	}
}

// SetButtonStyle sets css styles on a button that is already in the dialog
func (c *Dialog) SetButtonStyles(id string, a *html.Style) {
	if ctrl := c.buttonBar.Child(id); ctrl != nil {
		ctrl.SetStyles(a)
	}
}

func (c *Dialog) HasCloseBox() page.ControlI {
	c.addCloseBox()
	return c
}

func (c *Dialog) addCloseBox() {
	c.closeBox = NewButton(c.titleBar)
	c.closeBox.AddClass("gr-dialog-close")
	c.closeBox.SetText(`<i class="fa fa-times"></i>`)
	c.closeBox.SetEscapeText(false)
	c.closeBox.OnClick(action.Ajax(c.ID(), DialogClose))
}

// AddCloseButton adds a button to the list of buttons with the given label, but this button will trigger the DialogCloseEvent
// instead of the DialogButtonEvent. The button will also close the dialog (by hiding it).
func (c *Dialog) AddCloseButton(label string) {
	btn := NewButton(c.buttonBar)
	btn.SetLabel(label)
	btn.OnClick(action.Trigger(c.ID(), event.DialogCloseEvent, nil))
	// Note: We will also do the public doAction with a DialogCloseEvent
}

func (c *Dialog) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case DialogClose:
		c.Close()
	}
}

func (c *Dialog) Open() {
	c.SetVisible(true)
	c.isOpen = true
}

func (c *Dialog) Close() {
	c.SetVisible(false)
	c.isOpen = false
	parent := c.Parent()
	if len(parent.Children()) == 1 {
		parent.SetVisible(false)
	}
	c.Remove()
}

func (c *Dialog) SetDialogState(s int) *Dialog {
	c.dialogState = s
	c.Refresh()
	return c
}

/**
Alert creates a message dialog.

If you specify no buttons, a close box in the corner will be created that will just close the dialog. If you
specify just a string in buttons, or just one string as a slice of strings, one button will be shown that will just close the message.

If you specify more than one button, the first button will be the default button (the one pressed if the user presses the return key). In
this case, you will need to detect the button by adding a On(event.DialogButton(), action) to the dialog returned.
You will also be responsible for calling "Close()" on the dialog after detecting a button in this case.
*/
func Alert(form page.FormI, message string, buttons interface{}) *Dialog {
	dlg := NewDialog(form)
	dlg.SetText(message)
	if buttons != nil {
		switch b := buttons.(type) {
		case string:
			dlg.AddCloseButton(b)
		case []string:
			if len(b) == 1 {
				dlg.AddCloseButton(b[0])
			} else {
				dlg.AddButton(b[0], "", &DialogButtonOptions{IsPrimary: true})
				for _, l := range b[1:] {
					dlg.AddButton(l, "", nil)
				}
			}
		}
	} else {
		dlg.HasCloseBox()
	}
	dlg.Open()
	return dlg
}
