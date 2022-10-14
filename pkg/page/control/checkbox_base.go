package control

import (
	"context"
	"html"
	"io"

	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/html5tag"
)

// CheckboxI is the interface for all checkbox-like objects.
type CheckboxI interface {
	page.ControlI
	GetDrawingLabelAttributes() html5tag.Attributes
}

// CheckboxBase is a base class for checkbox-like objects, including html checkboxes and radio buttons.
type CheckboxBase struct {
	page.ControlBase
	checked bool
	// LabelMode describes where to place the label associating the text with the checkbox. The default is the
	// global page.DefaultCheckboxLabelDrawingMode, and you would normally set that instead so that all your checkboxes draw
	// the same way.
	LabelMode       html5tag.LabelDrawingMode
	labelAttributes html5tag.Attributes
}

// Init initializes a checkbox base class. It is called by checkbox implementations.
func (c *CheckboxBase) Init(parent page.ControlI, id string) {
	c.ControlBase.Init(parent, id)

	c.Tag = "input"
	c.IsVoidTag = true
	c.LabelMode = page.DefaultCheckboxLabelDrawingMode
	//c.SetHasFor(true)
	c.SetAttribute("autocomplete", "off") // fixes an html quirk
}

func (c *CheckboxBase) this() CheckboxI {
	return c.Self.(CheckboxI)
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (c *CheckboxBase) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := c.ControlBase.DrawingAttributes(ctx)
	if c.Text() != "" {
		a.AddValues("aria-labelledby", c.ID()+"_ilbl")
	}
	return a
}

// SetLabelDrawingMode determines how the label is drawn for the checkbox.
func (c *CheckboxBase) SetLabelDrawingMode(m html5tag.LabelDrawingMode) {
	c.LabelMode = m
	c.Refresh()
}

// DrawTag draws the checkbox tag. This can be quite tricky.
// Some CSS frameworks are very particular about how checkboxes get
// associated with labels. The Text value of the control will become the text directly associated with the checkbox,
// while the Label value is only shown when drawing a checkbox with a wrapper.
func (c *CheckboxBase) DrawTag(ctx context.Context, w io.Writer) {
	var ctrl string
	attributes := c.this().DrawingAttributes(ctx)
	if c.checked {
		attributes.Set("checked", "")
	}

	if text := c.Text(); text == "" {
		// there is no label to draw, just draw the input
		ctrl = html5tag.RenderVoidTag(c.Tag, attributes)
	} else if c.LabelMode == html5tag.LabelWrapAfter || c.LabelMode == html5tag.LabelWrapBefore {
		// Use the text as a label wrapper
		text = html.EscapeString(text)
		labelAttributes := c.this().GetDrawingLabelAttributes()

		labelAttributes.Set("id", c.ID()+"_ilbl")

		ctrl = html5tag.RenderVoidTag(c.Tag, attributes)
		ctrl = html5tag.RenderLabel(labelAttributes, text, ctrl, c.LabelMode)
	} else {
		// label does not wrap. We will put one after the other
		text = html.EscapeString(text)
		labelAttributes := c.this().GetDrawingLabelAttributes()

		labelAttributes.Set("for", c.ID())
		labelAttributes.Set("id", c.ID()+"_ilbl")

		ctrl = html5tag.RenderVoidTag(c.Tag, attributes)
		ctrl = html5tag.RenderLabel(labelAttributes, text, ctrl, c.LabelMode)
	}
	if _, err := io.WriteString(w, ctrl); err != nil {
		panic(err)
	}
}

// LabelAttributes returns a pointer to the input label attributes.
// Feel free to set the attributes directly on the returned object.
// The input label attributes are the attributes for the label tag that associates the Text with the checkbox.
// This is specific to checkbox style controls and is not the same as the label tag that appears when using a FormFieldWrapper wrapper.
// After setting attributes, be sure to call Refresh on the control if you do this during an Ajax response.
func (c *CheckboxBase) LabelAttributes() html5tag.Attributes {
	if c.labelAttributes == nil {
		c.labelAttributes = html5tag.NewAttributes()
	}
	return c.labelAttributes
}

// GetDrawingLabelAttributes is called by the framework to temporarily set the
// attributes of the label associated with the checkbox.
func (c *CheckboxBase) GetDrawingLabelAttributes() html5tag.Attributes {
	a := c.LabelAttributes().Copy()

	// copy tooltip to wrapping label
	if title := c.Attribute("title"); title != "" {
		a.Set("title", title)
	}

	if c.IsDisabled() {
		a.AddClass("disabled") // For styling the text associated with a disabled checkbox or control.
	}

	if !c.IsDisplayed() {
		a.SetStyle("display", "none")
	}

	a.SetData("grel", a.ID()) // make sure label gets replaced when drawing
	return a
}

