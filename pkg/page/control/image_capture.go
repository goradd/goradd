package control

import (
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/page"
	"context"
	"encoding/base64"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/log"
	"strings"
)

type ImageCaptureShape string

const (
	ImageCaptureShapeRect ImageCaptureShape = "rect"
	ImageCaptureShapeCircle ImageCaptureShape = "circle"
)

const imageCaptureScriptCommand = "imageCapture"

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
	Canvas        *Canvas
	CaptureButton *Button
	SwitchCameraButton *Button

	ErrText *Panel
	data []byte
	shape ImageCaptureShape
	typ string
	zoom int
	quality float32
}

// NewImageCapture creates a new image capture panel.
func NewImageCapture(parent page.ControlI, id string) *ImageCapture {
	i := &ImageCapture{}
	i.Init(i, parent, id)
	return i
}

// Init is called by subclasses.
func (i *ImageCapture) Init(self ImageCaptureI, parent page.ControlI, id string) {
	i.Control.Init(self, parent, id)
	i.Tag = "div"
	i.ParentForm().AddJavaScriptFile(config.GoraddAssets() + "/js/image-capture.js", false, nil)
	i.typ = "jpeg"
	i.quality = 0.92

	i.Canvas = NewCanvas(i, i.ID() + "_canvas")
	i.CaptureButton = NewButton(i, i.ID() + "_capture")
	i.CaptureButton.SetText(i.ΩT("New Image"))

	i.SwitchCameraButton = NewButton(i, i.ID() + "_switch")
	i.SwitchCameraButton.SetText(i.ΩT("Switch Camera"))
	i.SwitchCameraButton.SetDisplay("none")

	i.ErrText = NewPanel(i, i.ID() + "_err")
	i.ErrText.Tag = "p"
	i.ErrText.SetDisplay("none")
	i.ErrText.SetText(i.ΩT("This browser or device does not support image capture"))
}

func (i *ImageCapture) this() ImageCaptureI {
	return i.Self.(ImageCaptureI)
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

// ΩPutCustomScript is called by the framework.
func (i *ImageCapture) ΩPutCustomScript(ctx context.Context, response *page.Response) {
	options := map[string]interface{}{}
	d := base64.StdEncoding.EncodeToString(i.data)
	d = "data:image/" + i.typ + ";base64," + d
	options["data"] = d
	options["selectButtonName"] = i.ΩT("Capture")
	if i.zoom > 0 {
		options["zoom"] = i.zoom
	}
	if i.shape != "" {
		options["shape"] = string(i.shape)
	}
	options["mimeType"] = i.typ
	options["quality"] = i.quality

	response.ExecuteControlCommand(i.ID(), imageCaptureScriptCommand, page.PriorityHigh, options)
}

// TurnOff will turn off the camera and the image displayed in the control
func (i *ImageCapture) TurnOff() {
	i.ParentForm().Response().ExecuteControlCommand(i.ID(), imageCaptureScriptCommand, page.PriorityHigh, "turnOff")
}


// SetPixelSize sets the pixel size of the image that will be returned. Control the visible size of the canvas through
// setting css sizes.
func (i *ImageCapture) SetPixelSize(width int, height int) {
	i.Canvas.SetAttribute("width", width)
	i.Canvas.SetAttribute("height", height)
}

// SetMaskShape sets the masking shape for the image
func (i *ImageCapture) SetMaskShape (shape ImageCaptureShape) {
	i.shape = shape
}

// ΩDrawingAttributes is called by the framework.
func (i *ImageCapture) ΩDrawingAttributes() *html.Attributes {
	a := i.Control.ΩDrawingAttributes()
	if i.data != nil {
		// Turn the data into a source attribute
		d := base64.StdEncoding.EncodeToString(i.data)
		d = "data:image/" + i.typ + ";base64," + d
		a.Set("src", d)
	}
	return a
}

// ΩUpdateFormValues is called by the framework.
func (i *ImageCapture) ΩUpdateFormValues(ctx *page.Context) {
	if data := ctx.CustomControlValue(i.ID(), "data"); data != nil {
		s := data.(string)
		index := strings.Index(s, ",")
		if newdata,err := base64.StdEncoding.DecodeString(s[index+1:]); err == nil {
			i.data = newdata
		} else {
			log.Debug(err.Error())
		}
	}
}