package control

import (
	"github.com/spekary/goradd/pkg/page"
	"strconv"
	"github.com/spekary/goradd/pkg/html"
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

func NewImage(parent page.ControlI, id string) *Image {
	i := &Image{}
	i.Init(i, parent, id)
	return i
}

// Initializes a textbox. Normally you will not call this directly. However, sub controls should call this after
// creation to get the enclosed control initialized. Self is the newly created class. Like so:
// t := &MyTextBox{}
// t.Textbox.Init(t, parent, id)
// A parent control is isRequired. Leave id blank to have the system assign an id to the control.
func (i *Image) Init(self ImageI, parent page.ControlI, id string) {
	i.Control.Init(self, parent, id)
	i.Tag = "img"
	i.IsVoidTag = true
	i.typ = "jpeg"
}

func (i *Image) this() ImageI {
	return i.Self.(ImageI)
}

func (i *Image) Src() string {
	return i.Attribute("src")
}

func (i *Image) SetSrc(src string) {
	i.SetAttribute("src", src)
}

func (i *Image) Data() []byte {
	return i.data
}

func (i *Image) SetData(data []byte) {
	i.data = data
}

// Set the MIME type for the data, (jpeg, gif, png, etc.)
func (i *Image) SetMimeType(typ string) {
	i.typ = typ
}

func (i *Image) Alt() string {
	return i.Attribute("alt")
}

func (i *Image) SetAlt(alt string) {
	i.SetAttribute("alt", alt)
}

func (i *Image) Width() int {
	w := i.Attribute("width")
	if i,err := strconv.Atoi(w); err != nil {
		return 0
	} else {
		return i
	}
}

func (i *Image) SetWidth(width int) {
	i.SetAttribute("width", strconv.Itoa(width))
	i.Refresh()
}

func (i *Image) Height() int {
	w := i.Attribute("height")
	if i,err := strconv.Atoi(w); err != nil {
		return 0
	} else {
		return i
	}
}

func (i *Image) SetHeight(height int) {
	i.SetAttribute("width", strconv.Itoa(height))
	i.Refresh()
}

func (i *Image) DrawingAttributes() *html.Attributes {
	a := i.Control.DrawingAttributes()
	if i.data != nil {
		// Turn the data into a source attribute
		d := base64.StdEncoding.EncodeToString(i.data)
		d = "data:image/" + i.typ + ";base64," + d
		a.Set("src", d)
	}
	return a
}


