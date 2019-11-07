package control

import (
	"context"
	"fmt"
	config2 "github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
)

type ModalBackdropType int

const (
	ModalBackdrop ModalBackdropType = iota
	ModalNoBackdrop
	ModalStaticBackdrop
)

type ModalI interface {
	control.DialogI
}

// Modal is a bootstrap modal dialog.
// To use a custom template in a bootstrap modal, add a Panel child element or subclass of a panel
// child element. To use the grid system, add the container-fluid class to that embedded panel.
type Modal struct {
	control.Panel
	isOpen bool

	closeOnEscape bool
	sizeClass     string

	titleBar  *TitleBar
	buttonBar *control.Panel
	backdrop  ModalBackdropType

	foundRight bool // utility for adding buttons. No need to serialize this.
}

// event codes
const (
	ButtonClick = iota + 3000
	DialogClosing
	DialogClosed
)

func NewModal(parent page.ControlI, id string) *Modal {
	d := &Modal{}
	d.Self = d
	d.Init(parent, id)
	return d
}

func (d *Modal) Init(parent page.ControlI, id string) {
	if id == "" {
		panic("Modals must have an id")
	}

	d.Panel.Init(parent, id)
	d.Tag = "div"
	d.SetShouldAutoRender(true)

	d.SetValidationType(page.ValidateChildrenOnly) // allows sub items to validate and have validation stop here
	d.SetBlockParentValidation(true)
	d.On(event.Event("hide.bs.modal").Validate(page.ValidateNone).Private(), action.Trigger(d.ID(), event.DialogClosingEvent, nil))
	d.On(event.Event("hidden.bs.modal").Validate(page.ValidateNone).Private(), action.Trigger(d.ID(), event.DialogClosedEvent, nil))
	d.On(event.Event("hidden.bs.modal").Validate(page.ValidateNone).Private(), action.Ajax(d.ID(), DialogClosed))
	config2.LoadBootstrap(d.ParentForm())

	d.AddClass("modal fade").
		SetAttribute("tabindex", -1).
		SetAttribute("role", "dialog").
		SetAttribute("aria-labelledby", d.ID()+"-title").
		SetAttribute("aria-hidden", true)
	d.titleBar = NewTitleBar(d, d.ID()+"-titlebar")
	d.buttonBar = control.NewPanel(d, d.ID()+"-btnbar")
}

func (d *Modal) this() ModalI {
	return d.Self.(ModalI)
}

func (d *Modal) SetTitle(t string) {
	if d.titleBar.Title != t {
		d.titleBar.Title = t
		d.titleBar.Refresh()
	}
}

func (d *Modal) SetHasCloseBox(h bool) {
	if d.titleBar.HasCloseBox != h {
		d.titleBar.HasCloseBox = h
		d.titleBar.Refresh()
	}
}

func (d *Modal) SetDialogStyle(style control.DialogStyle) {
	var class string
	switch style {
	case control.DialogStyleDefault:
		class = BackgroundColorNone + " " + TextColorBody
	case control.DialogStyleWarning:
		class = BackgroundColorWarning + " " + TextColorBody
	case control.DialogStyleError:
		class = BackgroundColorDanger + " " + TextColorLight
	case control.DialogStyleSuccess:
		class = BackgroundColorSuccess + " " + TextColorLight
	case control.DialogStyleInfo:
		class = BackgroundColorInfo + " " + TextColorLight
	}
	d.titleBar.RemoveClassesWithPrefix("bg-")
	d.titleBar.RemoveClassesWithPrefix("text-")
	d.titleBar.AddClass(class)
}

func (d *Modal) SetBackdrop(b ModalBackdropType) {
	d.backdrop = b
	d.Refresh()
}

func (d *Modal) Title() string {
	return d.titleBar.Title
}

func (d *Modal) AddTitlebarClass(class string) {
	d.titleBar.AddClass(class)
}

func (d *Modal) DrawingAttributes(ctx context.Context) html.Attributes {
	a := d.Panel.DrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "bs-modal")
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
) {
	if id == "" {
		id = label
	}
	btn := NewButton(d.buttonBar, d.ID()+"-btn-"+id)
	btn.SetLabel(label)

	if options != nil {
		if options.Validates {
			btn.SetValidationType(page.ValidateContainer)
		}

		if !options.PushLeft && !d.foundRight {
			btn.AddClass("ml-auto")
			d.foundRight = true
		}

		if options.IsClose {
			btn.SetAttribute("data-dismiss", "modal") // make it a close button
		} else if options.ConfirmationMessage == "" {
			btn.On(event.Click(), action.Trigger(d.ID(), event.DialogButtonEvent, id))
		} else {
			btn.On(event.Click(),
				action.Group(
					action.Confirm(options.ConfirmationMessage),
					action.Trigger(d.ID(), event.DialogButtonEvent, id),
				),
			)
		}

		if options.Options != nil && len(options.Options) > 0 {
			if _, ok := options.Options["style"]; ok {
				btn.SetButtonStyle(options.Options["style"].(ButtonStyle))
			}
			if _, ok := options.Options["size"]; ok {
				btn.SetButtonSize(options.Options["size"].(ButtonSize))
			}
		}
	}

	d.buttonBar.Refresh()
}

func (d *Modal) RemoveButton(id string) {
	d.buttonBar.RemoveChild(d.ID() + "-btn-" + id)
	d.buttonBar.Refresh()
}

func (d *Modal) RemoveAllButtons() {
	d.buttonBar.RemoveChildren()
	d.Refresh()
}

func (d *Modal) SetButtonVisible(id string, visible bool) ModalI {
	if ctrl := d.buttonBar.Child(d.ID() + "-btn-" + id); ctrl != nil {
		ctrl.SetVisible(visible)
	}

	return d.this()
}

