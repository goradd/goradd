package panels

import (
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
)

type CheckboxPanel struct {
	Panel
}

func (p *CheckboxPanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case ButtonSubmit:
		var sel string
		if GetRadioButton(p,"radio1").Checked() {
			sel = "radio1"
		} else if GetRadioButton(p,"radio2").Checked() {
			sel = "radio2"
		} else if GetRadioButton(p,"radio3").Checked() {
			sel = "radio3"
		}
		GetPanel(p,"infoPanel").SetText(sel)
	}
}


func NewCheckboxPanel(ctx context.Context, parent page.ControlI) {
	p := &CheckboxPanel{}
	p.Panel.Init(p, parent, "checkboxPanel")
	p.AddControls(ctx,
		FormFieldWrapperCreator{
			ID:"checkbox1-ff",
			Label:"Checkbox 1:",
			For:"checkbox1",
			Instructions:"These are instructions for checkbox 1",
			Child:CheckboxCreator{
				ID:"checkbox1",
				Text:"My text is before",
				LabelMode:html.LabelBefore,
			},
		},
		FormFieldWrapperCreator{
			ID:"checkbox2-ff",
			Label:"Checkbox 2:",
			For:"checkbox2",
			Instructions:"These are instructions for checkbox 2",
			Child:CheckboxCreator{
				ID:"checkbox2",
				Text:"My text is after, and is wrapping the control",
				LabelMode:html.LabelWrapAfter,
			},
		},
		RadioButtonCreator{
			ID:"radio1",
			Group:"mygroup",
			Text:"Here",
		},
		RadioButtonCreator{
			ID:"radio2",
			Group:"mygroup",
			Text:"There",
		},
		RadioButtonCreator{
			ID:"radio3",
			Group:"mygroup",
			Text:"Everywhere",
		},
		PanelCreator{
			ID:"infoPanel",
		},
		ButtonCreator{
			ID:       "ajaxButton",
			Text:     "Submit Ajax",
			OnSubmit: action.Ajax("checkboxPanel", ButtonSubmit),
		},
		ButtonCreator{
			ID:       "serverButton",
			Text:     "Submit Server",
			OnSubmit: action.Ajax("checkboxPanel", ButtonSubmit),
		},

	)
}

func init() {
	browsertest.RegisterTestFunction("Checkbox Ajax Submit", testCheckboxAjaxSubmit)
	browsertest.RegisterTestFunction("Checkbox Server Submit", testCheckboxServerSubmit)
}

// testPlain exercises the plain text box
func testCheckboxAjaxSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "checkbox").SetValue("testing", 1).String()
	f := t.LoadUrl(myUrl)

	testCheckboxSubmit(t, f, f.Page().GetControl("ajaxButton"))

	t.Done("Complete")
}

func testCheckboxServerSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "checkbox").SetValue("testing", 1).String()
	f := t.LoadUrl(myUrl)

	testCheckboxSubmit(t, f, f.Page().GetControl("serverButton"))

	t.Done("Complete")
}

// testCheckboxSubmit does a variety of submits using the given button. We use this to double check the various
// results we might get after a submission, as well as nsure that the ajax and server submits produce
// the same results.
func testCheckboxSubmit(t *browsertest.TestForm, f page.FormI, btn page.ControlI) {
	checkbox1 := f.Page().GetControl("checkbox1").(*Checkbox)

	t.SetCheckbox("checkbox1", true)
	t.SetCheckbox("radio2", true)
	t.Click(btn)
	t.AssertEqual(true, checkbox1.Checked())

	checkbox2 := GetCheckbox(f, "checkbox2")
	radio1 := GetRadioButton(f,"radio1")
	radio2 := GetRadioButton(f,"radio2")
	radio3 := GetRadioButton(f,"radio3")
	info := GetPanel(f,"infoPanel")

	t.AssertEqual(false, checkbox2.Checked())
	t.AssertEqual("checkbox1-ff_lbl checkbox1_ilbl", t.ControlAttribute("checkbox1", "aria-labelledby"))

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

func init() {
	gob.Register(CheckboxPanel{})
}