package control
import (
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	localPage "goradd/page"
	"github.com/spekary/goradd/util/types"
	"fmt"
	"strings"
	"bytes"
	"context"
	"goradd/config"
	html2 "html"
)

type Checkbox struct {
	checkboxBase
}

func (c *Checkbox) DrawingAttributes() *html.Attributes {
	a := c.checkboxBase.DrawingAttributes()
	a.Set("name", c.Id())	// needed for posts
	a.Set("type", "checkbox")
	return a
}

// UpdateFormValues is an internal call that lets us reflect the value of the checkbox on the web page
func (c *Checkbox) UpdateFormValues(ctx *page.Context) {
	id := c.Id()

	if v,ok := ctx.CheckableValue(id); ok {
		c.checked = c.convertToBool(v)
	}
}

type RadioButton struct {
	checkboxBase
	group string
}

func (c *RadioButton) DrawingAttributes() *html.Attributes {
	a := c.checkboxBase.DrawingAttributes()
	a.Set("type", "radio")
	if c.group == "" {
		a.Set("name", c.Id())	// treat it like a checkbox if no group is specified
	} else {
		a.Set("name", c.group)
		a.Set("value", c.Id())
	}
	return a
}

func (c *RadioButton) UpdateFormValues(ctx *page.Context) {
	id := c.Id()

	if v,ok := ctx.CheckableValue(id); ok {
		c.checked = c.convertToBool(v)
	}
}

func (c *RadioButton) SetGroup (g string) page.ControlI {
	c.group = g
	c.Refresh()
	return c.This()
}

func (c *RadioButton) Group() string {
	return c.group
}

func (c *RadioButton) SetChecked(v bool) page.ControlI {
	if c.group != "" && v {
		if c.checked != v {
			c.checked = v
			c.Form().Response().ExecuteJsFunction("goradd.setRadioInGroup", page.PriorityStandard, c.Id());
		}
	} else {
		if c.checked != v {
			c.checked = v
			c.AddRenderScript("prop", "checked", v)
		}
	}
	return c.This()
}


type checkboxBase struct {
	localPage.Control
	checked bool
	labelMode	html.LabelDrawingMode		// how to draw the label associating the text with the checkbox
	labelAttributes *html.Attributes
}


// Initializes a textbox. Normally you will not call this directly. However, sub controls should call this after
// creation to get the enclosed control initialized. Self is the newly created class. Like so:
// t := &MyTextBox{}
// t.Textbox.Init(t, parent, id)
// A parent control is isRequired. Leave id blank to have the system assign an id to the control.
func (c *checkboxBase) Init(self TextboxI, parent page.ControlI) {
	c.Control.Init(self, parent)

	c.Tag = "input"
	c.IsVoidTag = true
	c.labelMode = config.DefaultCheckboxLabelDrawingMode
	c.SetHasFor(true)
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (c *checkboxBase) DrawingAttributes() *html.Attributes {
	a := c.Control.DrawingAttributes()
	if c.Text() != "" && (c.labelMode == html.LABEL_BEFORE || c.labelMode == html.LABEL_AFTER) {
		// Treat the closer text label as more important than the wrapper label
		a.Set("aria-labeledby", c.Id() + "_ilbl")
	}
	return a
}

func (c *checkboxBase) SetLabelDrawingMode(m html.LabelDrawingMode) {
	c.labelMode = m
	c.Refresh()
}

// Draw the checkbox tag. This can be quite tricky. Some CSS frameworks are very particular about how checkboxes get
// associated with labels. The Text value of the control will become the text directly associated with the checkbox,
// while the Label value is only shown when drawing a checkbox with a wrapper.
func (c *checkboxBase) DrawTag(ctx context.Context, buf *bytes.Buffer) (ctrl string) {
	attributes := c.This().DrawingAttributes()
	if c.checked {
		attributes.Set("checked", "")
	}

	if text := c.Text(); text == "" {
		// there is no label to draw, just draw the input
		if !c.HasWrapper() {
			if a := c.This().WrapperAttributes(); a != nil {
				attributes.Merge(a)
			}
		}
		ctrl = html.RenderVoidTag(c.Tag, attributes)
	} else if c.labelMode == html.WRAP_LABEL_AFTER || c.labelMode == html.WRAP_LABEL_BEFORE {
		// Use the text as a label wrapper
		text = html2.EscapeString(text)
		labelAttributes := c.getDrawingInputLabelAttributes()

		if !c.HasWrapper() {
			if a := c.This().WrapperAttributes(); a != nil {
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
			if a := c.This().WrapperAttributes(); a != nil {
				labelAttributes.Merge(a)
			}
		}

		labelAttributes2 := html.NewAttributes()
		labelAttributes2.Set("for", c.Id())
		labelAttributes2.Set("id", c.Id() + "_ilbl")

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
func (c *checkboxBase) InputLabelAttributes() *html.Attributes {
	if c.labelAttributes == nil {
		c.labelAttributes = html.NewAttributes()
	}
	return c.labelAttributes
}

func (c *checkboxBase) getDrawingInputLabelAttributes() *html.Attributes {
	a := c.InputLabelAttributes().Clone()

	// copy tooltip to wrapping label
	if title := c.Attribute("title"); title != "" {
		a.Set("title", title)
	}

	if c.IsDisabled() {
		a.AddClass("disabled")	// For styling the text associated with a disabled checkbox or control.
	}

	if !c.IsDisplayed() {
		a.SetStyle("display", "none")
	}
	return a
}


// Set the value of the checkbox. Returns itself for chaining.
func (c *checkboxBase) SetChecked(v bool) page.ControlI {
	if c.checked != v {
		c.checked = v
		c.AddRenderScript("prop", "checked", v)
	}

	return c.This()
}

func (c *checkboxBase) Checked() bool {
	return c.checked
}

func (c *checkboxBase) SetValue(v interface{}) page.ControlI {
	c.SetChecked(c.convertToBool(v))
	return c.Self.(page.ControlI)
}

func (c *checkboxBase) convertToBool(v interface{}) bool {
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

func (c *checkboxBase) Value() interface{} {
	return c.checked
}

/**
 * Puts the current state of the control to be able to restore it later.
 */
func (c *checkboxBase) MarshalState(m types.MapI) {
	m.Set("checked", c.checked)
}

/**
 * Restore the state of the control.
 * @param mixed $state Previously saved state as returned by GetState above.
 */
func (c *checkboxBase) UnmarshalState(m types.MapI) {
	if m.Has("checked") {
		v,_ := m.GetBool("checked")
		c.checked = v
	}
}
