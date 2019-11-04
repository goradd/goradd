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
	control.FormFieldWrapperI
	SetUseTooltips(use bool) FormGroupI
	UseTooltips() bool
	InnerDivAttributes() html.Attributes
}

// FormGroup is a Goradd control that wraps other controls, and provides common companion
// functionality like a form label, validation state display, and help text.
type FormGroup struct {
	control.FormFieldWrapper
	innerDivAttr html.Attributes
	useTooltips   bool // uses tooltips for the error class
}

func NewFormGroup(parent page.ControlI, id string) *FormGroup {
	p := &FormGroup{}
	p.Init(p, parent, id)
	return p
}

func (c *FormGroup) Init(self control.FormFieldWrapperI, parent page.ControlI, id string) {
	c.FormFieldWrapper.Init(self, parent, id)
	c.innerDivAttr = html.NewAttributes()
	c.AddClass("form-group") // to get a wrapper without this, just remove it after initialization
	c.InstructionAttributes().AddClass("form-text")
}

func (c *FormGroup) this() FormGroupI {
	return c.Self.(FormGroupI)
}

func (c *FormGroup) Validate(ctx context.Context) bool {
	c.setChildValidation()
	c.FormFieldWrapper.Validate(ctx)

	return true
}

func (c *FormGroup) ChildValidationChanged() {
	c.setChildValidation()
	c.FormFieldWrapper.ChildValidationChanged()
}

func (c *FormGroup) setChildValidation() {
	child := c.Page().GetControl(c.For())
	if child.ValidationMessage() != "" {
		child.RemoveClass("is-valid")
		child.AddClass("is-invalid")
	} else {
		child.AddClass("is-valid")
		child.RemoveClass("is-invalid")
	}
}


// SetUseTooltips sets whether to use tooltips to display validation messages.
func (c *FormGroup) SetUseTooltips(use bool) FormGroupI {
	c.useTooltips = use
	return c
}

func (c *FormGroup) UseTooltips() bool {
	return c.useTooltips
}

func (c *FormGroup) DrawingAttributes(ctx context.Context) html.Attributes {
	a := c.FormFieldWrapper.DrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "formGroup")
	if c.useTooltips {
		// bootstrap requires that parent of a tool-tipped object has position relative
		c.SetStyle("position", "relative")
	}
	return a
}

func (c *FormGroup) DrawTag(ctx context.Context) string {
	log.FrameworkDebug("Drawing FormFieldWrapper: " + c.ID())

	attributes := c.this().DrawingAttributes(ctx)
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
		buf.WriteString("<")
		buf.WriteString(c.Subtag)
		buf.WriteString(" ")
		buf.WriteString(c.innerDivAttr.String())
		buf.WriteString(">")
	}
	if err := c.this().DrawInnerHtml(ctx, buf); err != nil {
		panic(err)
	}
	if hasInnerDiv {
		buf.WriteString("</")
		buf.WriteString(c.Subtag)
		buf.WriteString(">")
	}
	if c.Instructions() != "" {
		buf.WriteString(html.RenderTag("small", c.InstructionAttributes(), html2.EscapeString(c.Instructions())))
	}
	if subControl.ValidationState() != page.ValidationNever {
		c.ErrorAttributes().SetClass(c.getValidationClass(subControl))
		buf.WriteString(html.RenderTag(c.Subtag, c.ErrorAttributes(), errorMessage))
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

func (c *FormGroup) InnerDivAttributes() html.Attributes {
	return c.innerDivAttr
}

// Use FormGroupCreator to create a FormGroup,
// which wraps a control with a div or span that also has a label, validation error
// text and optional instructions. Pass the creator of the control you
// are wrapping as the Child item.
type FormGroupCreator struct {
	// ID is the optional control id on the html form. If you do not specify this, it
	// will create on for you that is the ID of the child control + "-fg"
	ID string
	// Label is the text that will be in the html label tag associated with the Child control.
	Label string
	// Child is the creator of the child control you want to wrap
	Child page.Creator
	// Instructions is help text that will follow the control and that further describes its purpose or use.
	Instructions string
	// For specifies the id of the control that the label is for, and that is the control that we are wrapping.
	// You normally do not need to specify this, as it will default to the first child control, but if for some reason
	// that control is wrapped, you should explicitly specify the For control id here.
	For string
	// LabelAttributes are additional attributes to add to the label tag.
	LabelAttributes html.Attributes
	// ErrorAttributes are additional attributes to add to the tag that displays the error.
	ErrorAttributes html.Attributes
	// InstructionAttributes are additional attributes to add to the tag that displays the instructions.
	InstructionAttributes html.Attributes
	// Set IsInline to true to use a "span" instead of a "div" in the wrapping tag.
	IsInline bool
	// ControlOptions are additional options for the wrapper tag
	ControlOptions page.ControlOptions
	// InnerDivAttributes are the attributes for the additional div wrapper of the control
	// To achieve certain effects, Bootstrap needs this addition div. To display the div, you
	// must specify its attributes here. Otherwise no inner div will be displayed.
	InnerDivAttributes html.Attributes
	// UseTooltips will cause validation errors to be displayed with tooltips, a specific
	// feature of Bootstrap
	UseTooltips   bool
}

// Create is called by the framework to create the control. You do not
// normally need to call it.
func (f FormGroupCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	id := control.CalcWrapperID(f.ID, f.Child, "fg")
	c := NewFormGroup(parent, id)
	f.Init(ctx, c)
	return c
}

// Init is called by implementations of a FormFieldWrapper to initialize
// the creator. You do not normally need to call this.
func (f FormGroupCreator) Init(ctx context.Context, c FormGroupI) {
	// Reuse parent creator
	ff := control.FormFieldWrapperCreator{
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
		c.InnerDivAttributes().Merge(f.InnerDivAttributes)
	}
	c.SetUseTooltips(f.UseTooltips)
}

// GetFormGroup is a convenience method to return the form group with the given id from the page.
func GetFormGroup(c page.ControlI, id string) *FormGroup {
	return c.Page().GetControl(id).(*FormGroup)
}

func init() {
	page.RegisterControl(FormGroup{})
}