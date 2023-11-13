package panels

import (
	"context"
	. "github.com/goradd/goradd/pkg/bootstrap/control"
	"github.com/goradd/goradd/pkg/bootstrap/examples"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
)

// shared
const ()

type Forms2Panel struct {
	control.Panel
}

func NewForms2Panel(ctx context.Context, parent page.ControlI) {
	p := new(Forms2Panel)
	p.Init(p, ctx, parent, "textboxPanel")

}

func (p *Forms2Panel) Init(self any, ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, id)
	p.Panel.AddControls(ctx,
		FormGroupCreator{
			Label: "Name",
			Child: TextboxCreator{
				ID: "nameText",
				ControlOptions: page.ControlOptions{
					IsRequired: true,
				},
			},
		},
		FormGroupCreator{
			Label: "Street",
			Child: TextboxCreator{
				ID: "streetText",
			},
		},
		FormGroupCreator{
			Label: "City",
			Child: TextboxCreator{
				ID: "cityText",
			},
		},
		FormGroupCreator{
			Label: "State",
			Child: TextboxCreator{
				ColumnCount: 2,
				MaxLength:   2,
				ID:          "stateText",
			},
		},
		FormGroupCreator{
			Label: "Zip",
			Child: TextboxCreator{
				ColumnCount: 10,
				MaxLength:   10,
				ID:          "zipText",
			},
		},

		ButtonCreator{
			ID:       "ajaxButton",
			Text:     "Submit Ajax",
			OnSubmit: action.Do().ID(AjaxSubmit),
		},
		ButtonCreator{
			ID:       "serverButton",
			Text:     "Submit Server",
			OnSubmit: action.Do().ID(PostSubmit).Post(),
		},
	)
}

func init() {
	examples.RegisterPanel("forms2", "Forms 2", NewForms2Panel, 2)
	page.RegisterControl(new(Forms2Panel))
	//browsertest.RegisterTestFunction("Bootstrap Standard Form Do Submit", testForms1AjaxSubmit)
	//browsertest.RegisterTestFunction("Bootstrap Standard Form Server Submit", testForms1ServerSubmit)
}
