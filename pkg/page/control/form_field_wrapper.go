package control

import (
	"context"
	"github.com/goradd/goradd/pkg/html5tag"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/pool"
	"html"
	"io"
	"reflect"
	"strings"
)

type LabelAttributer interface {
	LabelAttributes() html5tag.Attributes
}

type FormFieldWrapperI interface {
	page.ControlI
	SetFor(relatedId string) FormFieldWrapperI
	For() string
	Instructions() string
	SetInstructions(string) FormFieldWrapperI
	LabelAttributes() html5tag.Attributes
	ErrorAttributes() html5tag.Attributes
	InstructionAttributes() html5tag.Attributes
}

// FormFieldWrapper is a Goradd control that wraps other controls, and provides common companion
// functionality like a form label, validation state display, and help text.
type FormFieldWrapper struct {
	page.ControlBase

	// instructions is text associated with the control for extra explanation. You could also try adding a tooltip to the wrapper.
	instructions string
	// labelAttributes are the attributes that will be directly put on the Label tag. The label tag itself comes
	// from the "Text" item in the control.
	labelAttributes       html5tag.Attributes
	errorAttributes       html5tag.Attributes
	instructionAttributes html5tag.Attributes
	forID                 string
	// savedMessage is what we use to determine if the subcontrol changed validation state. This needs to be serialized.
	savedMessage string
	// subtag is the tag to used for instructions and error
	subtag string
}

func NewFormField(parent page.ControlI, id string) *FormFieldWrapper {
	p := &FormFieldWrapper{}
	p.Self = p
	p.Init(parent, id)
	return p
}

func (c *FormFieldWrapper) Init(parent page.ControlI, id string) {
	c.ControlBase.Init(parent, id)
	c.Tag = "div"
	c.subtag = "div"
	c.labelAttributes = html5tag.NewAttributes().
		SetID(c.ID() + "_lbl").
		SetClass("goradd-lbl")
	c.errorAttributes = html5tag.NewAttributes().
		SetID(c.ID() + "_err").
		SetClass("goradd-error")
	c.instructionAttributes = html5tag.NewAttributes().
		SetID(c.ID() + "_inst").
		SetClass("goradd-instructions")
}

func (c *FormFieldWrapper) this() FormFieldWrapperI {
	return c.Self.(FormFieldWrapperI)
}

// SetFor associates the form field with a sub control. The relatedId
// is the ID that the form field is associated with. Most browsers allow you to click on the
// label in order to give focus to the related control
func (c *FormFieldWrapper) SetFor(relatedId string) FormFieldWrapperI {
	if relatedId == "" {
		panic("A For id is required.")
	}
	c.forID = relatedId
	return c.this()
}

func (c *FormFieldWrapper) For() string {
	return c.forID
}

// SetInstructions sets the instructions that will be printed with the control. Instructions only get rendered
// by wrappers, so if there is no wrapper, or the wrapper does not render the instructions, this will not appear.
func (c *FormFieldWrapper) SetInstructions(i string) FormFieldWrapperI {
	if i != c.instructions {
		c.instructions = i
		c.Refresh()
	}
	return c.this()
}

// Instructions returns the instructions to be printed with the control
func (c *FormFieldWrapper) Instructions() string {
	return c.instructions
}

func (c *FormFieldWrapper) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := c.ControlBase.DrawingAttributes(ctx)
	a.SetData("grctl", "formField")
	return a
}

