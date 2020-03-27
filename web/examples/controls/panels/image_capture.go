package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
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
			ID:"ic1-ff",
			Label:"Default ImageCapture",
			For:"ic1",
			Instructions:"Click to capture",
			Child:ImageCaptureCreator{
				ID:"ic1",
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
