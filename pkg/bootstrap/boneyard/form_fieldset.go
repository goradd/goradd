package boneyard

import (
"context"
	control2 "github.com/goradd/goradd/pkg/bootstrap/control"
	"github.com/goradd/goradd/pkg/html"
"github.com/goradd/goradd/pkg/log"
"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/pool"
html2 "html"
"strings"
)

type FormFieldsetI interface {
	control2.FormGroupI
}

// FormFieldset is a FormGroup kind of wrapper that is specific to using a fieldset as
// a wrapper. See https://getbootstrap.com/docs/4.3/components/forms/#horizontal-form.
// You will need to coordinate with whatever you are drawing internally to get the formatting right.
type FormFieldset struct {
	control2.FormGroup
}

func NewFormFieldset(parent page.ControlI, id string) *FormFieldset {
	p := &FormFieldset{}
	p.Init(p, parent, id)
	return p
}

func (c *FormFieldset) Init(self control.FormFieldWrapperI, parent page.ControlI, id string) {
	control2.Init(self, parent, id)
	c.Tag = "fieldset"
	c.LabelAttributes().AddClass("col-form-label").AddClass("pt-0") // helps with alignment
}

func (c *FormFieldset) this() FormFieldsetI {
	return c.Self.(FormFieldsetI)
}

func (c *FormFieldset) ΩDrawingAttributes() *html.Attributes {
	a := control2.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "formFieldset")
	return a
}

func (c *FormFieldset) ΩDrawTag(ctx context.Context) string {
	log.FrameworkDebug("Drawing FormFieldWrapper: " + c.ID())

	attributes := c.this().ΩDrawingAttributes()
	subControl := c.Page().GetControl(c.For())
	errorMessage := subControl.ValidationMessage()
	if errorMessage != "" {
		attributes.AddClass("error")
		errorMessage = html2.EscapeString(errorMessage)
	} else {
		errorMessage = "&nbsp;"
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)

	// We reuse the innerdiv to determine of we are going to wrap the inner control with a div
	// In particular, this can be used to give the control the extra row div if needed, which is
	// an idiosyncracy of the fieldset bootstrap wrapper. See the bootstrap doc for the example.
	hasInnerDiv := c.innerDivAttr.Len() > 0
	if hasInnerDiv {
		buf.WriteString("<div ")
		buf.WriteString(c.innerDivAttr.String())
		buf.WriteString(">")
	}

	text := c.Text()

	if text != "" {
		buf.WriteString(html.RenderTag("legend", c.LabelAttributes(), html2.EscapeString(text)))
	}

	var describedBy string

	describedBy = c.ID() + "_lbl"
	if errorMessage != "" {
		describedBy += " " + c.ID() + "_err"
	}
	if c.Instructions() != "" {
		describedBy += " " + c.ID() + "_inst"
	}
	describedBy = strings.TrimSpace(describedBy)

	subControl.SetAttribute("aria-describedby", describedBy)

	if err := c.this().ΩDrawInnerHtml(ctx, buf); err != nil {
		panic(err)
	}
	if subControl.ValidationState() != page.ValidationNever {
		c.ErrorAttributes().SetClass(control2.getValidationClass(subControl))
		buf.WriteString(html.RenderTag("div", c.ErrorAttributes(), errorMessage))
	}
	if c.Instructions() != "" {
		buf.WriteString(html.RenderTag("div", c.InstructionAttributes(), html2.EscapeString(c.Instructions())))
	}
	if hasInnerDiv {
		buf.WriteString("</div>")
	}
	return html.RenderTag(c.Tag, attributes, buf.String())
}

// Use FormFieldsetCreator to create a FormFieldset,
// which wraps a group of controls with a fieldset, validation error
// text and optional instructions. Pass the creator of the control you
// are wrapping as the Child item.
type FormFieldsetCreator struct {
	// ID is the control id on the html form.
	ID string
	// Label is the text that will be in the html label tag associated with the Child control.
	Label string
	// Children are the child creators declaring the controls wrapped by the fieldset
	Children []page.Creator
	// Instructions is help text that will follow the control and that further describes its purpose or use.
	Instructions string
	// LabelAttributes are additional attributes to add to the label tag.
	LabelAttributes html.AttributeCreator
	// ErrorAttributes are additional attributes to add to the tag that displays the error.
	ErrorAttributes html.AttributeCreator
	// InstructionAttributes are additional attributes to add to the tag that displays the instructions.
	InstructionAttributes html.AttributeCreator
	// Set IsInline to true to use a "span" instead of a "div" in the wrapping tag.
	IsInline bool
	// ControlOptions are additional options for the wrapper tag
	ControlOptions page.ControlOptions
	// UseTooltips will cause validation errors to be displayed with tooltips, a specific
	// feature of Bootstrap
	UseTooltips   bool
	// Set InnerRow to true to add the inner row for the fieldset
	InnerRow bool
}

func (f FormFieldsetCreator) Create(ctx context.Context, parent page.ControlI) FormFieldsetI {
	c := NewFormFieldset(parent,f.ID)
	f.Init(ctx, c)
	return c
}

func (f FormFieldsetCreator) Init(ctx context.Context, c FormFieldsetI) {
	// Reuse parent creator
	fg := control2.FormGroupCreator{
		control2.ControlOptions:        f.ControlOptions,
		control2.Label:                 f.Label,
		control2.Instructions:          f.Instructions,
		control2.LabelAttributes:       f.LabelAttributes,
		control2.ErrorAttributes:       f.ErrorAttributes,
		control2.InstructionAttributes: f.InstructionAttributes,
		control2.Child:                 f.Child,
		control2.UseTooltips:           f.UseTooltips,
	}

	if f.InnerRow {
		control2.InnerDivAttributes = html.AttributeCreator{"class": "row"}
	}

	control2.Init(ctx, c)
}