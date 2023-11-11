// Package dialog contains dialog controls.
//
// Dialogs use javascript to popup a "window" above the current page.
package dialog

import (
	"context"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/button"
	"io"

	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/html5tag"
)

const (
	ClosedAction = iota + 3000
)

const OverlayID = "gr-dlg-overlay"

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
	control.PanelI

	SetTitle(string)
	SetDialogStyle(state Style)
	SetHasCloseBox(bool)
	Show()
	Hide()
	AddButton(label string, id string, options *ButtonOptions)
	AddCloseButton(label string, id string)
	SetButtonText(id string, text string)
	SetButtonVisible(id string, visible bool)
	SetButtonStyle(id string, a html5tag.Style)
	MergeButtonAttributes(id string, a html5tag.Attributes)
	RemoveButton(id string)
	RemoveAllButtons()
}

// Dialog is the default implementation of a dialog in GoRADD. You should not normally call this directly, but
// rather call GetDialogPanel to create a dialog. GetDialogPanel will then call NewDialogI to create a dialog
// that wraps the panel. To change the default dialog style to a different one, call SetNewDialogFunction()
type Dialog struct {
	control.Panel
	buttonBarID string
	titleBarID  string
	closeBoxID  string
	dialogStyle Style
	title       string
}

// NewDialog creates a new dialog.
func NewDialog(parent page.ControlI, id string) *Dialog {
	d := new(Dialog)
	d.Init(d, parent, id) // parent is always the form
	return d
}

// Init is called by subclasses of the dialog.
func (d *Dialog) Init(self any, parent page.ControlI, id string) {
	// Our strategy here is to create a dialog overlay that is a container for the currently shown dialogs. This
	// container is owned by the form itself, even if sub-controls create the dialog.
	var overlay *control.Panel

	if id == "" {
		panic("Dialogs must have an id.")
	}

	if !parent.Page().HasControl(OverlayID) {
		overlay = control.NewPanel(parent.ParentForm(), OverlayID)
		overlay.SetShouldAutoRender(true)
	} else {
		overlay = control.GetPanel(parent, OverlayID)
	}

	// Make the overlay our parent
	d.Panel.Init(self, overlay, id)
	d.Tag = "div"

	d.titleBarID = d.ID() + "-title"
	tb := control.NewPanel(d, d.titleBarID)
	tb.AddClass("gr-dialog-title")

	d.buttonBarID = d.ID() + "-buttons"
	bb := control.NewPanel(d, d.buttonBarID)
	bb.AddClass("gr-dialog-buttons")
	d.SetValidationType(event.ValidateChildrenOnly) // allows sub items to validate and have validation stop here
	d.On(event.DialogClosed().Validate(event.ValidateNone).Private(), action.Do(d.ID(), ClosedAction))
}

func (d *Dialog) TitleBar() *control.Panel {
	return control.GetPanel(d, d.titleBarID)
}

func (d *Dialog) ButtonBar() *control.Panel {
	return control.GetPanel(d, d.buttonBarID)
}

func (d *Dialog) CloseBox() *button.Button {
	if d.closeBoxID == "" {
		return nil
	}
	return button.GetButton(d, d.closeBoxID)
}

// SetTitle sets the title of the dialog
func (d *Dialog) SetTitle(t string) {
	d.TitleBar().SetText(t)
}

// Title returns the title of the dialog
func (d *Dialog) Title() string {
	return d.TitleBar().Text()
}

// DrawingAttributes is called by the framework to set temporary attributes just before drawing.
func (d *Dialog) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := d.Panel.DrawingAttributes(ctx)
	a.SetData("grctl", "dialog")
	return a
}

func (d *Dialog) DrawInnerHtml(ctx context.Context, w io.Writer) {
	control.GetPanel(d, d.titleBarID).Draw(ctx, w)
	control.GetPanel(d, d.buttonBarID).Draw(ctx, w)
	page.WriteString(w, `<div class="gr-dlg-content">`)
	d.Panel.DrawInnerHtml(ctx, w)
	page.WriteString(w, `</div>`)

	return
}

