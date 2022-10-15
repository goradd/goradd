package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	. "github.com/goradd/goradd/pkg/page/control/button"
)

type ImageCapturePanel struct {
	Panel
}

func NewImageCapturePanel(ctx context.Context, parent page.ControlI) {
	p := &ImageCapturePanel{}
	p.Self = p
	p.Init(ctx, parent, "imageCapturePanel")
}

func (p *ImageCapturePanel) Init(ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(parent, "imageCapturePanel")
	p.AddControls(ctx,
		FormFieldWrapperCreator{
			ID:    "ic1-ff",
			Label: "Circle ImageCapture",
			For:   "ic1",
			Child: ImageCaptureCreator{
				ID:        "ic1",
				MaskShape: ImageCaptureShapeCircle,
				SaveState: true,
				ControlOptions: page.ControlOptions{
					On: page.EventList{
						{ImageCaptureEvent(), action.Ajax(p.ID(), 0)}, // Just get the data.
					},
				},
			},
		},
		ButtonCreator{
			ID:       "ajaxButton",
			Text:     "Submit Ajax",
			OnSubmit: action.Ajax("checkboxPanel", ButtonSubmit),
		},
		ButtonCreator{
			ID:       "serverButton",
			Text:     "Submit Server",
			OnSubmit: action.Server("checkboxPanel", ButtonSubmit),
		},
	)
}

func init() {
	page.RegisterControl(&ImageCapturePanel{})
}
