package panels

import (
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
)

// shared
const controlsFormPath = "/goradd/examples/controls.g"

const (
	AjaxSubmit int = iota + 1
	ServerSubmit
)

type TextboxPanel struct {
	Panel
	PlainText   *Textbox
	IntegerText *IntegerTextbox
	FloatText   *FloatTextbox
	SubmitAjax      *Button
	SubmitServer    *Button
}

func NewTextboxPanel(parent page.ControlI, id string) *TextboxPanel {
	p := &TextboxPanel{}
	p.Panel.Init(p, parent, id)

	p.PlainText = NewTextbox(p, "plainText")
	p.PlainText.SetLabel("Plain Text")
	p.IntegerText = NewIntegerTextbox(p, "intText")
	p.IntegerText.SetLabel("Integer Text")
	p.IntegerText.SetMinValue(5, "")

	p.FloatText = NewFloatTextbox(p, "floatText")
	p.FloatText.SetLabel("Float Text")
	p.FloatText.SetMaxValue(6, "Hey this must be less than 6")

	p.SubmitAjax = NewButton(p, "ajaxButton")
	p.SubmitAjax.SetText("Submit Ajax")
	p.SubmitAjax.OnSubmit(action.Ajax(p.ID(), AjaxSubmit))

	p.SubmitServer = NewButton(p, "serverButton")
	p.SubmitServer.SetText("Submit Server")
	p.SubmitServer.OnSubmit(action.Server(p.ID(), ServerSubmit))

	return p
}


func init() {
	browsertest.RegisterTestFunction("Plain Textbox", testPlain)
}

// testPlain exercises the plain text box
func testPlain(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).AddValue("control", "textbox").String()
	t.LoadUrl(myUrl)
	t.AssertEqual("plainText", t.JqueryAttribute("plainText", "id")) // a sanity check


	t.ChangeVal("plainText", "me")
	t.Click("ajaxButton")
	t.AssertEqual("me", t.JqueryValue("plainText"))

	t.Done("Complete")
	/*
		t.AssertEquals("A value is required", t.SelectorInnerText("#user-name_err"))*/
}
