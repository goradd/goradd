package panels

import (
	"context"
	. "github.com/goradd/goradd/pkg/bootstrap/control"
	"github.com/goradd/goradd/pkg/bootstrap/examples"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
)

// shared
const controlsFormPath = "/goradd/examples/bootstrap.g"

const (
	AjaxSubmit int = iota + 1
	ServerSubmit
	ButtonSubmit
)

type Forms1Panel struct {
	control.Panel
}


func NewForms1Panel(ctx context.Context, parent page.ControlI) {
	p := &Forms1Panel{}
	p.Panel.Init(p, parent, "textboxPanel")
	p.Panel.AddControls(ctx,
		FormGroupCreator{
			Label:"Name",
			Child: TextboxCreator{
				ID: "nameText",
			},
		},
		FormGroupCreator{
			Label:"Child Count",
			Instructions:"How many children do you have?",
			Child: IntegerTextboxCreator{
				ID: "childrenText",
				ColumnCount:2,
				ControlOptions: page.ControlOptions{
					IsRequired: true,
				},
			},
		},
		// Normally you would use a radio list for radio buttons.
		// This is just a demonstration of how you can do it without a radio list for special situations.
		RadioButtonCreator{
			ID: "singleRadio",
			Group:"m",
			Text:"Single",
			Checked:true,
		},
		RadioButtonCreator{
			ID: "marriedRadio",
			Group:"m",
			Text:"Married",
		},
		RadioButtonCreator{
			ID: "divorcedRadio",
			Group:"m",
			Text:"Divorced",
		},
		FormGroupCreator{
			Label: "Dog",
			Instructions:"Do you have a dog?",
			Child: CheckboxCreator {
				ID:   "dogCheck",
				Text: "I have a dog",
			},
		},
		ButtonCreator {
			ID: "ajaxButton",
			Text: "Submit Ajax",
			OnSubmit:action.Ajax(p.ID(), AjaxSubmit),
		},
		ButtonCreator {
			ID: "serverButton",
			Text: "Submit Server",
			OnSubmit:action.Ajax(p.ID(), ServerSubmit),
		},
	)
}

func init() {
	examples.RegisterPanel("forms1", "Forms 1", NewForms1Panel, 2)
	page.RegisterControl(Forms1Panel{})
	//browsertest.RegisterTestFunction("Bootstrap Standard Form Ajax Submit", testForms1AjaxSubmit)
	//browsertest.RegisterTestFunction("Bootstrap Standard Form Server Submit", testForms1ServerSubmit)
}

func testForms1AjaxSubmit(t *browsertest.TestForm) {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "textbox").String()
	f := t.LoadUrl(myUrl)

	testForms1Submit(t, f, f.Page().GetControl("ajaxButton"))

	t.Done("Complete")
}

func testForms1ServerSubmit(t *browsertest.TestForm) {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "textbox").String()
	f := t.LoadUrl(myUrl)

	testForms1Submit(t, f, f.Page().GetControl("serverButton"))

	t.Done("Complete")
}

// testForms1Submit does a variety of submits using the given button. We use this to double check the various
// results we might get after a submission, as well as nsure that the ajax and server submits produce
// the same results.
func testForms1Submit(t *browsertest.TestForm, f page.FormI, btn page.ControlI) {

	t.Click(btn)

}