// AddButton adds the given button to the dialog.
func (d *Dialog) AddButton(
	label string,
	id string,
	options *ButtonOptions,
) {
	if label == "" {
		id = label
	}
	btn := button.NewButton(d.ButtonBar(), id)
	btn.SetLabel(label)

	if options != nil {
		if options.Validates {
			//d.validators[id] = true
			btn.SetValidationType(event.ValidateContainer)
		}

		if options.IsPrimary {
			btn.SetIsPrimary(true)
		}

		if options.PushLeft {
			btn.AddClass("push-left")
			btn.SetAttribute("tabindex", 10000)
		} else {
			btn.SetAttribute("tabindex", 10001) // make sure right buttons tab after left buttons
		}

		if options.ConfirmationMessage == "" {
			if options.OnClick != nil {
				btn.On(event.Click(), options.OnClick)
			} else {
				btn.On(event.Click(), action.Trigger(d.ID(), event.DialogButtonEvent, id))
			}
		} else {
			if options.OnClick != nil {
				btn.On(event.Click(),
					action.Confirm(options.ConfirmationMessage),
					options.OnClick,
				)
			} else {
				btn.On(event.Click(),
					action.Confirm(options.ConfirmationMessage),
					action.Trigger(d.ID(), event.DialogButtonEvent, id),
				)
			}
		}
	} else {
		btn.On(event.Click(), action.Trigger(d.ID(), event.DialogButtonEvent, id))
	}

	d.Refresh()
}

// RemoveButton removes the given button from the dialog
func (d *Dialog) RemoveButton(id string) {
	d.ButtonBar().RemoveChild(id)
	d.Refresh()
	//delete(d.validators, id)
}

// RemoveAllButtons removes all the buttons from the dialog
func (d *Dialog) RemoveAllButtons() {
	bb := d.ButtonBar()
	bb.RemoveChildren()
	bb.Refresh()
	//delete(d.validators, id)
}

// SetButtonVisible sets the visible state of the button. Hidden buttons are still rendered, but are
// styled so that they are not shown.
func (d *Dialog) SetButtonVisible(id string, visible bool) {
	bb := d.ButtonBar()
	if ctrl := bb.Child(id); ctrl != nil {
		ctrl.SetVisible(visible)
	}
}

// SetButtonText sets the text of a button that was previously created
func (d *Dialog) SetButtonText(id string, text string) {
	bb := d.ButtonBar()
	if ctrl := bb.Child(id); ctrl != nil {
		ctrl.SetText(text)
	}
}

// SetButtonStyle sets css styles on a button that is already in the dialog
func (d *Dialog) SetButtonStyle(id string, a html5tag.Style) {
	bb := d.ButtonBar()
	if ctrl := bb.Child(id); ctrl != nil {
		ctrl.SetStyles(a)
	}
}

// MergeButtonAttributes merges the given attributes into the button's attributes.
func (d *Dialog) MergeButtonAttributes(id string, a html5tag.Attributes) {
	bb := d.ButtonBar()
	if ctrl := bb.Child(id); ctrl != nil {
		ctrl.MergeAttributes(a)
	}
}

// SetHasCloseBox adds a close box so that the dialog can be closed in a way that is independent of buttons.
// Often this is an X button in the upper right corner of the dialog.
func (d *Dialog) SetHasCloseBox(h bool) {
	cb := d.CloseBox()
	if h && cb == nil {
		d.addCloseBox()
	} else if !h && cb != nil {
		cb.Remove()
		d.closeBoxID = ""
	}
}

func (d *Dialog) addCloseBox() {
	d.closeBoxID = d.ID() + "-cb"
	cb := button.NewButton(d.TitleBar(), d.closeBoxID)
	cb.AddClass("gr-dialog-close")
	cb.SetText(`<span">X</span>`)
	cb.SetTextIsHtml(true)
	cb.On(event.Click(), action.Trigger(d.ID(), event.DialogClosedEvent, nil))
}

