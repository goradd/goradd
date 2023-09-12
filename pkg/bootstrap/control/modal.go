package control

import (
	"context"
	"fmt"
	config2 "github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/dialog"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/html5tag"
)

type ModalBackdropType int

const (
	ModalBackdrop       ModalBackdropType = iota // Standard bootstrap backdrop. Clicking the backdrop closes the modal.
	ModalStaticBackdrop                          // Clicking the backdrop will not close the modal.
)

type ModalI interface {
	dialog.DialogI
}

// Modal is a Bootstrap modal dialog control.
//
// To create Modal, you have a few options:
//   - Call NewModal()
//   - Pass a ModalCreator object to the form's AddControls() function.
//   - Since Modal implements dialog.DialogI, you can also call the dialog.Alert function. If you have previously
//     called setupBootstrap() in your project's config/goradd.go file, then that function
//     will call NewModal to create a Bootstrap style modal dialog.
//
// To use a custom template in a bootstrap modal, add a Panel child element or subclass of a panel
// child element. To use the grid system, add the container-fluid class to that embedded panel.
//
// A modal dialog starts out hidden. Call Show() on the modal dialog to display it.
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

const (
	DialogClosed = iota + 10000 // The event code for a dialog that is closing
)

// NewModal creates a new Modal dialog control.
func NewModal(parent page.ControlI, id string) *Modal {
	d := new(Modal)
	d.Init(d, parent, id)
	return d
}

// Init is called by the framework to initialize a modal dialog.
//
// Subclasses should call Init after creating themselves.
func (m *Modal) Init(self any, parent page.ControlI, id string) {

	if id == "" {
		panic("Modals must have an id")
	}

	m.Panel.Init(self, parent, id)
	m.Tag = "div"
	m.SetShouldAutoRender(true)

	m.SetValidationType(event.ValidateChildrenOnly) // allows sub items to validate and have validation stop here
	m.SetBlockParentValidation(true)
	config2.LoadBootstrap(m.ParentForm())

	m.AddClass("modal fade").
		SetAttribute("tabindex", -1).
		SetAttribute("role", "dialog").
		SetAttribute("aria-labelledby", m.ID()+"-title").
		SetAttribute("aria-hidden", true)
	m.titleBar = NewTitleBar(m, m.ID()+"-titlebar")
	m.buttonBar = control.NewPanel(m, m.ID()+"-btnbar")

	m.On(event.DialogClosed().Validate(event.ValidateNone).Private(), action.Do(m.ID(), DialogClosed))
}

func (m *Modal) this() ModalI {
	return m.Self().(ModalI)
}

// SetTitle sets the content of the title of the modal dialog.
func (m *Modal) SetTitle(t string) {
	if m.titleBar.Title != t {
		m.titleBar.Title = t
		m.titleBar.Refresh()
	}
}

// SetHasCloseBox determines if the modal has a close box in the upper corner which will close the dialog.
func (m *Modal) SetHasCloseBox(h bool) {
	if m.titleBar.HasCloseBox != h {
		m.titleBar.HasCloseBox = h
		m.titleBar.Refresh()
	}
}

// SetDialogStyle sets the style of the dialog.
//
// These styles are mapped to Bootstrap TextColor* styles.
func (m *Modal) SetDialogStyle(style dialog.Style) {
	var class string
	switch style {
	case dialog.DefaultStyle:
		class = BackgroundColorNone + " " + TextColorBody
	case dialog.WarningStyle:
		class = BackgroundColorWarning + " " + TextColorBody
	case dialog.ErrorStyle:
		class = BackgroundColorDanger + " " + TextColorLight
	case dialog.SuccessStyle:
		class = BackgroundColorSuccess + " " + TextColorLight
	case dialog.InfoStyle:
		class = BackgroundColorInfo + " " + TextColorLight
	}
	m.titleBar.RemoveClassesWithPrefix("bg-")
	m.titleBar.RemoveClassesWithPrefix("text-")
	m.titleBar.AddClass(class)
}

// SetBackdrop determines whether the modal dialog will close when clicking on the backdrop.
func (m *Modal) SetBackdrop(b ModalBackdropType) {
	m.backdrop = b
	m.Refresh()
}

// Title returns the title of the dialog.
func (m *Modal) Title() string {
	return m.titleBar.Title
}