// SetChecked sets the value of the checkbox. Returns itself for chaining.
func (c *CheckboxBase) SetChecked(v bool) CheckboxI {
	if c.checked != v {
		c.checked = v
		c.AddRenderScript("prop", "checked", v)
	}

	return c.this()
}

// SetCheckedNoRefresh is used internally to update values without causing a refresh loop.
func (c *CheckboxBase) SetCheckedNoRefresh(v interface{}) {
	c.checked = page.ConvertToBool(v)
}

// Checked returns true if the checkbox is checked.
func (c *CheckboxBase) Checked() bool {
	return c.checked
}

// SetValue sets the checked status of checkbox. The given value can be:
//
// For True
// "1", "true", "TRUE", "on", "ON", 1(int), true(bool)
//
// For False
// "0", "false", "FALSE", "off, "OFF", ""(empty string), 0(int), false(bool)
//
// Other values will cause a panic.
func (c *CheckboxBase) SetValue(v interface{}) CheckboxI {
	c.SetChecked(page.ConvertToBool(v))
	return c.this()
}

// Value returns the boolean checked status of the checkbox.
func (c *CheckboxBase) Value() interface{} {
	return c.checked
}

// MarshalState is called by the framework to save the state of the checkbox between form
// views. Call SetState(true) to enable state saving.
func (c *CheckboxBase) MarshalState(m page.SavedState) {
	m.Set("checked", c.checked)
}

// UnmarshalState restores the state of the checkbox if coming back to a form in the same session.
func (c *CheckboxBase) UnmarshalState(m page.SavedState) {
	if v, ok := m.Load("checked"); ok {
		if v2, ok2 := v.(bool); ok2 {
			c.checked = v2
		}
	}
}

// TextIsLabel is called by the framework to determine that the Text of the control
// is used in a label tag. You do not normally need to call this unless you are creating
// the template for a custom control.
func (c *CheckboxBase) TextIsLabel() bool {
	return true
}

// Serialize is called by the framework during pagestate serialization.
func (c *CheckboxBase) Serialize(e page.Encoder) {
	c.ControlBase.Serialize(e)

	if err := e.Encode(c.checked); err != nil {
		panic(err)
	}

	if err := e.Encode(c.LabelMode); err != nil {
		panic(err)
	}

	if err := e.Encode(c.labelAttributes); err != nil {
		panic(err)
	}
}

// Deserialize is called by the framework during page state serialization.
func (c *CheckboxBase) Deserialize(d page.Decoder) {
	c.ControlBase.Deserialize(d)

	if err := d.Decode(&c.checked); err != nil {
		panic(err)
	}

	if err := d.Decode(&c.LabelMode); err != nil {
		panic(err)
	}

	if err := d.Decode(&c.labelAttributes); err != nil {
		panic(err)
	}
}

// utility code for subclasses

// UpdateCheckboxFormValues is used by subclasses of CheckboxBase to update their internal state
// if they are a checkbox type of control.
func (c *CheckboxBase) UpdateCheckboxFormValues(ctx context.Context) {
	id := c.ID()
	grctx := page.GetContext(ctx)

	if v, ok := grctx.FormValue(id); ok {
		c.SetCheckedNoRefresh(v)
	} else if grctx.RequestMode() == page.Server && c.IsOnPage() {
		// We will not get a value if an item is not checked. But since this is a POST, all values on page
		// should send something if its checked, therefore we know its not checked.
		c.SetCheckedNoRefresh(false)
	}
}

// UpdateRadioFormValues is used by subclasses of CheckboxBase to update their internal state
// if they are a radioButton type of control.
func (c *CheckboxBase) UpdateRadioFormValues(ctx context.Context, group string) {
	id := c.ID()
	grctx := page.GetContext(ctx)

	if group != "" {
		if v, ok := grctx.FormValue(group); ok {
			c.SetCheckedNoRefresh(v == c.ID())
		}
	} else {
		// a radio button without a group makes little sense. This is here in case this is the basis for some javascript control.
		if v, ok := grctx.FormValue(id); ok {
			c.SetCheckedNoRefresh(v)
		} else if grctx.RequestMode() == page.Server && c.IsOnPage() {
			c.SetCheckedNoRefresh(false)
		}
	}
}
