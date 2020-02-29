package control

import (
	"context"
	"github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
)

type RadioListGroupI interface {
	RadioListI
	SetButtonStyle(string) RadioListGroupI
}

// RadioListGroup is a RadioList styled as a group.
//
// See https://getbootstrap.com/docs/4.4/components/buttons/#checkbox-and-radio-buttons
type RadioListGroup struct {
	RadioList
	buttonStyle string
}

func NewRadioListGroup(parent page.ControlI, id string) *RadioListGroup {
	l := new(RadioListGroup)
	l.Self = l
	l.Init(parent, id)
	return l
}

func (l *RadioListGroup) Init(parent page.ControlI, id string) {
	l.RadioList.Init(parent, id)
	l.SetLabelDrawingMode(html.LabelWrapAfter)
	l.SetRowClass("")
	l.buttonStyle = ButtonStyleSecondary
	config.LoadBootstrap(l.ParentForm())
}

func (l *RadioListGroup) this() RadioListGroupI {
	return l.Self.(RadioListGroupI)
}

func (l *RadioListGroup) SetButtonStyle(buttonStyle string) RadioListGroupI {
	l.buttonStyle = buttonStyle
	return l
}


// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *RadioListGroup) DrawingAttributes(ctx context.Context) html.Attributes {
	a := l.ControlBase.DrawingAttributes(ctx) // skip default checkbox list attributes
	a.SetDataAttribute("grctl", "bs-RadioListGroup")
	a.AddClass("btn-group btn-group-toggle")
	a.SetDataAttribute("toggle", "buttons")
	return a
}

// RenderItem is called by the framework to render a single item in the list.
func (l *RadioListGroup) RenderItem(item *control.ListItem) (h string) {
	selected := l.SelectedItem().ID() == item.ID()
	attributes := html.NewAttributes()
	attributes.SetID(item.ID())
	attributes.Set("name", l.ID())
	attributes.Set("value", item.Value())
	attributes.Set("type", "radio")
	if selected {
		attributes.Set("checked", "")
	}
	ctrl := html.RenderVoidTag("input", attributes)
	labelAttributes := html.NewAttributes().Set("for", item.ID()).AddClass("btn").AddClass(l.buttonStyle)
	if selected {
		labelAttributes.AddClass("active")
	}
	return html.RenderLabel(labelAttributes, item.Label(), ctrl, html.LabelWrapAfter)
}

func (l *RadioListGroup) Serialize(e page.Encoder) (err error) {
	if err = l.RadioList.Serialize(e); err != nil {
		return
	}

	if err = e.Encode(l.buttonStyle); err != nil {
		return err
	}

	return
}


func (l *RadioListGroup) Deserialize(d page.Decoder) (err error) {
	if err = l.RadioList.Deserialize(d); err != nil {
		return
	}

	if err = d.Decode(&l.buttonStyle); err != nil {
		return
	}

	return
}


type RadioListGroupCreator struct {
	ID string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []control.ListValue
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider control.DataBinder
	// DataProviderID is the id of a control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProviderID string
	// Value is the initial value of the textbox. Often its best to load the value in a separate Load step after creating the control.
	Value string
	// SaveState saves the selected value so that it is restored if the form is returned to.
	ButtonStyle string
	OnChange action.ActionI
	SaveState bool
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c RadioListGroupCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewRadioListGroup(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c RadioListGroupCreator) Init(ctx context.Context, ctrl RadioListGroupI) {
	sub := RadioListCreator{
		ID: c.ID,
		Items: c.Items,
		DataProvider: c.DataProvider,
		Value: c.Value,
		SaveState: c.SaveState,
		ControlOptions: c.ControlOptions,
	}
	sub.Init(ctx, ctrl)
	if c.ButtonStyle != "" {
		ctrl.SetButtonStyle(c.ButtonStyle)
	}
	if c.OnChange != nil {
		ctrl.On(event.Change(), c.OnChange)
	}
}

// GetRadioList is a convenience method to return the control with the given id from the page.
func GetRadioListGroup(c page.ControlI, id string) *RadioListGroup {
	return c.Page().GetControl(id).(*RadioListGroup)
}

func init() {
	page.RegisterControl(&RadioListGroup{})
}