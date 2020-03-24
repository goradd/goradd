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
	DialogClosed = iota + 10000
)

func NewModal(parent page.ControlI, id string) *Modal {
	d := &Modal{}
	d.Self = d
	d.Init(parent, id)
	return d
}

func (m *Modal) Init(parent page.ControlI, id string) {

	if id == "" {
		panic("Modals must have an id")
	}

	m.Panel.Init(parent, id)
	m.Tag = "div"
	m.SetShouldAutoRender(true)

	m.SetValidationType(page.ValidateChildrenOnly) // allows sub items to validate and have validation stop here
	m.SetBlockParentValidation(true)
	config2.LoadBootstrap(m.ParentForm())

	m.AddClass("modal fade").
		SetAttribute("tabindex", -1).
		SetAttribute("role", "dialog").
		SetAttribute("aria-labelledby", m.ID()+"-title").
		SetAttribute("aria-hidden", true)
	m.titleBar = NewTitleBar(m, m.ID()+"-titlebar")
	m.buttonBar = control.NewPanel(m, m.ID()+"-btnbar")

	m.On(event.DialogClosed().Validate(page.ValidateNone).Private(), action.Ajax(m.ID(), DialogClosed))
}

func (m *Modal) this() ModalI {
	return m.Self.(ModalI)
}

func (m *Modal) SetTitle(t string) {
	if m.titleBar.Title != t {
		m.titleBar.Title = t
		m.titleBar.Refresh()
	}
}

func (m *Modal) SetHasCloseBox(h bool) {
	if m.titleBar.HasCloseBox != h {
		m.titleBar.HasCloseBox = h
		m.titleBar.Refresh()
	}
}

func (m *Modal) SetDialogStyle(style control.DialogStyle) {
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
	m.titleBar.RemoveClassesWithPrefix("bg-")
	m.titleBar.RemoveClassesWithPrefix("text-")
	m.titleBar.AddClass(class)
}

func (m *Modal) SetBackdrop(b ModalBackdropType) {
	m.backdrop = b
	m.Refresh()
}

func (m *Modal) Title() string {
	return m.titleBar.Title
}

func (m *Modal) AddTitlebarClass(class string) {
	m.titleBar.AddClass(class)
}

func (m *Modal) DrawingAttributes(ctx context.Context) html.Attributes {
	a := m.Panel.DrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "bs-modal")
	return a
}

// AddButton adds a button to the modal. Buttons should be added in the order to appear.
// Styling options you can include in options.Options:
//  style - ButtonStyle value
//  size - ButtonSize value
func (m *Modal) AddButton(
	label string,
	id string,
	options *control.DialogButtonOptions,
) {
	if id == "" {
		id = label
	}
	btn := NewButton(m.buttonBar, m.ID()+"-btn-"+id)
	btn.SetLabel(label)

	if options != nil {
		if options.IsClose {
			btn.SetAttribute("data-dismiss", "modal") // make it a close button
		} else if options.ConfirmationMessage == "" {
			if options.OnClick != nil {
				btn.On(event.Click(), options.OnClick)
			} else {
				btn.On(event.Click(), action.Trigger(m.ID(), event.DialogButtonEvent, id))
			}
		} else {
			if options.OnClick != nil {
				btn.On(event.Click(),
					action.Group(
						action.Confirm(options.ConfirmationMessage),
						options.OnClick,
					),
				)
			} else {
				btn.On(event.Click(),
					action.Group(
						action.Confirm(options.ConfirmationMessage),
						action.Trigger(m.ID(), event.DialogButtonEvent, id),
					),
				)
			}
		}
	} else {
		btn.On(event.Click(), action.Trigger(m.ID(), event.DialogButtonEvent, id))
	}

	if (options == nil || !options.PushLeft) && !m.foundRight {
		btn.AddClass("ml-auto")
		m.foundRight = true
	}


	if options != nil {
		if options.Validates {
			btn.SetValidationType(page.ValidateContainer)
		}

		if options.Options != nil && len(options.Options) > 0 {
			if _, ok := options.Options["style"]; ok {
				btn.SetButtonStyle(options.Options["style"].(ButtonStyle))
			}
			if _, ok := options.Options["size"]; ok {
				btn.SetButtonSize(options.Options["size"].(ButtonSize))
			}
		}

		if options.IsPrimary {
			btn.SetIsPrimary(true)
		}
	}

	m.buttonBar.Refresh()
}