func (c *FormFieldWrapper) DrawTag(ctx context.Context, w io.Writer) {
	log.FrameworkDebug("Drawing FormFieldWrapper: " + c.ID())

	attributes := c.this().DrawingAttributes(ctx)
	var child page.ControlI
	var errorMessage string

	if c.Page().HasControl(c.forID) {
		child = c.Page().GetControl(c.forID)
		errorMessage = child.ValidationMessage()
		if errorMessage != "" {
			attributes.AddClass("error")
		}
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)

	text := c.Text()
	if text == "" && child != nil {
		text = child.Attribute("placeholder")
		if text != "" {
			c.labelAttributes.SetStyle("display", "none") // make a hidden label for screen readers
		}
	}

	if text != "" {
		if c.forID != "" {
			c.labelAttributes.Set("for", c.forID)
		}
		buf.WriteString(html5tag.RenderTag("label", c.labelAttributes, html.EscapeString(text)))
		if child != nil {
			child.SetAttribute("aria-labelledby", c.ID() + "_lbl")
		}
	}

	var describedBy string

	if errorMessage != "" {
		describedBy = c.ID() + "_err"
	}
	if c.instructions != "" {
		describedBy += " " + c.ID() + "_inst"
	}
	describedBy = strings.TrimSpace(describedBy)

	if describedBy != "" && child != nil {
		child.SetAttribute("aria-describedby", describedBy)
	}
	c.this().DrawInnerHtml(ctx, buf)
	if child != nil && child.ValidationState() != page.ValidationNever {
		page.WriteString(buf, html5tag.RenderTag(c.subtag, c.errorAttributes, html.EscapeString(errorMessage)))
	}
	if c.instructions != "" {
		page.WriteString(buf, html5tag.RenderTag(c.subtag, c.instructionAttributes, html.EscapeString(c.instructions)))
	}
	if _,err := io.WriteString(w, html5tag.RenderTag(c.Tag, attributes, buf.String())); err != nil {panic(err)}
}

func (c *FormFieldWrapper) LabelAttributes() html5tag.Attributes {
	return c.labelAttributes
}

func (c *FormFieldWrapper) SetLabelAttributes(a html5tag.Attributes) FormFieldWrapperI {
	c.labelAttributes = a
	return c.this()
}

func (c *FormFieldWrapper) ErrorAttributes() html5tag.Attributes {
	return c.errorAttributes
}

func (c *FormFieldWrapper) SetErrorAttributes(a html5tag.Attributes) FormFieldWrapperI {
	c.errorAttributes = a
	return c.this()
}

func (c *FormFieldWrapper) InstructionAttributes() html5tag.Attributes {
	return c.instructionAttributes
}

func (c *FormFieldWrapper) SetInstructionAttributes(a html5tag.Attributes) FormFieldWrapperI {
	c.instructionAttributes = a
	return c.this()
}

func (c *FormFieldWrapper) Validate(_ context.Context) bool {
	c.checkChildValidation()
	return true
}

func (c *FormFieldWrapper) ChildValidationChanged() {
	c.checkChildValidation()
	c.ControlBase.ChildValidationChanged()
}

func (c *FormFieldWrapper) checkChildValidation() {
	child := c.Page().GetControl(c.forID)
	m := child.ValidationMessage()
	if m != c.savedMessage {
		c.savedMessage = m // store the message to see if it changes between validations
		c.Refresh()
	}
}

func (c *FormFieldWrapper) SubTag() string {
	return c.subtag
}

func (c *FormFieldWrapper) Serialize(e page.Encoder) {
	c.ControlBase.Serialize(e)

	if err := e.Encode(c.instructions); err != nil {
		panic(err)
	}
	if err := e.Encode(c.labelAttributes); err != nil {
		panic(err)
	}
	if err := e.Encode(c.errorAttributes); err != nil {
		panic(err)
	}
	if err := e.Encode(c.instructionAttributes); err != nil {
		panic(err)
	}
	if err := e.Encode(c.forID); err != nil {
		panic(err)
	}
	if err := e.Encode(c.savedMessage); err != nil {
		panic(err)
	}
	if err := e.Encode(c.subtag); err != nil {
		panic(err)
	}
}

