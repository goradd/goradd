package panels

import (
	. "github.com/goradd/goradd/pkg/bootstrap/control"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
)

// shared
const controlsFormPath = "/goradd/examples/controls.g"

const (
	AjaxSubmit int = iota + 1
	ServerSubmit
	ButtonSubmit
)

type Forms1Panel struct {
	control.Panel
	Name   *Textbox
	ChildrenCount   *IntegerTextbox
	MStatusSingle *RadioButton
	MStatusMarried *RadioButton
	MStatusDivorced *RadioButton
	Dog *Checkbox

	SubmitAjax      *Button
	SubmitServer    *Button
}

func NewForms1Panel(parent page.ControlI) *Forms1Panel {
	p := &Forms1Panel{}
	p.Panel.Init(p, parent, "textboxPanel")

	p.Name = NewTextbox(p, "nameText")
	p.Name.SetLabel("Name")

	p.ChildrenCount = NewIntegerTextbox(p, "childrenText")
	p.ChildrenCount.SetLabel("Child Count")
	p.ChildrenCount.SetInstructions("How many children do you have?")
	p.ChildrenCount.SetIsRequired(true)
	p.ChildrenCount.SetColumnCount(2)

	// Normally you would use a radio list for radio buttons.
	// This is just a demonstration of how you can do it without a radio list for special situations.
	p.MStatusSingle = NewRadioButton(p, "singleRadio")
	p.MStatusSingle.SetGroup("m")
	p.MStatusSingle.SetText("Single")
	p.MStatusSingle.SetChecked(true) // default a value

	p.MStatusMarried = NewRadioButton(p, "marriedRadio")
	p.MStatusMarried.SetGroup("m")
	p.MStatusMarried.SetText("Married")

	p.MStatusDivorced = NewRadioButton(p, "divorcedRadio")
	p.MStatusDivorced.SetGroup("m")
	p.MStatusDivorced.SetText("Divorced")

	p.Dog = NewCheckbox(p, "dogCheck")
	p.Dog.SetText("I have a dog")
	p.Dog.SetInstructions("Do you have a dog?")

	p.SubmitAjax = NewButton(p, "ajaxButton")
	p.SubmitAjax.SetText("Submit Ajax")
	p.SubmitAjax.OnSubmit(action.Ajax(p.ID(), AjaxSubmit))

	p.SubmitServer = NewButton(p, "serverButton")
	p.SubmitServer.SetText("Submit Server")
	p.SubmitServer.OnSubmit(action.Server(p.ID(), ServerSubmit))

	return p
}


func init() {
	//browsertest.RegisterTestFunction("Bootstrap Standard Form Ajax Submit", testForms1AjaxSubmit)
	//browsertest.RegisterTestFunction("Bootstrap Standard Form Server Submit", testForms1ServerSubmit)
}

func testForms1AjaxSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "textbox").String()
	f := t.LoadUrl(myUrl)

	testForms1Submit(t, f, f.Page().GetControl("ajaxButton"))

	t.Done("Complete")
}

func testForms1ServerSubmit(t *browsertest.TestForm)  {
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

