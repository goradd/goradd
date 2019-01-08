package control

import (
	"context"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	html2 "html"
	"reflect"
)

type CheckboxI interface {
	page.ControlI
	GetDrawingInputLabelAttributes() *html.Attributes
}


// Checkbox is a base class for checkbox-like objects, including html checkboxes and radio buttons.
type CheckboxBase struct {
	page.Control
	checked         bool
	LabelMode       html.LabelDrawingMode // how to draw the label associating the text with the checkbox
	labelAttributes *html.Attributes
}

// Init initializes a checbox base class. Normally you will not call this directly. However, sub controls should call this after
// creation to get the enclosed control initialized.
func (c *CheckboxBase) Init(self page.ControlI, parent page.ControlI, id string) {
	c.Control.Init(self, parent, id)

	c.Tag = "input"
	c.IsVoidTag = true
	c.LabelMode = page.DefaultCheckboxLabelDrawingMode
	c.SetHasFor(true)
	c.SetAttribute("autocomplete", "off") // fixes an html quirk
}

func (c *CheckboxBase) this() CheckboxI {
	return c.Self.(CheckboxI)
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (c *CheckboxBase) DrawingAttributes() *html.Attributes {
	a := c.Control.DrawingAttributes()
	if c.Text() != "" && (c.LabelMode == html.LabelBefore || c.LabelMode == html.LabelAfter) {
		// Treat the closer text label as more important than the wrapper label
		a.Set("aria-labeledby", c.ID()+"_ilbl")
	}
	return a
}

func (c *CheckboxBase) SetLabelDrawingMode(m html.LabelDrawingMode) {
	c.LabelMode = m
	c.Refresh()
}

// Draw the checkbox tag. This can be quite tricky. Some CSS frameworks are very particular about how checkboxes get
// associated with labels. The Text value of the control will become the text directly associated with the checkbox,
// while the Label value is only shown when drawing a checkbox with a wrapper.
func (c *CheckboxBase) DrawTag(ctx context.Context) (ctrl string) {
	attributes := c.this().DrawingAttributes()
	if c.checked {
		attributes.Set("checked", "")
	}

	if text := c.Text(); text == "" {
		// there is no label to draw, just draw the input
		if !c.HasWrapper() {
			if a := c.this().WrapperAttributes(); a != nil {
				attributes.Merge(a)
			}
		}
		ctrl = html.RenderVoidTag(c.Tag, attributes)
	} else if c.LabelMode == html.LabelWrapAfter || c.LabelMode == html.LabelWrapBefore {
		// Use the text as a label wrapper
		text = html2.EscapeString(text)
		labelAttributes := c.this().GetDrawingInputLabelAttributes()

		if !c.HasWrapper() {
			if a := c.this().WrapperAttributes(); a != nil {
				labelAttributes.Merge(a)
			}
		}

		ctrl = html.RenderVoidTag(c.Tag, attributes)
		ctrl = html.RenderLabel(labelAttributes, text, ctrl, c.LabelMode)
	} else {
		// label does not wrap. We will put one after the other
		text = html2.EscapeString(text)
		labelAttributes := c.this().GetDrawingInputLabelAttributes()

		if !c.HasWrapper() {
			if a := c.this().WrapperAttributes(); a != nil {
				labelAttributes.Merge(a)
			}
		}

		labelAttributes.Set("for", c.ID())
		labelAttributes.Set("id", c.ID()+"_ilbl")

		ctrl = html.RenderVoidTag(c.Tag, attributes)
		ctrl = html.RenderLabel(labelAttributes, text, ctrl, c.LabelMode)
	}
	return ctrl
}

// Returns a pointer to the input label attributes. Feel free to set the attributes directly on the returned object.
// The input label attributes are the attributes for the label tag that associates the Text with the checkbox.
// This is specific to checkbox style controls and is not the same as the label tag that appears when using a name wrapper.
// After setting attributes, be sure to call Refresh on the control if you do this during an Ajax response.
func (c *CheckboxBase) InputLabelAttributes() *html.Attributes {
	if c.labelAttributes == nil {
		c.labelAttributes = html.NewAttributes()
	}
	return c.labelAttributes
}

func (c *CheckboxBase) GetDrawingInputLabelAttributes() *html.Attributes {
	a := c.InputLabelAttributes().Clone()

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

	a.SetDataAttribute("grel", a.ID())	// make sure label gets replaced when drawing
	return a
}

// Set the value of the checkbox. Returns itself for chaining.
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

func (c *CheckboxBase) Checked() bool {
	return c.checked
}

func (c *CheckboxBase) SetValue(v interface{}) CheckboxI {
	c.SetChecked(page.ConvertToBool(v))
	return c.this()
}


func (c *CheckboxBase) Value() interface{} {
	return c.checked
}

/**
 * Puts the current state of the control to be able to restore it later.
 */
func (c *CheckboxBase) MarshalState(m maps.Setter) {
	m.Set("checked", c.checked)
}

/**
 * Restore the state of the control.
 * @param mixed $state Previously saved state as returned by GetState above.
 */
func (c *CheckboxBase) UnmarshalState(m maps.Loader) {
	if v,ok := m.Load("checked"); ok {
		if v2,ok := v.(bool); ok {
			c.checked = v2
		}
	}
}

func (c *CheckboxBase) TextIsLabel() bool {
	return true
}

func (c *CheckboxBase) Serialize(e page.Encoder) (err error) {
	if err = c.Control.Serialize(e); err != nil {
		return
	}

	if err = e.Encode(c.checked); err != nil {
		return err
	}

	if err = e.Encode(c.LabelMode); err != nil {
		return err
	}

	if err = e.Encode(c.labelAttributes); err != nil {
		return err
	}

	return
}

// ΩisSerializer is used by the automated control serializer to determine how far down the control chain the control
// has to go before just calling serialize and deserialize
func (c *CheckboxBase) ΩisSerializer(i page.ControlI) bool {
	return reflect.TypeOf(c) == reflect.TypeOf(i)
}


func (c *CheckboxBase) Deserialize(d page.Decoder, p *page.Page) (err error) {
	if err = c.Control.Deserialize(d, p); err != nil {
		return
	}

	if err = d.Decode(&c.checked); err != nil {
		return
	}

	if err = d.Decode(&c.LabelMode); err != nil {
		return
	}

	if err = d.Decode(&c.labelAttributes); err != nil {
		return
	}

	return
}
