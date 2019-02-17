package control

import (
	"github.com/goradd/goradd/pkg/page"
	"strconv"
	"github.com/goradd/goradd/pkg/html"
	"encoding/base64"
)

type ImageI interface {
	page.ControlI
}

// Image is an img tag. You can display either a URL, or direct image information by setting the Src or the Data values.
type Image struct {
	page.Control
	data []byte			// slice of data itself
	typ string			// the image MIME type (jpeg, gif, etc.) for data. Default is jpeg.
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
func (i *Image) SetSrc(src string) {
	i.SetAttribute("src", src)
}

// Data returns the data of the image if provided.
func (i *Image) Data() []byte {
	return i.data
}

// SetData sets the raw data of the image.
func (i *Image) SetData(data []byte) {
	i.data = data
}

// Set the MIME type for the data, (jpeg, gif, png, etc.)
func (i *Image) SetMimeType(typ string) {
	i.typ = typ
}

// Alt returns the text that will be used for the alt tag.
func (i *Image) Alt() string {
	return i.Attribute("alt")
}

// SetAlt will set the alt tag. The html standard requires the alt tag. Alt tags are used to display
// a descirption of an image when the browser cannot display an image, and is very important
// for assistive technologies.
func (i *Image) SetAlt(alt string) {
	i.SetAttribute("alt", alt)
}

// Width returns the number that will be used as the width of the image
func (i *Image) Width() int {
	w := i.Attribute("width")
	if i,err := strconv.Atoi(w); err != nil {
		return 0
	} else {
		return i
	}
}

// SetWidth sets the width attribute.
func (i *Image) SetWidth(width int) {
	i.SetAttribute("width", strconv.Itoa(width))
	i.Refresh()
}

// Height returns the number that will be used in the height attribute.
func (i *Image) Height() int {
	w := i.Attribute("height")
	if i,err := strconv.Atoi(w); err != nil {
		return 0
	} else {
		return i
	}
}

// SetHeight sets the height attribute.
func (i *Image) SetHeight(height int) {
	i.SetAttribute("width", strconv.Itoa(height))
	i.Refresh()
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


