package control

import (
	"context"
	"github.com/goradd/goradd/pkg/config"
	"html"
	"io"
	"strings"

	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/pool"
	"github.com/goradd/html5tag"
)

type FormGroupI interface {
	control.FormFieldWrapperI
	SetUseTooltips(use bool) FormGroupI
	UseTooltips() bool
	InnerDivAttributes() html5tag.Attributes
}

// FormGroup is a Goradd control that wraps other controls, and provides common companion
// functionality like a form label, validation state display, and help text.
type FormGroup struct {
	control.FormFieldWrapper
	innerDivAttr html5tag.Attributes
	useTooltips  bool // uses tooltips for the error class
}

func NewFormGroup(parent page.ControlI, id string) *FormGroup {
	p := &FormGroup{}
	p.Self = p
	p.Init(parent, id)
	return p
}

func (c *FormGroup) Init(parent page.ControlI, id string) {
	c.FormFieldWrapper.Init(parent, id)
	c.innerDivAttr = html5tag.NewAttributes()
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
	if child.ValidationState() == page.ValidationWaiting {
		child.RemoveClass("is-valid")
		child.RemoveClass("is-invalid")
	} else if child.ValidationMessage() != "" {
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

func (c *FormGroup) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := c.FormFieldWrapper.DrawingAttributes(ctx)
	a.SetData("grctl", "formGroup")
	if c.useTooltips {
		// bootstrap requires that parent of a tool-tipped object has position relative
		c.SetStyle("position", "relative")
	}
	return a
}

func (c *FormGroup) DrawTag(ctx context.Context, w io.Writer) {
	log.FrameworkDebug("Drawing FormFieldWrapper: " + c.ID())

	attributes := c.this().DrawingAttributes(ctx)
	if c.For() == "" {
		panic("a FormGroup MUST have a sub control")
	}
	subControl := c.Page().GetControl(c.For())
	errorMessage := subControl.ValidationMessage()
	if errorMessage != "" {
		attributes.AddClass("error")
		errorMessage = html.EscapeString(errorMessage)
	} else {
		errorMessage = "&nbsp;"
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)

	text := c.Text()
	if text == "" {
		text = subControl.Attribute("placeholder")
		if text != "" {
			c.LabelAttributes().SetClass("visually-hidden") // make a hidden label for screen readers
		}
	}

	if text != "" {
		c.LabelAttributes().Set("for", c.For())
		buf.WriteString(html5tag.RenderTag("label", c.LabelAttributes(), html.EscapeString(text)))
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
		buf.WriteString(c.SubTag())
		buf.WriteString(" ")
		buf.WriteString(c.innerDivAttr.String())
		buf.WriteString(">")
	}
	c.this().DrawInnerHtml(ctx, buf)
	if hasInnerDiv {
		buf.WriteString("</")
		buf.WriteString(c.SubTag())
		buf.WriteString(">")
	}
	if c.Instructions() != "" {
		buf.WriteString(html5tag.RenderTag("div", c.InstructionAttributes(), html.EscapeString(c.Instructions())))
	}
	if subControl.ValidationState() != page.ValidationNever {
		c.ErrorAttributes().SetClass(c.getValidationClass(subControl))
		buf.WriteString(html5tag.RenderTag(c.SubTag(), c.ErrorAttributes(), errorMessage))
	}

	if _, err := io.WriteString(w, html5tag.RenderTag(c.Tag, attributes, buf.String())); err != nil {
		panic(err)
	}
}

func (c *FormGroup) getValidationClass(subcontrol page.ControlI) (class string) {
	switch subcontrol.ValidationState() {
	case page.ValidationWaiting:
		fallthrough
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

func (c *FormGroup) InnerDivAttributes() html5tag.Attributes {
	return c.innerDivAttr
}

func (c *FormGroup) Serialize(e page.Encoder) {
	c.FormFieldWrapper.Serialize(e)
	if err := e.Encode(c.innerDivAttr); err != nil {
		panic(err)
	}
	if err := e.Encode(c.useTooltips); err != nil {
		panic(err)
	}
}

func (c *FormGroup) Deserialize(dec page.Decoder) {
	c.FormFieldWrapper.Deserialize(dec)

	if err := dec.Decode(&c.innerDivAttr); err != nil {
		panic(err)
	}
	if err := dec.Decode(&c.useTooltips); err != nil {
		panic(err)
	}
}

// FormGroupCreator creates a FormGroup,
// which wraps a control with a div or span that also has a label, validation error
// text and optional instructions. Pass the creator of the control you
// are wrapping as the Child item.
type FormGroupCreator struct {
	// ID is the optional control id on the html form. If you do not specify this, it
	// will create on for you that is the ID of the child control + "-ff"
	ID string
	// Label is the text that will be in the html label tag associated with the Child control.
	Label string
	// Child is the creator of the child control you want to wrap
	Child page.Creator
	// Instructions contains help text that will follow the control and that further describes its purpose or use.
	Instructions string
	// For specifies the id of the control that the label is for, and that is the control that we are wrapping.
	// You normally do not need to specify this, as it will default to the first child control, but if for some reason
	// that control is wrapped, you should explicitly specify the For control id here.
	For string
	// LabelAttributes are additional attributes to add to the label tag.
	LabelAttributes html5tag.Attributes
	// ErrorAttributes are additional attributes to add to the tag that displays the error.
	ErrorAttributes html5tag.Attributes
	// InstructionAttributes are additional attributes to add to the tag that displays the instructions.
	InstructionAttributes html5tag.Attributes
	// Set IsInline to true to use a "span" instead of a "div" in the wrapping tag.
	IsInline bool
	// ControlOptions are additional options for the wrapper tag
	ControlOptions page.ControlOptions
	// InnerDivAttributes are the attributes for the additional div wrapper of the control
	// To achieve certain effects, Bootstrap needs this addition div. To display the div, you
	// must specify its attributes here. Otherwise, no inner div will be displayed.
	InnerDivAttributes html5tag.Attributes
	// UseTooltips will cause validation errors to be displayed with tooltips, a specific
	// feature of Bootstrap
	UseTooltips bool
}

// Create is called by the framework to create the control. You do not
// normally need to call it.
func (f FormGroupCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	id := control.MakeCreatorWrapperID(f.ID, f.Child, config.DefaultFormFieldWrapperIdSuffix)
	c := NewFormGroup(parent, id)
	f.Init(ctx, c)
	return c
}

// Init is called by implementations of a FormFieldWrapper to initialize
// the creator.
func (f FormGroupCreator) Init(ctx context.Context, c FormGroupI) {
	// Reuse parent creator
	ff := control.FormFieldWrapperCreator{
		ControlOptions:        f.ControlOptions,
		Label:                 f.Label,
		For:                   f.For,
		Instructions:          f.Instructions,
		LabelAttributes:       f.LabelAttributes,
		ErrorAttributes:       f.ErrorAttributes,
		InstructionAttributes: f.InstructionAttributes,
		Child:                 f.Child,
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
	page.RegisterControl(&FormGroup{})
}
