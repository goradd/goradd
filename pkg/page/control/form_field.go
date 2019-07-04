package control

import (
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/pool"
	html2 "html"
	"strings"
)

type FormFieldI interface {
	page.ControlI
	SetFor(relatedId string) FormFieldI
	For() string
	Instructions() string
	SetInstructions(string) FormFieldI
	LabelAttributes() *html.Attributes
	ErrorAttributes() *html.Attributes
	InstructionAttributes() *html.Attributes
}

// FormField is a Goradd control that wraps other controls, and provides common companion
// functionality like a form label, validation state display, and help text.
type FormField struct {
	page.Control

	// instructions is text associated with the control for extra explanation. You could also try adding a tooltip to the wrapper.
	instructions string
	// labelAttributes are the attributes that will be directly put on the Label tag. The label tag itself comes
	// from the "Text" item in the control.
	labelAttributes *html.Attributes
	errorAttributes *html.Attributes
	instructionAttributes *html.Attributes
	forID string
}

func NewFormField(parent page.ControlI, id string) *FormField {
	p := &FormField{}
	p.Init(p, parent, id)
	return p
}

func (c *FormField) Init(self PanelI, parent page.ControlI, id string) {
	c.Control.Init(self, parent, id)
	c.Tag = "div"
	c.labelAttributes = html.NewAttributes().
		SetID(c.ID() + "_lbl").
		SetClass("goradd-lbl")
	c.errorAttributes = html.NewAttributes().
		SetID(c.ID() + "_err").
		SetClass("goradd-error")
	c.instructionAttributes = html.NewAttributes().
		SetID(c.ID() + "_inst").
		SetClass("goradd-instructions")

}

func (c *FormField) this() FormFieldI {
	return c.Self.(FormFieldI)
}

// SetFor associates the form field with a sub control. The relatedId
// is the ID that the form field is associated with. Most browsers allow you to click on the
// label in order to give focus to the related control
func (c *FormField) SetFor(relatedId string) FormFieldI {
	if relatedId == "" {
		panic("A For id is required.")
	}
	c.forID = relatedId
	return c.this()
}

func (c *FormField) For() string {
	return c.forID
}

// SetInstructions sets the instructions that will be printed with the control. Instructions only get rendered
// by wrappers, so if there is no wrapper, or the wrapper does not render the instructions, this will not appear.
func (c *FormField) SetInstructions(i string) FormFieldI {
	if i != c.instructions {
		c.instructions = i
		c.Refresh()
	}
	return c.this()
}

// Instructions returns the instructions to be printed with the control
func (c *FormField) Instructions() string {
	return c.instructions
}

func (c *FormField) ΩDrawingAttributes() *html.Attributes {
	a := c.Control.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "formField")
	return a
}

func (c *FormField) ΩDrawTag(ctx context.Context) string {
	log.FrameworkDebug("Drawing FormField: " + c.ID())

	attributes := c.this().ΩDrawingAttributes()
	subControl := c.Page().GetControl(c.forID)
	errorMessage := subControl.ValidationMessage()
	if errorMessage != "" {
		attributes.AddClass("error")
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)

	text := c.Text()
	if text == "" {
		text = subControl.Attribute("placeholder")
		if text != "" {
			c.labelAttributes.SetStyle("display", "none") // make a hidden label for screen readers
		}
	}

	if text != "" {
		c.labelAttributes.Set("for", c.forID)
		buf.WriteString(html.RenderTag("label", c.labelAttributes, html2.EscapeString(text)))
	}

	var describedBy string

	if errorMessage != "" {
		describedBy = c.ID() + "_err"
	}
	if c.instructions != "" {
		describedBy += " " + c.ID() + "_inst"
	}
	describedBy = strings.TrimSpace(describedBy)

	if describedBy != "" {
		subControl.SetAttribute("aria-describedby", describedBy)
	}
	if err := c.this().ΩDrawInnerHtml(ctx, buf); err != nil {
		panic(err)
	}
	if subControl.ValidationState() != page.ValidationNever {
		buf.WriteString(html.RenderTag("div", c.errorAttributes, html2.EscapeString(errorMessage)))
	}
	if c.instructions != "" {
		buf.WriteString(html.RenderTag("div", c.instructionAttributes, html2.EscapeString(c.instructions)))
	}
	return html.RenderTag(c.Tag, attributes, buf.String())
}

func (c *FormField) LabelAttributes() *html.Attributes {
	return c.labelAttributes
}

func (c *FormField) SetLabelAttributes(a *html.Attributes) FormFieldI {
	c.labelAttributes = a
	return c.this()
}

func (c *FormField) ErrorAttributes() *html.Attributes {
	return c.errorAttributes
}

func (c *FormField) SetErrorAttributes(a *html.Attributes) FormFieldI {
	c.errorAttributes = a
	return c.this()
}

func (c *FormField) InstructionAttributes() *html.Attributes {
	return c.instructionAttributes
}

func (c *FormField) SetInstructionAttributes(a *html.Attributes) FormFieldI {
	c.instructionAttributes = a
	return c.this()
}

type FormFieldCreator struct {
	ID string
	Label string
	For string
	IsInline bool
	Instructions string
	LabelAttributes html.AttributeCreator
	ErrorAttributes html.AttributeCreator
	InstructionAttributes html.AttributeCreator
	Child page.Creator
	ControlOptions page.ControlOptions
}

func (f FormFieldCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	c := NewFormField(parent, f.ID)
	f.Init(ctx, c)
	if f.IsInline { // subclasses might deal with this issue differently
		c.Tag = "span"
	}

	return c
}

func (f FormFieldCreator) Init(ctx context.Context, c FormFieldI) {
	c.ApplyOptions(f.ControlOptions)
	c.SetText(f.Label)
	c.SetInstructions(f.Instructions)
	c.SetFor(f.For)
	if f.LabelAttributes != nil {
		c.LabelAttributes().Merge(html.NewAttributesFromMap(f.LabelAttributes))
	}
	if f.ErrorAttributes != nil {
		c.ErrorAttributes().Merge(html.NewAttributesFromMap(f.ErrorAttributes))
	}
	if f.InstructionAttributes != nil {
		c.InstructionAttributes().Merge(html.NewAttributesFromMap(f.InstructionAttributes))
	}

	c.AddControls(ctx, f.Child)
}