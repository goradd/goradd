package control

import (
"context"
"github.com/goradd/goradd/pkg/html"
"github.com/goradd/goradd/pkg/log"
"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/pool"
html2 "html"
"strings"
)

type FormFieldsetI interface {
	FormGroupI
}

// FormFieldset is FormGroup kind of wrapper that is specific to using a fieldset as
// a wrapper. See https://getbootstrap.com/docs/4.3/components/forms/#horizontal-form.
// You will need to coordinate with whatever you are drawing internally to get the formatting right.
type FormFieldset struct {
	FormGroup
}

func NewFormFieldset(parent page.ControlI, id string) *FormFieldset {
	p := &FormFieldset{}
	p.Init(p, parent, id)
	return p
}

func (c *FormFieldset) Init(self control.FormFieldI, parent page.ControlI, id string) {
	c.FormGroup.Init(self, parent, id)
	c.Tag = "fieldset"
	c.LabelAttributes().AddClass("col-form-label").AddClass("pt-0") // helps with alignment
}

func (c *FormFieldset) this() FormFieldsetI {
	return c.Self.(FormFieldsetI)
}

func (c *FormFieldset) ΩDrawingAttributes() *html.Attributes {
	a := c.FormGroup.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "formFieldset")
	return a
}

func (c *FormFieldset) ΩDrawTag(ctx context.Context) string {
	log.FrameworkDebug("Drawing FormField: " + c.ID())

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
		c.ErrorAttributes().SetClass(c.getValidationClass(subControl))
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


type FormFieldsetCreator struct {
	ControlOptions page.ControlOptions
	ID string
	Label string
	For string
	Instructions string
	LabelAttributes html.AttributeCreator
	ErrorAttributes html.AttributeCreator
	InstructionAttributes html.AttributeCreator
	Child page.Creator

	// Set InnerRow to true to add the inner row for the fieldset
	InnerRow bool
	UseTooltips   bool // uses tooltips for the error class
}

func (f FormFieldsetCreator) Create(ctx context.Context, parent page.ControlI) FormFieldsetI {
	c := NewFormFieldset(parent, f.ID)
	f.Init(ctx, c)
	return c
}

func (f FormFieldsetCreator) Init(ctx context.Context, c FormFieldsetI) {
	// Reuse parent creator
	fg := FormGroupCreator {
		ControlOptions: f.ControlOptions,
		Label: f.Label,
		For: f.For,
		Instructions: f.Instructions,
		LabelAttributes: f.LabelAttributes,
		ErrorAttributes: f.ErrorAttributes,
		InstructionAttributes: f.InstructionAttributes,
		Child: f.Child,
		UseTooltips: f.UseTooltips,
	}

	if f.InnerRow {
		fg.InnerDivAttributes = html.AttributeCreator{"class":"row"}
	}

	fg.Init(ctx, c)
}