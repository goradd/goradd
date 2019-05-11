package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
)

type CheckboxPanel struct {
	Panel
	Checkbox1   *Checkbox
	Checkbox2   *Checkbox

	Radio1		*RadioButton
	Radio2		*RadioButton
	Radio3		*RadioButton

	Info 		*Panel

	SubmitAjax      *Button
	SubmitServer    *Button
}

func NewCheckboxPanel(ctx context.Context, parent page.ControlI) {
	p := &CheckboxPanel{}
	p.Panel.Init(p, parent, "checkboxPanel")

	p.Checkbox1 = NewCheckbox(p, "checkbox1")
	p.Checkbox1.SetLabel("Checkbox 1:")
	p.Checkbox1.SetText("My text is before")
	p.Checkbox1.SetLabelDrawingMode(html.LabelBefore)

	p.Checkbox2 = NewCheckbox(p, "checkbox2")
	p.Checkbox2.SetLabel("Checkbox 2:")
	p.Checkbox2.SetLabelDrawingMode(html.LabelWrapAfter)
	p.Checkbox2.SetText("My text is after, and is wrapping the control")

	p.Radio1 = NewRadioButton(p, "radio1")
	p.Radio1.SetGroup("mygroup")
	p.Radio1.SetText("Here")

	p.Radio2 = NewRadioButton(p, "radio2")
	p.Radio2.SetGroup("mygroup")
	p.Radio2.SetText("There")

	p.Radio3 = NewRadioButton(p, "radio3")
	p.Radio3.SetGroup("mygroup")
	p.Radio3.SetText("Everywhere")

	p.Info = NewPanel(p, "infoPanel")

	p.SubmitAjax = NewButton(p, "ajaxButton")
	p.SubmitAjax.SetText("Submit Ajax")
	p.SubmitAjax.OnSubmit(action.Ajax(p.ID(), ButtonSubmit))

	p.SubmitServer = NewButton(p, "serverButton")
	p.SubmitServer.SetText("Submit Server")
	p.SubmitServer.OnSubmit(action.Server(p.ID(), ButtonSubmit))
}


func init() {
	browsertest.RegisterTestFunction("Checkbox Ajax Submit", testCheckboxAjaxSubmit)
	browsertest.RegisterTestFunction("Checkbox Server Submit", testCheckboxServerSubmit)
}

// testPlain exercises the plain text box
func testCheckboxAjaxSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).AddValue("control", "checkbox").AddValue("testing", 1).String()
	f := t.LoadUrl(myUrl)

	testCheckboxSubmit(t, f, f.Page().GetControl("ajaxButton"))

	t.Done("Complete")
}

func testCheckboxServerSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).AddValue("control", "checkbox").AddValue("testing", 1).String()
	f := t.LoadUrl(myUrl)

	testCheckboxSubmit(t, f, f.Page().GetControl("serverButton"))

	t.Done("Complete")
}

// testCheckboxSubmit does a variety of submits using the given button. We use this to double check the various
// results we might get after a submission, as well as nsure that the ajax and server submits produce
// the same results.
func testCheckboxSubmit(t *browsertest.TestForm, f page.FormI, btn page.ControlI) {
	t.SetCheckbox("checkbox1", true)
	t.SetCheckbox("radio2", true)

	t.Click(btn)

	checkbox1 := f.Page().GetControl("checkbox1").(*Checkbox)
	checkbox2 := f.Page().GetControl("checkbox2").(*Checkbox)

	radio1 := f.Page().GetControl("radio1").(*RadioButton)
	radio2 := f.Page().GetControl("radio2").(*RadioButton)
	radio3 := f.Page().GetControl("radio3").(*RadioButton)

	info := f.Page().GetControl("infoPanel").(*Panel)


	t.AssertEqual("checkbox1_lbl checkbox1_ilbl", t.JqueryAttribute("checkbox1", "aria-labelledby"))
	t.AssertEqual(true, checkbox1.Checked())
	t.AssertEqual(false, checkbox2.Checked())

	t.AssertEqual(false, radio1.Checked())
	t.AssertEqual(true, radio2.Checked())

	t.AssertEqual("radio2", info.Text())

	t.SetCheckbox("radio3", true)
	t.SetCheckbox("checkbox1", false)
	t.Click(btn)
	t.AssertEqual(false, checkbox1.Checked())
	t.AssertEqual("radio3", info.Text())

	radio1.SetChecked(true);
	t.Click(btn)
	t.Click(btn) // two clicks are required to get the response back
	t.AssertEqual("radio1", info.Text())
	t.AssertEqual(true, radio1.Checked())
	t.AssertEqual(false, radio2.Checked())
	t.AssertEqual(false, radio3.Checked())

}

func (p *CheckboxPanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case ButtonSubmit:
		var sel string
		if p.Radio1.Checked() {
			sel = p.Radio1.ID()
		} else if  p.Radio2.Checked() {
			sel = p.Radio2.ID()
		} else if  p.Radio3.Checked() {
			sel = p.Radio3.ID()
		}
		p.Info.SetText(sel)
	}
}
