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
	p := new(ImageCapturePanel)
	p.Init(p, ctx, parent, "imageCapturePanel")
}

func (p *ImageCapturePanel) Init(self any, ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, id)
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
						ImageCaptureEvent(), // Just get the data.
					},
				},
			},
		},
		ButtonCreator{
			ID:       "ajaxButton",
			Text:     "Submit Ajax",
			OnSubmit: action.Do().ControlID("checkboxPanel").ID(ButtonSubmit),
		},
		ButtonCreator{
			ID:       "serverButton",
			Text:     "Submit Post",
			OnSubmit: action.Do().ControlID("checkboxPanel").ID(ButtonSubmit).Post(),
		},
	)
}

func init() {
	page.RegisterControl(&ImageCapturePanel{})
}
