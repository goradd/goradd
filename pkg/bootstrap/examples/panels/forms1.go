package panels

import (
	"context"
	. "github.com/goradd/goradd/pkg/bootstrap/control"
	"github.com/goradd/goradd/pkg/bootstrap/examples"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/list"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
)

// shared
const controlsFormPath = "/goradd/examples/bootstrap.g"

const (
	AjaxSubmit int = iota + 1
	ServerSubmit
	ButtonSubmit
	RadioChange
)

type Forms1Panel struct {
	control.Panel
}

func NewForms1Panel(ctx context.Context, parent page.ControlI) {
	p := &Forms1Panel{}
	p.Self = p
	p.Init(ctx, parent, "textboxPanel")

}

func (p *Forms1Panel) Init(ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(parent, id)
	p.Panel.AddControls(ctx,
		FormGroupCreator{
			ID:    "nameText-ff",
			Label: "Name",
			Child: TextboxCreator{
				ID: "nameText",
			},
		},
		FormGroupCreator{
			ID:           "childrenText-ff",
			Label:        "Child Count",
			Instructions: "How many children do you have?",
			Child: IntegerTextboxCreator{
				ID:          "childrenText",
				ColumnCount: 2,
				ControlOptions: page.ControlOptions{
					IsRequired: true,
				},
			},
		},
		RadioListGroupCreator{
			ID: "status",
			Items: []list.ListValue{
				{"Single", "Single"},
				{"Married", "Married"},
				{"Divorced", "Divorced"},
			},
			Value:    "Single",
			OnChange: action.Ajax(p.ID(), RadioChange),
		},
		control.SpanCreator{
			ID: "radioResult",
		},
		FormGroupCreator{
			ID:           "dogCheck-ff",
			Label:        "Dog",
			Instructions: "Do you have a dog?",
			Child: CheckboxCreator{
				ID:   "dogCheck",
				Text: "I have a dog",
			},
		},
		ButtonCreator{
			ID:       "ajaxButton",
			Text:     "Submit Ajax",
			OnSubmit: action.Ajax(p.ID(), AjaxSubmit),
		},
		ButtonCreator{
			ID:       "serverButton",
			Text:     "Submit Server",
			OnSubmit: action.Server(p.ID(), ServerSubmit),
		},
	)
}

func (p *Forms1Panel) DoAction(ctx context.Context, a action.Params) {
	switch a.ID {
	case RadioChange:
		c := control.GetSpan(p, "radioResult")
		r := GetRadioListGroup(p, "status")
		s := r.StringValue()
		c.SetText(s)
	}
}

func init() {
	examples.RegisterPanel("forms1", "Forms 1", NewForms1Panel, 2)
	page.RegisterControl(&Forms1Panel{})
	//browsertest.RegisterTestFunction("Bootstrap Standard Form Ajax Submit", testForms1AjaxSubmit)
	//browsertest.RegisterTestFunction("Bootstrap Standard Form Server Submit", testForms1ServerSubmit)
}

func testForms1AjaxSubmit(t *browsertest.TestForm) {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "textbox").String()
	t.LoadUrl(myUrl)

	testForms1Submit(t, "ajaxButton")

	t.Done("Complete")
}

func testForms1ServerSubmit(t *browsertest.TestForm) {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "textbox").String()
	t.LoadUrl(myUrl)

	testForms1Submit(t, "serverButton")

	t.Done("Complete")
}

// testForms1Submit does a variety of submits using the given button. We use this to double check the various
// results we might get after a submission, as well as nsure that the ajax and server submits produce
// the same results.
func testForms1Submit(t *browsertest.TestForm, id string) {
	t.Click(id)
}