func (m *Modal) RemoveButton(id string) {
	m.buttonBar.RemoveChild(m.ID() + "-btn-" + id)
	m.buttonBar.Refresh()
}

func (m *Modal) RemoveAllButtons() {
	m.buttonBar.RemoveChildren()
	m.Refresh()
}

func (m *Modal) SetButtonVisible(id string, visible bool) {
	if ctrl := m.buttonBar.Child(m.ID() + "-btn-" + id); ctrl != nil {
		ctrl.SetVisible(visible)
	}
}

// SetButtonStyle sets css styles on a button that is already in the dialog
func (m *Modal) SetButtonStyle(id string, a html.Style) {
	if ctrl := m.buttonBar.Child(m.ID() + "-btn-" + id); ctrl != nil {
		ctrl.SetStyles(a)
	}
}

// AddCloseButton adds a button to the list of buttons with the given label, but this button will trigger the DialogCloseEvent
// instead of the DialogButtonEvent. The button will also close the dialog (by hiding it).
func (m *Modal) AddCloseButton(label string, id string) {
	m.AddButton(label, id, &control.DialogButtonOptions{IsClose: true})
}

func (m *Modal) PrivateAction(_ context.Context, a page.ActionParams) {
	switch a.ID {
	case DialogClosed:
		m.closed()
	}
}

func (m *Modal) Show() {
	if m.Parent() == nil {
		m.SetParent(m.ParentForm()) // This is a saved modal which has previously been created and removed. Insert it back into the form.
	}
	m.SetVisible(true)
	m.isOpen = true
	//d.Refresh()
	m.ParentForm().Response().ExecuteJqueryCommand(m.ID(), "modal", page.PriorityLow, "show")
}

func (m *Modal) Hide() {
	m.ParentForm().Response().ExecuteJqueryCommand(m.ID(), "modal", page.PriorityLow, "hide")
}

func (m *Modal) closed() {
	m.isOpen = false
	m.ResetValidation()
	m.SetVisible(false)
}

func (m *Modal) PutCustomScript(_ context.Context, response *page.Response) {
	var backdrop interface{}

	switch m.backdrop {
	case ModalBackdrop:
		backdrop = true
	case ModalNoBackdrop:
		backdrop = false
	case ModalStaticBackdrop:
		backdrop = "static"
	}

	script := fmt.Sprintf(
		`jQuery("#%s").modal({backdrop: %#v, keyboard: %t, focus: true, show: %t});`,
		m.ID(), backdrop, m.closeOnEscape, m.isOpen)
	script += fmt.Sprintf(
		`jQuery("#%s").on("hidden.bs.modal", function(){g$("%[1]s").trigger("grdlgclosed")});`, m.ID())

	response.ExecuteJavaScript(script, page.PriorityStandard)
}

func (m *Modal) Serialize(e page.Encoder) (err error) {
	if err = m.Panel.Serialize(e); err != nil {
		return
	}

	if err = e.Encode(m.isOpen); err != nil {
		panic(err)
	}
	if err = e.Encode(m.closeOnEscape); err != nil {
		panic(err)
	}

	if err = e.Encode(m.backdrop); err != nil {
		panic(err)
	}
	if err = e.Encode(m.foundRight); err != nil {
		panic(err)
	}

	return
}


func (m *Modal) Deserialize(d page.Decoder) (err error) {
	if err = m.Panel.Deserialize(d); err != nil {
		return
	}

	if err = d.Decode(&m.isOpen); err != nil {
		return
	}
	if err = d.Decode(&m.closeOnEscape); err != nil {
		return
	}

	m.titleBar = m.Page().GetControl(m.ID()+"-titlebar").(*TitleBar)
	m.buttonBar = m.Page().GetControl(m.ID()+"-btnbar").(*control.Panel)

	if err = d.Decode(&m.backdrop); err != nil {
		return
	}
	if err = d.Decode(&m.foundRight); err != nil {
		return
	}

	return
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
