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

type FormGroupI interface {
	control.FormFieldI
	SetUseTooltips(use bool) FormGroupI
	UseTooltips() bool
	InnerDivAttributes() *html.Attributes
}

// FormGroup is a Goradd control that wraps other controls, and provides common companion
// functionality like a form label, validation state display, and help text.
type FormGroup struct {
	control.FormField
	innerDivAttr *html.Attributes
	useTooltips   bool // uses tooltips for the error class
}

func NewFormGroup(parent page.ControlI, id string) *FormGroup {
	p := &FormGroup{}
	p.Init(p, parent, id)
	return p
}

func (c *FormGroup) Init(self control.FormFieldI, parent page.ControlI, id string) {
	c.FormField.Init(self, parent, id)
	c.innerDivAttr = html.NewAttributes()
	c.AddClass("form-group") // to get a wrapper with out this, just remove it after initialization
	c.InstructionAttributes().AddClass("form-text")
}

func (c *FormGroup) this() FormGroupI {
	return c.Self.(FormGroupI)
}

// SetUseTooltips sets whether to use tooltips to display validation messages.
func (c *FormGroup) SetUseTooltips(use bool) FormGroupI {
	c.useTooltips = use
	return c
}

func (c *FormGroup) UseTooltips() bool {
	return c.useTooltips
}

func (c *FormGroup) ΩDrawingAttributes() *html.Attributes {
	a := c.FormField.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "formGroup")
	if c.useTooltips {
		// bootstrap requires that parent of a tool-tipped object has position relative
		c.SetStyle("position", "relative")
	}
	return a
}

func (c *FormGroup) ΩDrawTag(ctx context.Context) string {
	log.FrameworkDebug("Drawing FormField: " + c.ID())

	attributes := c.this().ΩDrawingAttributes()
	if c.For() == "" {
		panic("a FormGroup MUST have a sub control")
	}
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

	text := c.Text()
	if text == "" {
		text = subControl.Attribute("placeholder")
		if text != "" {
			c.LabelAttributes().SetClass("sr-only") // make a hidden label for screen readers
		}
	}

	if text != "" {
		c.LabelAttributes().Set("for", c.For())
		buf.WriteString(html.RenderTag("label", c.LabelAttributes(), html2.EscapeString(text)))
	}

	var describedBy string

	if errorMessage != "" {
		describedBy = c.ID() + "_err"
	}
	if c.Instructions() != "" {
		describedBy += " " + c.ID() + "_inst"
	}
	describedBy = strings.TrimSpace(describedBy)

	if describedBy != "" {
		subControl.SetAttribute("aria-describedby", describedBy)
	}

	hasInnerDiv := c.innerDivAttr.Len() > 0
	if hasInnerDiv {
		buf.WriteString("<div ")
		buf.WriteString(c.innerDivAttr.String())
		buf.WriteString(">")
	}
	if err := c.this().ΩDrawInnerHtml(ctx, buf); err != nil {
		panic(err)
	}
	if hasInnerDiv {
		buf.WriteString("</div>")
	}
	if c.Instructions() != "" {
		buf.WriteString(html.RenderTag("small", c.InstructionAttributes(), html2.EscapeString(c.Instructions())))
	}
	if subControl.ValidationState() != page.ValidationNever {
		c.ErrorAttributes().SetClass(c.getValidationClass(subControl))
		buf.WriteString(html.RenderTag("div", c.ErrorAttributes(), errorMessage))
	}
	return html.RenderTag(c.Tag, attributes, buf.String())
}

func (c *FormGroup) getValidationClass(subcontrol page.ControlI) (class string) {
	switch subcontrol.ValidationState() {
	case page.ValidationWaiting: fallthrough // we need to correctly style
	case page.ValidationValid:
		if c.UseTooltips() {
			class = "valid-tooltip"
		} else {
			class = "valid-feedback"
		}

	case page.ValidationInvalid:
		if c.UseTooltips() {
			class = "invalid-tooltip"
		} else {
			class = "invalid-feedback"
		}
	}
	return
}

func (c *FormGroup) InnerDivAttributes() *html.Attributes {
	return c.innerDivAttr
}

type FormGroupCreator struct {
	ID string
	Label string

	// For specifies the id of the control that the label is for, and that is the control that we are wrapping.
	// You normally do not need this, as it will simply look at the first child control, but if for some reason
	// that control is wrapped, you should explicitly sepecify the For control id here.
	For string
	Instructions string
	LabelAttributes html.AttributeCreator
	ErrorAttributes html.AttributeCreator
	InstructionAttributes html.AttributeCreator
	Child page.Creator

	InnerDivAttributes html.AttributeCreator
	UseTooltips   bool // uses tooltips for the error class
	ControlOptions page.ControlOptions
}

func (f FormGroupCreator) Create(ctx context.Context, parent page.ControlI) FormGroupI {
	c := NewFormGroup(parent, f.ID)
	f.Init(ctx, c)
	return c
}

func (f FormGroupCreator) Init(ctx context.Context, c FormGroupI) {
	// Reuse parent creator
	ff := control.FormFieldCreator {
		ControlOptions: f.ControlOptions,
		Label: f.Label,
		For: f.For,
		Instructions: f.Instructions,
		LabelAttributes: f.LabelAttributes,
		ErrorAttributes: f.ErrorAttributes,
		InstructionAttributes: f.InstructionAttributes,
		Child: f.Child,
	}

	ff.Init(ctx, c)

	if f.InnerDivAttributes != nil {
		c.InnerDivAttributes().Merge(html.NewAttributesFromMap(f.InnerDivAttributes))
	}
	c.SetUseTooltips(f.UseTooltips)

}