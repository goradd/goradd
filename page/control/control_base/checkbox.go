package control_base

import (
	"context"
	"fmt"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/util/types"
	localPage "goradd/page"
	html2 "html"
	"strings"
)

type CheckboxI interface {
	localPage.ControlI
}


// Checkbox is a base class for checkbox-like objects, including html checkboxes and radio buttons.
type Checkbox struct {
	localPage.Control
	checked         bool
	labelMode       html.LabelDrawingMode // how to draw the label associating the text with the checkbox
	labelAttributes *html.Attributes
}

// Initializes a textbox. Normally you will not call this directly. However, sub controls should call this after
// creation to get the enclosed control initialized. Self is the newly created class. Like so:
// t := &MyTextBox{}
// t.Textbox.Init(t, parent, id)
// A parent control is isRequired. Leave id blank to have the system assign an id to the control.
func (c *Checkbox) Init(self page.ControlI, parent page.ControlI) {
	c.Control.Init(self, parent)

	c.Tag = "input"
	c.IsVoidTag = true
	c.labelMode = page.DefaultCheckboxLabelDrawingMode
	c.SetHasFor(true)
	c.SetAttribute("autocomplete", "off") // fixes an html quirk
}

func (c *Checkbox) this() CheckboxI {
	return c.Self.(CheckboxI)
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (c *Checkbox) DrawingAttributes() *html.Attributes {
	a := c.Control.DrawingAttributes()
	if c.Text() != "" && (c.labelMode == html.LABEL_BEFORE || c.labelMode == html.LABEL_AFTER) {
		// Treat the closer text label as more important than the wrapper label
		a.Set("aria-labeledby", c.ID()+"_ilbl")
	}
	return a
}

func (c *Checkbox) SetLabelDrawingMode(m html.LabelDrawingMode) {
	c.labelMode = m
	c.Refresh()
}

// Draw the checkbox tag. This can be quite tricky. Some CSS frameworks are very particular about how checkboxes get
// associated with labels. The Text value of the control will become the text directly associated with the checkbox,
// while the Label value is only shown when drawing a checkbox with a wrapper.
func (c *Checkbox) DrawTag(ctx context.Context) (ctrl string) {
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
	} else if c.labelMode == html.WRAP_LABEL_AFTER || c.labelMode == html.WRAP_LABEL_BEFORE {
		// Use the text as a label wrapper
		text = html2.EscapeString(text)
		labelAttributes := c.getDrawingInputLabelAttributes()

		if !c.HasWrapper() {
			if a := c.this().WrapperAttributes(); a != nil {
				labelAttributes.Merge(a)
			}
		}

		ctrl = html.RenderVoidTag(c.Tag, attributes)
		ctrl = html.RenderLabel(labelAttributes, text, ctrl, c.labelMode)
	} else {
		// label does not wrap. We will create a span to wrap the label and input together
		text = html2.EscapeString(text)
		labelAttributes := c.getDrawingInputLabelAttributes()

		if !c.HasWrapper() {
			if a := c.this().WrapperAttributes(); a != nil {
				labelAttributes.Merge(a)
			}
		}

		labelAttributes2 := html.NewAttributes()
		labelAttributes2.Set("for", c.ID())
		labelAttributes2.Set("id", c.ID()+"_ilbl")

		ctrl = html.RenderVoidTag(c.Tag, attributes)
		ctrl = html.RenderLabel(labelAttributes2, text, ctrl, c.labelMode)
		ctrl = html.RenderTag("span", labelAttributes, ctrl)
	}
	return ctrl
}

// Returns a pointer to the input label attributes. Feel free to set the attributes directly on the returned object.
// The input label attributes are the attributes for the label tag that associates the Text with the checkbox.
// This is specific to checkbox style controls and is not the same as the label tag that appears when using a name wrapper.
// After setting attributes, be sure to call Refresh on the control if you do this during an Ajax response.
func (c *Checkbox) InputLabelAttributes() *html.Attributes {
	if c.labelAttributes == nil {
		c.labelAttributes = html.NewAttributes()
	}
	return c.labelAttributes
}

func (c *Checkbox) getDrawingInputLabelAttributes() *html.Attributes {
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
	return a
}

// Set the value of the checkbox. Returns itself for chaining.
func (c *Checkbox) SetChecked(v bool) page.ControlI {
	if c.checked != v {
		c.checked = v
		c.AddRenderScript("prop", "checked", v)
	}

	return c.this()
}

// SetCheckedNoRefresh is used internally to update values without causing a refresh loop.
func (c *Checkbox) SetCheckedNoRefresh(v interface{}) {
	c.checked = ConvertToBool(v)
}

func (c *Checkbox) Checked() bool {
	return c.checked
}

func (c *Checkbox) SetValue(v interface{}) CheckboxI {
	c.SetChecked(ConvertToBool(v))
	return c.this()
}

func ConvertToBool(v interface{}) bool {
	var val bool
	switch s := v.(type) {
	case string:
		slower := strings.ToLower(s)
		if slower == "true" || slower == "on" || slower == "1" {
			val = true
		} else if slower == "false" || slower == "off" || slower == "" || slower == "0" {
			val = false
		} else {
			panic(fmt.Errorf("unknown checkbox string value: %s", s))
		}
	case int:
		if s == 0 {
			val = false
		} else {
			val = true
		}
	case bool:
		val = s
	default:
		panic(fmt.Errorf("unknown checkbox value: %v", v))
	}

	return val
}

func (c *Checkbox) Value() interface{} {
	return c.checked
}

/**
 * Puts the current state of the control to be able to restore it later.
 */
func (c *Checkbox) MarshalState(m types.MapI) {
	m.Set("checked", c.checked)
}

/**
 * Restore the state of the control.
 * @param mixed $state Previously saved state as returned by GetState above.
 */
func (c *Checkbox) UnmarshalState(m types.MapI) {
	if m.Has("checked") {
		v, _ := m.GetBool("checked")
		c.checked = v
	}
}