func (c *FormFieldWrapper) Deserialize(dec page.Decoder) {
	c.ControlBase.Deserialize(dec)

	if err := dec.Decode(&c.instructions); err != nil {
		panic(err)
	}

	if err := dec.Decode(&c.labelAttributes); err != nil {
		panic(err)
	}

	if err := dec.Decode(&c.errorAttributes); err != nil {
		panic(err)
	}

	if err := dec.Decode(&c.instructionAttributes); err != nil {
		panic(err)
	}

	if err := dec.Decode(&c.forID); err != nil {
		panic(err)
	}

	if err := dec.Decode(&c.savedMessage); err != nil {
		panic(err)
	}
	if err := dec.Decode(&c.subtag); err != nil {
		panic(err)
	}
}


// Use FormFieldWrapperCreator to create a FormFieldWrapper,
// which wraps a control with a div or span that also has a label, validation error
// text and optional instructions. Pass the creator of the control you
// are wrapping as the Child item.
type FormFieldWrapperCreator struct {
	// ID is the optional control id on the html form. If you do not specify this, it
	// will create on for you that is the ID of the child control + "-ff"
	ID string
	// Label is the text that will be in the html label tag associated with the Child control.
	Label string
	// Child is the creator of the child control you want to wrap
	Child page.Creator
	// Instructions is help text that will follow the control and that further describes its purpose or use.
	Instructions string
	// For specifies the id of the control that the label is for, and that is the control that we are wrapping.
	// You normally do not need this, as it will simply look at the first child control, but if for some reason
	// that control is wrapped, you should explicitly sepecify the For control id here.
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
}

// Create is called by the framework to create the control. You do not
// normally need to call it directly. Instead either pass this creator to
// AddControls for the parent control you want to add this to, or add this to
// the Children of the parent control's creator.
func (f FormFieldWrapperCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	id := CalcWrapperID(f.ID, f.Child, "ff")
	c := NewFormField(parent, id)
	f.Init(ctx, c)
	if f.IsInline { // subclasses might deal with this issue differently
		c.Tag = "span"
		c.subtag = "span"
	}

	return c
}

// Init is called by implementations of a FormFieldWrapper to initialize
// the creator. You do not normally need to call this.
func (f FormFieldWrapperCreator) Init(ctx context.Context, c FormFieldWrapperI) {
	c.ApplyOptions(ctx, f.ControlOptions)
	c.SetText(f.Label)
	c.SetInstructions(f.Instructions)
	if f.LabelAttributes != nil {
		c.LabelAttributes().Merge(f.LabelAttributes)
	}
	if f.ErrorAttributes != nil {
		c.ErrorAttributes().Merge(f.ErrorAttributes)
	}
	if f.InstructionAttributes != nil {
		c.InstructionAttributes().Merge(f.InstructionAttributes)
	}

	if f.Child == nil {
		panic("FormFieldWrapper controls require a child control")
	}
	c.AddControls(ctx, f.Child)
	if f.For != "" {
		c.SetFor(f.For)
	} else {
		childId := c.Children()[0].ID()
		c.SetFor(childId)
	}
}

// GetFormFieldWrapper is a convenience method to return the form field with the given id from the page.
func GetFormFieldWrapper(c page.ControlI, id string) *FormFieldWrapper {
	return c.Page().GetControl(id).(*FormFieldWrapper)
}

// GetCreatorID uses reflection to get the id of the given creator
func GetCreatorID(c page.Creator) string {
	v := reflect.ValueOf(c)
	f := v.FieldByName("ID")
	return f.String()
}

// CalcWrapperID returns the computed id of a control that wraps another control
// This would be the id of the child control followed by the postfix
func CalcWrapperID(wrapperId string, childCreator page.Creator, postfix string) string {
	id := wrapperId
	if id == ""  &&
		childCreator != nil {
		childId := GetCreatorID(childCreator)
		if childId != "" {
			id = childId + "-" + postfix
		}
	}
	return id
}

func init() {
	page.RegisterControl(&FormFieldWrapper{})
}