// SetButtonStyle sets css styles on a button that is already in the dialog
func (d *Modal) SetButtonStyle(id string, a html.Style) ModalI {
	if ctrl := d.buttonBar.Child(d.ID() + "-btn-" + id); ctrl != nil {
		ctrl.SetStyles(a)
	}
	return d.this()
}

// AddCloseButton adds a button to the list of buttons with the given label, but this button will trigger the DialogCloseEvent
// instead of the DialogButtonEvent. The button will also close the dialog (by hiding it).
func (d *Modal) AddCloseButton(label string) ModalI {
	d.AddButton(label, "", &control.DialogButtonOptions{IsClose: true})
	return d.this()
}

func (d *Modal) PrivateAction(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case DialogClosed:
		d.closed()
	}
}

func (d *Modal) Show() {
	if d.Parent() == nil {
		d.SetParent(d.ParentForm()) // This is a saved modal which has previously been created and removed. Insert it back into the form.
	}
	d.SetVisible(true)
	d.isOpen = true
	//d.Refresh()
	d.ParentForm().Response().ExecuteJqueryCommand(d.ID(), "modal", page.PriorityLow, "show")
}

func (d *Modal) Hide() {
	d.ParentForm().Response().ExecuteJqueryCommand(d.ID(), "modal", page.PriorityLow, "hide")
}

func (d *Modal) closed() {
	d.isOpen = false
	//d.Remove()
	d.SetVisible(false)
}

func (d *Modal) PutCustomScript(ctx context.Context, response *page.Response) {
	var backdrop interface{}

	switch d.backdrop {
	case ModalBackdrop:
		backdrop = true
	case ModalNoBackdrop:
		backdrop = false
	case ModalStaticBackdrop:
		backdrop = "static"
	}

	script := fmt.Sprintf(
		`jQuery("#%s").modal({backdrop: %#v, keyboard: %t, focus: true, show: %t});`,
		d.ID(), backdrop, d.closeOnEscape, d.isOpen)
	response.ExecuteJavaScript(script, page.PriorityStandard)
}

/**
Alert creates a message dialog.

If you specify no buttons, a close box in the corner will be created that will just close the dialog. If you
specify just a string in buttons, or just one string as a slice of strings, one button will be shown that will just close the message.

If you specify more than one button, the first button will be the default button (the one pressed if the user presses the return key). In
this case, you will need to detect the button by adding a On(event.DialogButton(), action) to the dialog returned.
You will also be responsible for calling "Close()" on the dialog after detecting a button in this case.
*/
func BootstrapAlert(form page.FormI, message string, buttons interface{}) control.DialogI {
	dlg := NewModal(form, "")
	dlg.SetText(message)
	if buttons != nil {
		switch b := buttons.(type) {
		case string:
			dlg.AddCloseButton(b)
		case []string:
			if len(b) == 1 {
				dlg.AddCloseButton(b[0])
			} else {
				dlg.AddButton(b[0], "", &control.DialogButtonOptions{Options: map[string]interface{}{"style": ButtonStylePrimary}})
				for _, l := range b[1:] {
					dlg.AddButton(l, "", nil)
				}
			}
		}
	} else {
		dlg.SetHasCloseBox(true)
	}
	dlg.Show()
	return dlg
}

type TitleBar struct {
	control.Panel
	HasCloseBox bool
	Title       string
}

func NewTitleBar(parent page.ControlI, id string) *TitleBar {
	d := &TitleBar{}
	d.Self = d
	d.Panel.Init(parent, id)
	return d
}

func init() {
	control.SetAlertFunction(BootstrapAlert)
	page.RegisterControl(&Modal{})
	page.RegisterControl(&TitleBar{})
}

type ModalButtonCreator struct {
	Label string
	ID string
	Validates bool
	ConfirmationMessage string
	PushLeft bool
	IsClose bool
	Options map[string]interface{}
}

func ModalButtons (buttons ...ModalButtonCreator) []ModalButtonCreator {
	return buttons
}


type ModalCreator struct {
	ID string
	Title string
	TitlebarClass string
	HasCloseBox bool
	Style control.DialogStyle
	Backdrop ModalBackdropType
	Buttons []ModalButtonCreator
	OnButton action.ActionI
	page.ControlOptions
	Children []page.Creator
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c ModalCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewModal(parent.ParentForm(), c.ID) // modals always use the form as the parent

	if c.Title != "" {
		ctrl.SetTitle(c.Title)
	}
	if c.TitlebarClass != "" {
		ctrl.AddTitlebarClass(c.TitlebarClass)
	}
	if c.HasCloseBox {
		ctrl.SetHasCloseBox(true)
	}
	if c.Style != control.DialogStyleDefault {
		ctrl.SetDialogStyle(c.Style)
	}
	if c.Backdrop != ModalBackdrop {
		ctrl.SetBackdrop(c.Backdrop)
	}
	if c.Buttons != nil {
		for _, button := range c.Buttons {
			ctrl.AddButton(button.Label, button.ID, &control.DialogButtonOptions{
				Validates:           button.Validates,
				ConfirmationMessage: button.ConfirmationMessage,
				PushLeft:            button.PushLeft,
				IsClose:             button.IsClose,
				Options:             button.Options,
			})
		}
	}
	if c.OnButton != nil {
		ctrl.On(event.DialogButton(), c.OnButton)
	}
	if c.Children != nil {
		ctrl.AddControls(ctx, c.Children...)
	}

	return ctrl
}


// GetListGroup is a convenience method to return the control with the given id from the page.
func GetModal(c page.ControlI, id string) *Modal {
	return c.Page().GetControl(id).(*Modal)
}
