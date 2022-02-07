package control

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/page"
	"path"
	"strings"
)

type ImageCaptureShape string

const (
	ImageCaptureShapeRect   ImageCaptureShape = "rect"
	ImageCaptureShapeCircle ImageCaptureShape = "circle"
)

const imageCaptureScriptCommand = "imageCapture"

// CaptureEvent triggers when the capture button has been pressed, and the image has been captured.
func CaptureEvent() *page.Event {
	return &page.Event{JsEvent: "capture"}
}

type ImageCaptureI interface {
	page.ControlI
}

// ImageCapture is a panel that has both an image and button to help you capture images from the user's camera.
// It is a kind of composite control that exports the image so that you can further manipulate it after
// creation. It also has javascript to manage the actual image capture process. It does not currently allow
// the user to upload an image in place of capturing an image from the camera.
// It only captures images from devices and browsers that support image capture.
type ImageCapture struct {
	Panel

	ErrTextID string
	data    []byte
	shape   ImageCaptureShape
	typ     string
	zoom    int
	quality float32
}

// NewImageCapture creates a new image capture panel.
func NewImageCapture(parent page.ControlI, id string) *ImageCapture {
	i := &ImageCapture{}
	i.Self = i
	i.Init(parent, id)
	return i
}

// Init is called by subclasses.
func (i *ImageCapture) Init(parent page.ControlI, id string) {
	i.Panel.Init(parent, id)
	i.ParentForm().AddJavaScriptFile(path.Join(config.AssetPrefix, "goradd", "/js/image-capture.js"), false, nil)
	i.typ = "jpeg"
	i.quality = 0.92

	NewCanvas(i, i.canvasID())
	
	NewButton(i, i.captureID()).
		SetText(i.GT("New Image"))

	NewButton(i, i.switchID()).
		SetDisplay("none").
		SetText(i.GT("Switch Camera"))

	i.ErrTextID = i.ID()+"-err"
	et := NewPanel(i, i.ErrTextID)
	et.Tag = "p"
	et.SetDisplay("none")
	et.SetText(i.GT("This browser or device does not support image capture"))
}

func (i *ImageCapture) this() ImageCaptureI {
	return i.Self.(ImageCaptureI)
}

func (i *ImageCapture) canvasID() string {
	return i.ID() + "-canvas"
}

func (i *ImageCapture) captureID() string {
	return i.ID() + "-capture"
}

func (i *ImageCapture) switchID() string {
	return i.ID() + "-switch"
}


func (i *ImageCapture) Data() []byte {
	return i.data // clone?
}

// SetData sets the binary picture data. The data must be in the mime type format.
func (i *ImageCapture) SetData(data []byte) {
	i.data = data
	i.AddRenderScript("option", "data", data) // Set just the data through javascript if possible
}

func (i *ImageCapture) SetMimeType(typ string) {
	i.typ = typ
}

// SetQuality specifies a number between 0 and 1 used as the quality value for capturing jpegs or webp images.
func (i *ImageCapture) SetQuality(quality float32) {
	i.quality = quality
}

// SetZoom zooms the camera by the given percent, i.e. 50 is 50% closer and 100 would be a 2x zoom.
func (i *ImageCapture) SetZoom(zoom int) {
	i.zoom = zoom
}

// TurnOff will turn off the camera and the image displayed in the control
func (i *ImageCapture) TurnOff() {
	i.ExecuteWidgetFunction("turnOff")
}

// SetPixelSize sets the pixel size of the image that will be returned. ControlBase the visible size of the canvas through
// setting css sizes.
func (i *ImageCapture) SetPixelSize(width int, height int) {
	canvas := GetCanvas(i, i.canvasID())
	canvas.SetAttribute("width", width)
	canvas.SetAttribute("height", height)
}

// SetMaskShape sets the masking shape for the image
func (i *ImageCapture) SetMaskShape(shape ImageCaptureShape) {
	i.shape = shape
}

/*
// PutCustomScript is called by the framework.

// The code below is being preserved to show an example of how you could connect an html object to a different
// javascript library using javascript functions. This control was initially attached using the JQuery UI Widget library, which is no longer
// in active development.
func (i *ImageCapture) PutCustomScript(ctx context.Context, response *page.Response) {
	options := map[string]interface{}{}
	d := base64.StdEncoding.EncodeToString(i.data)
	d = "data:image/" + i.typ + ";base64," + d
	options["data"] = d
	options["selectImageCaptureName"] = i.GT("Capture")
	if i.zoom > 0 {
		options["zoom"] = i.zoom
	}
	if i.shape != "" {
		options["shape"] = string(i.shape)
	}
	options["mimeType"] = i.typ
	options["quality"] = i.quality

	response.ExecuteJqueryCommand(i.ID(), imageCaptureScriptCommand, page.PriorityHigh, options)
}
*/

