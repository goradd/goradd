package control

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/pool"
	html2 "html"
)

type FormFieldsetI interface {
	control.PanelI
	LegendAttributes() html.Attributes
	SetAsRow(r bool) FormFieldsetI
	SetInstructions(instructions string) FormFieldsetI
}

// FormFieldset is a FormGroup kind of wrapper that is specific to using a fieldset as
// a wrapper. See https://getbootstrap.com/docs/4.3/components/forms/#horizontal-form.
// You will need to coordinate with whatever you are drawing internally to get the formatting right.
type FormFieldset struct {
	control.Panel
	legendAttributes html.Attributes
	asRow bool
	instructions string
	instructionAttributes html.Attributes
}

func NewFormFieldset(parent page.ControlI, id string) *FormFieldset {
	p := &FormFieldset{}
	p.Self = p
	p.Init(parent, id)
	return p
}

func (c *FormFieldset) Init(parent page.ControlI, id string) {
	c.Panel.Init(parent, id)
	c.Tag = "fieldset"
	c.legendAttributes = html.NewAttributes()
	c.legendAttributes.AddClass("pt-0") // helps with alignment. Remove if needed
	c.instructionAttributes = html.NewAttributes().
		SetID(c.ID() + "_inst").
		SetClass("form-text")
}

func (c *FormFieldset) this() FormFieldsetI {
	return c.Self.(FormFieldsetI)
}

func (c *FormFieldset) LegendAttributes() html.Attributes {
	return c.legendAttributes
}

func (c *FormFieldset) SetAsRow(r bool) FormFieldsetI{
	c.asRow = r
	return c.this()
}

func (c *FormFieldset) SetInstructions(instructions string) FormFieldsetI {
	if c.instructions != instructions {
		c.instructions = instructions
		c.Refresh()
	}
	return c.this()
}

func (c *FormFieldset) InstructionAttributes() html.Attributes {
	return c.instructionAttributes
}


func (c *FormFieldset) DrawingAttributes(ctx context.Context) html.Attributes {
	a := c.Panel.DrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "formFieldset")
	return a
}

func (c *FormFieldset) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	var s string

	buf2 := pool.GetBuffer()
	defer pool.PutBuffer(buf2)

	if c.Text() != "" {
		buf2.WriteString(html.RenderTag("legend", c.legendAttributes, html2.EscapeString(c.Text())))
	}
	if err = c.Panel.DrawInnerHtml(ctx, buf2); err != nil {
		return
	}
	if c.instructions != "" {
		s = html.RenderTag("small", c.instructionAttributes, html2.EscapeString(c.instructions))
		buf2.WriteString(s)
	}

	if c.asRow {
		s = html.RenderTag("div", html.NewAttributes().AddClass("row"), buf2.String())
		buf.WriteString(s)
	} else {
		_,err = buf2.WriteTo(buf)
	}
	return
}

func (c *FormFieldset) Serialize(e page.Encoder) (err error) {
	if err = c.Panel.Serialize(e); err != nil {
		return
	}

	if err = e.Encode(c.legendAttributes); err != nil {
		return err
	}
	if err = e.Encode(c.asRow); err != nil {
		return err
	}
	if err = e.Encode(c.instructions); err != nil {
		return err
	}
	if err = e.Encode(c.instructionAttributes); err != nil {
		return err
	}

	return
}


func (c *FormFieldset) Deserialize(d page.Decoder) (err error) {
	if err = c.Panel.Deserialize(d); err != nil {
		return
	}

	if err = d.Decode(&c.legendAttributes); err != nil {
		return
	}
	if err = d.Decode(&c.asRow); err != nil {
		return
	}
	if err = d.Decode(&c.instructions); err != nil {
		return
	}
	if err = d.Decode(&c.instructionAttributes); err != nil {
		return
	}

	return
}


// Use FormFieldsetCreator to create a bootstrap fieldset,
// which wraps a control group with a fieldset.
// The Child item should be a panel or a control that groups other controls,
// like a RadioList or CheckboxList
type FormFieldsetCreator struct {
	// ID is the control id on the html form.
	ID string
	// Legend is the text that will be in the html label tag associated with the Child control.
	Legend string
	// Child is should be a panel, or a control that draws a group of controls,
	// like a RadioList or CheckboxList
	Child page.Creator
	// LegendAttributes are additional attributes to add to the label tag.
	LegendAttributes html.Attributes
	// Instructions is help text that accompanies the control
	Instructions string
	// Set AsRow to true to put the legend on the same row as the content
	AsRow bool
	// ControlOptions are additional options for the wrapper tag
	ControlOptions page.ControlOptions
}

func (f FormFieldsetCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	id := control.CalcWrapperID(f.ID, f.Child, "fs")
	c := NewFormFieldset(parent,id)
	f.Init(ctx, c)
	return c
}

func (f FormFieldsetCreator) Init(ctx context.Context, c FormFieldsetI) {
	if f.Legend != "" {
		c.SetText(f.Legend)
	}
	if f.LegendAttributes != nil {
		c.LegendAttributes().Merge(f.LegendAttributes)
	}
	if f.AsRow {
		c.SetAsRow(true)
	}
	if f.Instructions != "" {
		c.SetInstructions(f.Instructions)
	}
	c.AddControls(ctx, f.Child)
	c.ApplyOptions(ctx, f.ControlOptions)
}

// GetFormFieldset is a convenience method to return the fieldset with the given id from the page.
func GetFormFieldset(c page.ControlI, id string) *FormFieldset {
	return c.Page().GetControl(id).(*FormFieldset)
}

func init() {
	page.RegisterControl(&FormFieldset{})
}