// AddTitlebarClass adds a css class to the class of the title bar.
func (m *Modal) AddTitlebarClass(class string) {
	m.titleBar.AddClass(class)
}

func (m *Modal) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := m.Panel.DrawingAttributes(ctx)
	a.SetData("grctl", "bs-modal")
	if m.backdrop == ModalStaticBackdrop {
		a.SetData("bsBackdrop", "static")
	}
	return a
}

// AddButton adds a button to the modal. Buttons should be added in the order to appear.
// Styling options you can include in options.Options:
//
//	style - ButtonStyle value
//	size - ButtonSize value
func (m *Modal) AddButton(
	label string,
	id string,
	options *dialog.ButtonOptions,
) {
	if id == "" {
		id = label
	}
	btn := NewButton(m.buttonBar, m.ID()+"-btn-"+id)
	btn.SetLabel(label)

	if options != nil {
		if options.IsClose {
			btn.SetDataAttribute("bsDismiss", "modal") // make it a close button
			btn.SetDataAttribute("bsTarget", m.ID())   // make it a close button
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
			btn.SetValidationType(event.ValidateContainer)
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

// RemoveButton removes the button from the dialog with the given id.
func (m *Modal) RemoveButton(id string) {
	m.buttonBar.RemoveChild(m.ID() + "-btn-" + id)
	m.buttonBar.Refresh()
}

// RemoveAllButtons removes all the buttons from the modal.
func (m *Modal) RemoveAllButtons() {
	m.buttonBar.RemoveChildren()
	m.Refresh()
}

// SetButtonVisible sets the visible state of the button with the given id.
func (m *Modal) SetButtonVisible(id string, visible bool) {
	if ctrl := m.buttonBar.Child(m.ID() + "-btn-" + id); ctrl != nil {
		ctrl.SetVisible(visible)
	}
}

// SetButtonStyle sets css styles on a button that is already in the dialog
func (m *Modal) SetButtonStyle(id string, a html5tag.Style) {
	if ctrl := m.buttonBar.Child(m.ID() + "-btn-" + id); ctrl != nil {
		ctrl.SetStyles(a)
	}
}

// SetButtonText sets the text of a button that is already on the dialog
func (m *Modal) SetButtonText(id string, text string) {
	if ctrl := m.buttonBar.Child(m.ID() + "-btn-" + id); ctrl != nil {
		ctrl.SetText(text)
	}
}

// MergeButtonAttributes merges the give attributes into the button's current attributes.
func (m *Modal) MergeButtonAttributes(id string, a html5tag.Attributes) {
	if ctrl := m.buttonBar.Child(m.ID() + "-btn-" + id); ctrl != nil {
		ctrl.MergeAttributes(a)
	}
}

// AddCloseButton adds a button to the list of buttons with the given label, but this button will trigger the DialogCloseEvent
// instead of the DialogButtonEvent. The button will also close the dialog (by hiding it).
func (m *Modal) AddCloseButton(label string, id string) {
	m.AddButton(label, id, &dialog.ButtonOptions{IsClose: true})
}

// DoPrivateAction is called by the framework to record that a dialog was closed.
func (m *Modal) DoPrivateAction(_ context.Context, a action.Params) {
	switch a.ID {
	case DialogClosed:
		m.closed()
	}
}

// Show will cause the modal to appear.
func (m *Modal) Show() {
	if m.Parent() == nil {
		m.SetParent(m.ParentForm()) // This is a saved modal which has previously been created and removed. Insert it back into the form.
	}
	m.SetVisible(true)
	m.isOpen = true
	//d.Refresh()
	m.ParentForm().Response().ExecuteJavaScript(fmt.Sprintf("bootstrap.Modal.getInstance(document.getElementById('%s')).show();", m.ID()), page.PriorityLow)
}

// Hide will visibly hide the modal, but will keep its html and javascript code in the client.
func (m *Modal) Hide() {
	m.ParentForm().Response().ExecuteJavaScript(fmt.Sprintf("bootstrap.Modal.getInstance(document.getElementById('%s')).hide()", m.ID()), page.PriorityLow)
}

// closed is used by the framework to record that a dialog was closed.
func (m *Modal) closed() {
	m.isOpen = false
	m.ResetValidation()
	m.SetVisible(false)
}

// PutCustomScript is called by the framework to insert the javascript required to manage the Bootstrap modal.
func (m *Modal) PutCustomScript(_ context.Context, response *page.Response) {

	script := fmt.Sprintf(`
var m = new bootstrap.Modal(document.getElementById('%s') , {keyboard: %t});
`, m.ID(), m.closeOnEscape)
	script += fmt.Sprintf(
		`g$("%s").on("hidden.bs.modal", function(){g$("%[1]s").trigger("%s")});`, m.ID(), event.DialogClosedEvent)

	response.ExecuteJavaScript(script, page.PriorityStandard)
}

// Serialize is called by the framework to record the state of the modal dialog object.
func (m *Modal) Serialize(e page.Encoder) {
	m.Panel.Serialize(e)

	if err := e.Encode(m.isOpen); err != nil {
		panic(err)
	}
	if err := e.Encode(m.closeOnEscape); err != nil {
		panic(err)
	}

	if err := e.Encode(m.backdrop); err != nil {
		panic(err)
	}
	if err := e.Encode(m.foundRight); err != nil {
		panic(err)
	}
}

// Deserialize is called by the framework to restore the state of the dialog.
func (m *Modal) Deserialize(d page.Decoder) {
	m.Panel.Deserialize(d)

	if err := d.Decode(&m.isOpen); err != nil {
		panic(err)
	}
	if err := d.Decode(&m.closeOnEscape); err != nil {
		panic(err)
	}

	m.titleBar = m.Page().GetControl(m.ID() + "-titlebar").(*TitleBar)
	m.buttonBar = m.Page().GetControl(m.ID() + "-btnbar").(*control.Panel)

	if err := d.Decode(&m.backdrop); err != nil {
		panic(err)
	}
	if err := d.Decode(&m.foundRight); err != nil {
		panic(err)
	}
}

// TitleBar is a control that displays the title bar portion of a modal dialog.
type TitleBar struct {
	control.Panel
	HasCloseBox bool
	Title       string
}

// NewTitleBar creates a new TitleBar control.
func NewTitleBar(parent page.ControlI, id string) *TitleBar {
	d := new(TitleBar)
	d.Panel.Init(d, parent, id)
	return d
}

func init() {
	page.RegisterControl(&Modal{})
	page.RegisterControl(&TitleBar{})
}

// ModalButtonCreator declares a dialog button to put on a modal. Pass the structure to
// the ModalButtons function.
type ModalButtonCreator struct {
	Label               string
	ID                  string
	Validates           bool
	ConfirmationMessage string
	PushLeft            bool
	IsClose             bool
	Options             map[string]interface{}
}

// ModalButtons is a helper for declaring buttons on a Modal dialog control.
// Pass it to the Buttons parameter of a ModalCreator.
func ModalButtons(buttons ...ModalButtonCreator) []ModalButtonCreator {
	return buttons
}

// ModalCreator declares a Bootstrap modal dialog. Pass this structure to the AddControls function of a form.
//
// For example, the following will create a dialog with a title, text, and two buttons.
//
//		 form.AddControls(
//		   ModalCreator {
//	      ID: "my-modal",
//	      Title: "Look Out!",
//	      Style: dialog.WarningStyle,
//	      Buttons: ModalButtons(
//	        ModalButtonCreator {
//	          Label: "OK",
//	          ID: "ok",
//	          IsClose: true,
//	        },
//	        ModalButtonCreator {
//	          Label: "Cancel",
//	          ID: "cancel",
//	          IsClose: true,
//	        }
//	      ),
//		   }
//		 )
type ModalCreator struct {
	ID            string
	Title         string
	TitlebarClass string
	HasCloseBox   bool
	Style         dialog.Style
	Backdrop      ModalBackdropType
	Buttons       []ModalButtonCreator
	OnButton      action.ActionI
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
	if c.Style != dialog.DefaultStyle {
		ctrl.SetDialogStyle(c.Style)
	}
	if c.Backdrop != ModalBackdrop {
		ctrl.SetBackdrop(c.Backdrop)
	}
	if c.Buttons != nil {
		for _, button := range c.Buttons {
			ctrl.AddButton(button.Label, button.ID, &dialog.ButtonOptions{
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

// GetModal is a convenience method to return the control with the given id from the page.
func GetModal(c page.ControlI, id string) *Modal {
	return c.Page().GetControl(id).(*Modal)
}
