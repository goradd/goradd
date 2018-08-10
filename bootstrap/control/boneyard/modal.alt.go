package boneyard

import (
	"github.com/spekary/goradd/page/control"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/event"
	"github.com/spekary/goradd/page/action"
	"github.com/spekary/goradd/html"
	"goradd-project/config"
	"goradd-project/app"
	"context"
	"github.com/spekary/goradd/util/types"
)

const MessageDialogId = "qAlertDialog" // The control id to use for the reusable global alert dialog.

// Modal is a bootstrap modal dialog.
// To use a custom template in a bootstrap modal, add a Panel child element or subclass of a panel
// child element. To use the grid system, add the container-fluid class to that embedded panel.
type Modal struct {
	control.Panel
	isOpen      bool
	hasCloseBox bool
	title string
	closeOnEscape bool
	sizeClass	string

	buttonOptions types.OrderedMap
}

type ButtonOptions struct {
	Label string `json:"label"`
	Confirmation string `json:"confirm"`
	IsPrimary bool `json:"primary"`
	Style string `json:"style"`
	Size string `json:"size"`
	Id string 	`json:"id"`
	IsClose bool 	`json:"close"`
	PushLeft bool	`json:"left"`
	validates bool
}

// event codes
const (
	ButtonClick = iota + 3000
	DialogClose
)

const DialogButtonEvent = "grdlgbtn"


func NewModal(parent page.ControlI, id string) *Modal {
	d := &Modal{}
	d.Init(d, parent) // parent is always the overlay
	return d
}

func (d *Modal) Init(self page.ControlI, parent page.ControlI, id string) {
	d.Panel.Init(self, parent, id)
	d.Tag = "div"

	d.SetValidationType(page.ValidateChildrenOnly) // allows sub items to validate and have validation stop here
	d.On(event.DialogClose(), action.Ajax(d.ID(), DialogClose), action.PrivateAction{})
	app.LoadBootstrap(d.ParentForm())
	d.ParentForm().AddJavaScriptFile(config.GoraddDir+"/bootstrap/assets/js/gr.bs.modal.js", false, nil)
	d.ParentForm().AddStyleSheetFile(config.GoraddDir+"/bootstrap/assets/css/bootstrap.min.css", nil)
}

// ValidationType overrides the default ValidationType function to determine if the button clicked causes validation
func (d *Modal) ValidationType(e page.EventI) page.ValidationType {
	id := e.GetActionValue().(string)

	if d.buttonOptions.Has(id) {
		options := d.buttonOptions.Get(id).(ButtonOptions)
		if options.validates {
			return page.ValidateChildrenOnly
		}
	}
	return page.ValidateNone
}


func (d *Modal) SetTitle(t string) *Modal {
	d.title = t
	return d
}

func (d *Modal) SetHasCloseBox(h bool) *Modal {
	d.hasCloseBox = h
}

func (d *Modal) Title() string {
	return d.title
}

func (d *Modal) DrawingAttributes() *html.Attributes {
	a := d.Panel.DrawingAttributes()
	a.SetDataAttribute("grctl", "dialog")
	return a
}

// AddButton adds a button to the modal. Buttons should be added in the order to appear.
// Styling options you can include in options.Options:
//  style - ButtonStyle value
//  size - ButtonSize value
func (d *Modal) AddButton(
	label string,
	id string,
	options *control.DialogButtonOptions,
)  {

	option := ButtonOptions{Label: label}

	if id == "" {
		option.Id = label
	}

	if options != nil {
		if options.IsPrimary {
			option.IsPrimary = true
		}

		if options.Validates {
			option.validates = true
		}

		if options.PushLeft {
			btn.AddClass("push-left")
		}

		if options.IsClose {
			if _,ok := options.Options["dismiss"]; ok {
				btn.SetAttribute("data-dismiss", "modal") // make it a close button
			}
		} else if options.ConfirmationMessage == "" {
			btn.OnClick(action.Trigger(d.ID(), DialogButtonEvent, id))
		} else {
			btn.OnClick(
				action.Confirm(options.ConfirmationMessage),
				action.Trigger(d.ID(), DialogButtonEvent, id),
			)
		}

		if options.Options != nil && len(options.Options) > 0 {
			if _,ok := options.Options["style"]; ok {
				btn.SetButtonStyle(options.Options["style"].(ButtonStyle))
			}
			if _,ok := options.Options["size"]; ok {
				btn.SetButtonSize(options.Options["size"].(ButtonSize))
			}
		}

	}

	d.Refresh()
	return
}

func (d *Modal) RemoveButton(id string) {
	d.buttonBar.RemoveChild(d.ID() + "_btn_" + id)
	d.buttonBar.Refresh()
	//delete(d.validators, id)

}

func (d *Modal) RemoveAllButtons() {
	d.buttonBar.RemoveChildren()
	d.Refresh()
	//delete(d.validators, id)
}

func (d *Modal) SetButtonVisible(id string, visible bool) {
	if ctrl := d.buttonBar.Child(d.ID() + "_btn_" + id); ctrl != nil {
		ctrl.SetVisible(false)
	}
}

// SetButtonStyle sets css styles on a button that is already in the dialog
func (d *Modal) SetButtonStyles(id string, a *html.Style) {
	if ctrl := d.buttonBar.Child(d.ID() + "_btn_" + id); ctrl != nil {
		ctrl.SetStyles(a)
	}
}

// AddCloseButton adds a button to the list of buttons with the given label, but this button will trigger the DialogCloseEvent
// instead of the DialogButtonEvent. The button will also close the dialog (by hiding it).
func (d *Modal) AddCloseButton(label string) {
	d.AddButton(label,"", &control.DialogButtonOptions{IsClose:true})
}

func (d *Modal) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case DialogClose:
		d.Close()
	}
}

func (d *Modal) Open() {
	d.SetVisible(true)
	d.isOpen = true
}

func (d *Modal) Close() {
	d.SetVisible(false)
	d.isOpen = false
	parent := d.Parent()
	if len(parent.Children()) == 1 {
		parent.SetVisible(false)
	}
	d.Remove()
}


/**
Alert creates a message dialog.

If you specify no buttons, a close box in the corner will be created that will just close the dialog. If you
specify just a string in buttons, or just one string as a slice of strings, one button will be shown that will just close the message.

If you specify more than one button, the first button will be the default button (the one pressed if the user presses the return key). In
this case, you will need to detect the button by adding a On(event.DialogButton(), action) to the dialog returned.
You will also be responsible for calling "Close()" on the dialog after detecting a button in this case.
*/
func Alert(form page.FormI, message string, buttons interface{}) *Modal {
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


type titleBar struct {
	control.Panel
	hasCloseBox bool
	title string
}

func NewTitleBar(parent page.ControlI, id string) *titleBar {
	d := &titleBar{}
	d.Init(d, parent)
	return d
}