// AddCloseButton adds a button to the list of buttons with the given label, but this button will trigger the DialogCloseEvent
// instead of the DialogButtonEvent. The button will also close the dialog.
func (d *Dialog) AddCloseButton(label string, id string) {
	btn := button.NewButton(d.ButtonBar(), id)
	btn.SetLabel(label)
	btn.SetValidationType(event.ValidateNone)
	btn.On(event.Click(), action.Trigger(d.ID(), event.DialogClosedEvent, nil))
}

// DoPrivateAction is called by the framework and will respond to the DialogClose action sent by any close buttons on the
// page to close the dialog.
func (d *Dialog) DoPrivateAction(_ context.Context, a action.Params) {
	switch a.ID {
	case ClosedAction:
		d.Hide()
	}
}

// Show will show the dialog.
func (d *Dialog) Show() {
	overlay := control.GetPanel(d, OverlayID)
	overlay.SetVisible(true)
	d.SetVisible(true)
}

// Hide will hide the dialog. The dialog will still be part of the form, just in a hidden state.
func (d *Dialog) Hide() {
	d.SetVisible(false)
	overlay := control.GetPanel(d, OverlayID)
	var vis bool
	for _, child := range overlay.Children() {
		if child.IsVisible() {
			vis = true
			break
		}
	}

	if !vis {
		overlay.SetVisible(false) // hide the overlay if all of the enclosed dialogs are not visible
	}
}

// SetDialogStyle sets the style of the dialog.
func (d *Dialog) SetDialogStyle(s Style) {
	d.dialogStyle = s
	d.Refresh()
}

func (d *Dialog) Serialize(e page.Encoder) {
	d.ControlBase.Serialize(e)
	if err := e.Encode(d.buttonBarID); err != nil {
		panic(err)
	}
	if err := e.Encode(d.titleBarID); err != nil {
		panic(err)
	}
	if err := e.Encode(d.closeBoxID); err != nil {
		panic(err)
	}
	if err := e.Encode(d.dialogStyle); err != nil {
		panic(err)
	}
	if err := e.Encode(d.title); err != nil {
		panic(err)
	}
}

func (d *Dialog) Deserialize(dec page.Decoder) {
	d.ControlBase.Deserialize(dec)
	if err := dec.Decode(&d.buttonBarID); err != nil {
		panic(err)
	}
	if err := dec.Decode(&d.titleBarID); err != nil {
		panic(err)
	}
	if err := dec.Decode(&d.closeBoxID); err != nil {
		panic(err)
	}
	if err := dec.Decode(&d.dialogStyle); err != nil {
		panic(err)
	}
	if err := dec.Decode(&d.title); err != nil {
		panic(err)
	}
}

// NewDialogI creates a new dialog in a css framework independent way, by returning a DialogI interface.
// Call SetNewDialogFunction() to set the function that controls how dialogs are created throughout the framework.
func NewDialogI(form page.FormI, id string) DialogI {
	return newDialogFunc(form, id)
}

type DialogIFuncType func(form page.FormI, id string) DialogI

var newDialogFunc DialogIFuncType = defaultNewDialogFunc // default to our built-in one

func defaultNewDialogFunc(form page.FormI, id string) DialogI {
	return NewDialog(form, id)
}

// SetNewDialogFunction sets the function that will create new dialogs.
// This is normally called by a CSS dialog implementation to set how dialogs are created in the application.
func SetNewDialogFunction(f DialogIFuncType) {
	newDialogFunc = f
}

// RestoreNewDialogFunction restores the new dialog function to the default one. This is primarily used by the example
// code, or in situations where you have multiple styles of dialog to demonstrate.
func RestoreNewDialogFunction() {
	newDialogFunc = defaultNewDialogFunc
}

// GetDialog returns the DialogI object corresponding to the given id. The id should be either the id of the
// dialog object, or the DialogPanel inside the dialog object.
func GetDialog(parent page.ControlI, id string) DialogI {
	c := parent.Page().GetControl(id)

	if dlg, ok := c.(DialogI); ok {
		return dlg
	} else if dlg, ok := c.(*DialogPanel); ok {
		return dlg.Parent().(DialogI)
	}
	return nil
}

func init() {
	page.RegisterControl(&Dialog{})
}
