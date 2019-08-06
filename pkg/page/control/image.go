package control

import (
	"context"
	"encoding/base64"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"strconv"
)

type ImageI interface {
	page.ControlI
	SetSrc(src string) ImageI
	SetAlt(alt string) ImageI
	SetWidth(width int) ImageI
	SetHeight(height int) ImageI
	SetMimeType(typ string) ImageI
	SetData(data []byte) ImageI
}

// Image is an img tag. You can display either a URL, or direct image information by setting the Src or the Data values.
type Image struct {
	page.Control
	data []byte // slice of data itself
	typ  string // the image MIME type (jpeg, gif, etc.) for data. Default is jpeg.
}

// NewImage creates a new image.
func NewImage(parent page.ControlI, id string) *Image {
	i := &Image{}
	i.Init(i, parent, id)
	return i
}

// Init is called by subclasses. Normally you will not call this directly.
func (i *Image) Init(self ImageI, parent page.ControlI, id string) {
	i.Control.Init(self, parent, id)
	i.Tag = "img"
	i.IsVoidTag = true
	i.typ = "jpeg"
}

func (i *Image) this() ImageI {
	return i.Self.(ImageI)
}

// Src returns the src attribute
func (i *Image) Src() string {
	return i.Attribute("src")
}

// SetSrc sets the src attribute.
func (i *Image) SetSrc(src string) ImageI {
	i.SetAttribute("src", src)
	return i.this()
}

// Data returns the data of the image if provided.
func (i *Image) Data() []byte {
	return i.data
}

// SetData sets the raw data of the image.
func (i *Image) SetData(data []byte) ImageI {
	i.data = data
	return i.this()
}

// Set the MIME type for the data, (jpeg, gif, png, etc.)
func (i *Image) SetMimeType(typ string) ImageI {
	i.typ = typ
	return i.this()
}

// Alt returns the text that will be used for the alt tag.
func (i *Image) Alt() string {
	return i.Attribute("alt")
}

// SetAlt will set the alt tag. The html standard requires the alt tag. Alt tags are used to display
// a descirption of an image when the browser cannot display an image, and is very important
// for assistive technologies.
func (i *Image) SetAlt(alt string) ImageI {
	i.SetAttribute("alt", alt)
	return i.this()
}

// Width returns the number that will be used as the width of the image
func (i *Image) Width() int {
	w := i.Attribute("width")
	if i, err := strconv.Atoi(w); err != nil {
		return 0
	} else {
		return i
	}
}

// SetWidth sets the width attribute.
func (i *Image) SetWidth(width int) ImageI {
	i.SetAttribute("width", strconv.Itoa(width))
	i.Refresh()
	return i.this()
}

// Height returns the number that will be used in the height attribute.
func (i *Image) Height() int {
	w := i.Attribute("height")
	if i, err := strconv.Atoi(w); err != nil {
		return 0
	} else {
		return i
	}
}

// SetHeight sets the height attribute.
func (i *Image) SetHeight(height int) ImageI {
	i.SetAttribute("width", strconv.Itoa(height))
	i.Refresh()
	return i.this()
}

// ΩDrawingAttributes is called by the framework.
func (i *Image) ΩDrawingAttributes() *html.Attributes {
	a := i.Control.ΩDrawingAttributes()
	if i.data != nil {
		// Turn the data into a source attribute
		d := base64.StdEncoding.EncodeToString(i.data)
		d = "data:image/" + i.typ + ";base64," + d
		a.Set("src", d)
	}
	return a
}


// ImageCreator is the initialization structure for declarative creation of buttons
type ImageCreator struct {
	// ID is the control id
	ID string
	// Src is the content of the source attribute, usually a url
	Src string
	// Alt is the text displayed for screen readers
	Alt string
	MimeType string
	Width int
	Height int
	Data []byte
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c ImageCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewImage(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of Images to initialize a control with the
// creator. You do not normally need to call this.
func (c ImageCreator) Init(ctx context.Context, ctrl ImageI) {
	ctrl.SetSrc(c.Src)
	ctrl.SetAlt(c.Alt)
	if c.MimeType != "" {
		ctrl.SetMimeType(c.MimeType)
	}
	if c.Width != 0 {
		ctrl.SetWidth(c.Width)
	}
	if c.Height != 0 {
		ctrl.SetHeight(c.Height)
	}
	if c.Data != nil {
		ctrl.SetData(c.Data)
	}
	ctrl.ApplyOptions(c.ControlOptions)
}

// GetImage is a convenience method to return the button with the given id from the page.
func GetImage(c page.ControlI, id string) *Image {
	return c.Page().GetControl(id).(*Image)
}