// DrawingAttributes is called by the framework.
func (i *ImageCapture) DrawingAttributes(ctx context.Context) html.Attributes {
	a := i.Panel.DrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "imagecapture")
	a.SetDataAttribute("grWidget", "goradd.ImageCapture")

	d := base64.StdEncoding.EncodeToString(i.data)
	d = "data:image/" + i.typ + ";base64," + d
	a.SetDataAttribute("grOptData", d)
	a.SetDataAttribute("grOptSelectButtonName", i.GT("Capture"))
	if i.zoom > 0 {
		a.SetDataAttribute("grOptZoom", fmt.Sprint(i.zoom))
	}
	if i.shape != "" {
		a.SetDataAttribute("grOptShape", fmt.Sprint(i.shape))
	}
	a.SetDataAttribute("grOptMimeType", i.typ)
	a.SetDataAttribute("grOptQuality", fmt.Sprint(i.quality))
/*
	if i.data != nil {
		// Turn the data into a source attribute
		d := base64.StdEncoding.EncodeToString(i.data)
		d = "data:image/" + i.typ + ";base64," + d
		a.Set("src", d)
	}*/
	return a
}

// UpdateFormValues is called by the framework.
func (i *ImageCapture) UpdateFormValues(ctx context.Context) {
	if data := page.GetContext(ctx).CustomControlValue(i.ID(), "data"); data != nil {
		s := data.(string)
		index := strings.Index(s, ",")
		if newdata, err := base64.StdEncoding.DecodeString(s[index+1:]); err == nil {
			i.data = newdata
		} else {
			log.Debug(err.Error())
		}
	}
}
// MarshalState is an internal function to save the state of the control
func (i *ImageCapture) MarshalState(m maps.Setter) {
	m.Set("data", i.Data())
}

// UnmarshalState is an internal function to restore the state of the control
func (i *ImageCapture) UnmarshalState(m maps.Loader) {
	if v, ok := m.Load("data"); ok {
		if s, ok := v.([]byte); ok {
			i.data = s
		}
	}
}

func (i *ImageCapture) Serialize(e page.Encoder) (err error) {
	if err = i.ControlBase.Serialize(e); err != nil {
		return
	}

	if err = e.Encode(i.ErrTextID); err != nil {
		return
	}
	if err = e.Encode(i.data); err != nil {
		return
	}
	if err = e.Encode(i.shape); err != nil {
		return
	}
	if err = e.Encode(i.typ); err != nil {
		return
	}
	if err = e.Encode(i.zoom); err != nil {
		return
	}
	if err = e.Encode(i.quality); err != nil {
		return
	}

	return
}

func (i *ImageCapture) Deserialize(dec page.Decoder) (err error) {
	if err = i.ControlBase.Deserialize(dec); err != nil {
		return
	}

	if err = dec.Decode(&i.ErrTextID); err != nil {
		return
	}

	if err = dec.Decode(&i.data); err != nil {
		return
	}
	if err = dec.Decode(&i.shape); err != nil {
		return
	}

	if err = dec.Decode(&i.typ); err != nil {
		return
	}
	if err = dec.Decode(&i.zoom); err != nil {
		return
	}

	if err = dec.Decode(&i.quality); err != nil {
		return
	}

	return
}

// ImageCaptureCreator is the initialization structure for declarative creation of buttons
type ImageCaptureCreator struct {
	// ID is the control id
	ID string
	MaskShape   	ImageCaptureShape
	MimeType    string
	Zoom    int
	Quality float32
	SaveState bool
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c ImageCaptureCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewImageCapture(parent, c.ID)
	if c.MaskShape != "" {
		ctrl.SetMaskShape(c.MaskShape)
	}
	if c.MimeType != "" {
		ctrl.SetMimeType(c.MimeType)
	}
	if c.Zoom != 0 {
		ctrl.SetZoom(c.Zoom)
	}
	if c.Quality != 0 {
		ctrl.SetQuality(c.Quality)
	}
	if c.SaveState {
		ctrl.SaveState(ctx, true)
	}
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	return ctrl
}

// GetImageCapture is a convenience method to return the button with the given id from the page.
func GetImageCapture(c page.ControlI, id string) *ImageCapture {
	return c.Page().GetControl(id).(*ImageCapture)
}

func init() {
	page.RegisterControl(&ImageCapture{})
}