package control

import (
	"context"
	"github.com/goradd/goradd/pkg/html5tag"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/pool"
	"html"
	"io"
)

type FormFieldsetI interface {
	control.PanelI
	LegendAttributes() html5tag.Attributes
	SetAsRow(r bool) FormFieldsetI
	SetInstructions(instructions string) FormFieldsetI
}

// FormFieldset is a FormGroup kind of wrapper that is specific to using a fieldset as
// a wrapper. See https://getbootstrap.com/docs/4.3/components/forms/#horizontal-form.
// You will need to coordinate with whatever you are drawing internally to get the formatting right.
type FormFieldset struct {
	control.Panel
	legendAttributes      html5tag.Attributes
	asRow                 bool
	instructions          string
	instructionAttributes html5tag.Attributes
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
	c.legendAttributes = html5tag.NewAttributes()
	c.legendAttributes.AddClass("pt-0") // helps with alignment. Remove if needed
	c.instructionAttributes = html5tag.NewAttributes().
		SetID(c.ID() + "_inst").
		SetClass("form-text")
}

func (c *FormFieldset) this() FormFieldsetI {
	return c.Self.(FormFieldsetI)
}

func (c *FormFieldset) LegendAttributes() html5tag.Attributes {
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

func (c *FormFieldset) InstructionAttributes() html5tag.Attributes {
	return c.instructionAttributes
}


func (c *FormFieldset) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := c.Panel.DrawingAttributes(ctx)
	a.SetData("grctl", "formFieldset")
	return a
}

func (c *FormFieldset) DrawInnerHtml(ctx context.Context, w io.Writer) {
	var s string

	buf2 := pool.GetBuffer()
	defer pool.PutBuffer(buf2)

	if c.Text() != "" {
		buf2.WriteString(html5tag.RenderTag("legend", c.legendAttributes, html.EscapeString(c.Text())))
	}
	c.Panel.DrawInnerHtml(ctx, buf2)
	if c.instructions != "" {
		s = html5tag.RenderTag("small", c.instructionAttributes, html.EscapeString(c.instructions))
		buf2.WriteString(s)
	}

	if c.asRow {
		s = html5tag.RenderTag("div", html5tag.NewAttributes().AddClass("row"), buf2.String())
		page.WriteString(w, s)
	} else {
		page.WriteString(w, buf2.String())
	}
	return
}

func (c *FormFieldset) Serialize(e page.Encoder) {
	c.Panel.Serialize(e)

	if err := e.Encode(c.legendAttributes); err != nil {
		panic(err)
	}
	if err := e.Encode(c.asRow); err != nil {
		panic(err)
	}
	if err := e.Encode(c.instructions); err != nil {
		panic(err)
	}
	if err := e.Encode(c.instructionAttributes); err != nil {
		panic(err)
	}
}


func (c *FormFieldset) Deserialize(d page.Decoder) {
	c.Panel.Deserialize(d)

	if err := d.Decode(&c.legendAttributes); err != nil {
		panic(err)
	}
	if err := d.Decode(&c.asRow); err != nil {
		panic(err)
	}
	if err := d.Decode(&c.instructions); err != nil {
		panic(err)
	}
	if err := d.Decode(&c.instructionAttributes); err != nil {
		panic(err)
	}
}


// FormFieldsetCreator creates a bootstrap fieldset,
// which wraps a control group with a fieldset.
// The Child item should be a panel or a control that groups other controls,
// like a RadioList or CheckboxList
type FormFieldsetCreator struct {
	// ID is the control id on the html form.
	ID string
	// Legend is the text that will be in the html label tag associated with the Child control.
	Legend string
	// Child should be a panel, or a control that draws a group of controls,
	// like a RadioList or CheckboxList
	Child page.Creator
	// LegendAttributes are additional attributes to add to the label tag.
	LegendAttributes html5tag.Attributes
	// Instructions contains help text that accompanies the control
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