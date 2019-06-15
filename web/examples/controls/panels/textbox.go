package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/datetime"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
	"github.com/goradd/goradd/web/examples/controls"
)

// shared
const controlsFormPath = "/goradd/examples/controls.g"

const (
	AjaxSubmit int = iota + 1
	ServerSubmit
	ButtonSubmit
	ResetStateSubmit
	ProxyClick
)

type TextboxPanel struct {
	Panel
	PlainText    *Textbox
	MultiText    *Textbox
	IntegerText  *IntegerTextbox
	FloatText    *FloatTextbox
	EmailText    *EmailTextbox
	PasswordText *Textbox
	SearchText   *Textbox
	DateTimeText *DateTextbox
	DateText     *DateTextbox
	TimeText     *DateTextbox

	SubmitAjax   *Button
	SubmitServer *Button
}

func NewTextboxPanel(ctx context.Context, parent page.ControlI) {
	p := &TextboxPanel{}
	p.Panel.Init(p, parent, "textboxPanel")

	p.PlainText = NewTextbox(p, "plainText")
	p.PlainText.SetLabel("Plain Text")
	p.PlainText.SaveState(ctx, true)

	p.MultiText = NewTextbox(p, "multiText")
	p.MultiText.SetLabel("Multi Text")
	p.MultiText.SetRowCount(2)
	p.PlainText.SaveState(ctx, true)

	p.IntegerText = NewIntegerTextbox(p, "intText")
	p.IntegerText.SetLabel("Integer Text")

	p.FloatText = NewFloatTextbox(p, "floatText")
	p.FloatText.SetLabel("Float Text")

	p.EmailText = NewEmailTextbox(p, "emailText")
	p.EmailText.SetLabel("Email Text")

	p.PasswordText = NewTextbox(p, "passwordText")
	p.PasswordText.SetLabel("Password")
	p.PasswordText.SetType(TextboxTypePassword)

	p.SearchText = NewTextbox(p, "searchText")
	p.SearchText.SetLabel("Search")
	p.SearchText.SetType(TextboxTypeSearch)

	p.DateTimeText = NewDateTextbox(p, "dateTimeText")
	p.DateTimeText.SetLabel("U.S. Date-time")
	p.DateText = NewDateTextbox(p, "dateText")
	p.DateText.SetFormat(datetime.EuroDate)
	p.DateText.SetLabel("Euro Date")
	p.TimeText = NewDateTextbox(p, "timeText")
	p.TimeText.SetFormat(datetime.UsTime)
	p.TimeText.SetLabel("U.S. Time")

	p.SubmitAjax = NewButton(p, "ajaxButton")
	p.SubmitAjax.SetText("Submit Ajax")
	p.SubmitAjax.OnSubmit(action.Ajax(p.ID(), ButtonSubmit))

	p.SubmitServer = NewButton(p, "serverButton")
	p.SubmitServer.SetText("Submit Server")
	p.SubmitServer.OnSubmit(action.Server(p.ID(), ButtonSubmit))

}

func (p *TextboxPanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case ButtonSubmit:
	}
}

func init() {
	browsertest.RegisterTestFunction("Textbox Ajax Submit", testTextboxAjaxSubmit)
	browsertest.RegisterTestFunction("Textbox Server Submit", testTextboxServerSubmit)
	controls.RegisterPanel("textbox", "Textboxes", NewTextboxPanel, 2)
}

func testTextboxAjaxSubmit(t *browsertest.TestForm) {
	testTextboxSubmit(t, "ajaxButton")

	t.Done("Complete")
}

func testTextboxServerSubmit(t *browsertest.TestForm) {
	testTextboxSubmit(t, "serverButton")

	t.Done("Complete")
}

// testTextboxSubmit does a variety of submits using the given button. We use this to double check the various
// results we might get after a submission, as well as nsure that the ajax and server submits produce
// the same results.
func testTextboxSubmit(t *browsertest.TestForm, btnName string) {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "textbox").SetValue("testing",1).String()
	f := t.LoadUrl(myUrl)
	btn := f.Page().GetControl(btnName)

	t.ChangeVal("plainText", "me")
	t.ChangeVal("multiText", "me\nyou")
	t.ChangeVal("intText", "me")
	t.ChangeVal("floatText", "me")
	t.ChangeVal("emailText", "me")
	t.ChangeVal("dateTimeText", "me")
	t.ChangeVal("dateText", "me")
	t.ChangeVal("timeText", "me")

	t.Click(btn)

	t.AssertEqual("me", t.ControlValue("plainText"))
	t.AssertEqual("me\nyou", t.ControlValue("multiText"))
	t.AssertEqual("me", t.ControlValue("intText"))
	t.AssertEqual("me", t.ControlValue("floatText"))
	t.AssertEqual("me", t.ControlValue("emailText"))
	t.AssertEqual("me", t.ControlValue("dateTimeText"))
	t.AssertEqual("me", t.ControlValue("dateText"))
	t.AssertEqual("me", t.ControlValue("timeText"))

	t.AssertEqual(true, t.HasClass("intText_ctl", "error"))
	t.AssertEqual(true, t.HasClass("floatText_ctl", "error"))
	t.AssertEqual(true, t.HasClass("emailText_ctl", "error"))
	t.AssertEqual(true, t.HasClass("dateText_ctl", "error"))
	t.AssertEqual(true, t.HasClass("timeText_ctl", "error"))
	t.AssertEqual(true, t.HasClass("dateTimeText_ctl", "error"))

	plainText := f.Page().GetControl("plainText").(*Textbox)
	intText := f.Page().GetControl("intText").(*IntegerTextbox)
	floatText := f.Page().GetControl("floatText").(*FloatTextbox)
	emailText := f.Page().GetControl("emailText").(*EmailTextbox)
	dateText := f.Page().GetControl("dateText").(*DateTextbox)
	timeText := f.Page().GetControl("timeText").(*DateTextbox)
	dateTimeText := f.Page().GetControl("dateTimeText").(*DateTextbox)

	plainText.SetInstructions("Sample instructions")
	t.ChangeVal("intText", 5)
	t.ChangeVal("floatText", 6.7)
	t.ChangeVal("emailText", "me@you.com")
	t.ChangeVal("dateText", "19/2/2018")
	t.ChangeVal("timeText", "4:59 am")
	t.ChangeVal("dateTimeText", "2/19/2018 4:23 pm")

	t.Click(btn)

	t.AssertEqual(5, intText.Int())
	t.AssertEqual(6.7, floatText.Float64())
	t.AssertEqual("me@you.com", emailText.Text())
	t.AssertEqual(datetime.NewDateTime("19/2/2018", datetime.EuroDate), dateText.Date())
	t.AssertEqual(datetime.NewDateTime("4:59 am", datetime.UsTime), timeText.Date())
	t.AssertEqual(datetime.NewDateTime("2/19/2018 4:23 pm", datetime.UsDateTime), dateTimeText.Date())

	t.AssertEqual(false, t.HasClass("intText_ctl", "error"))
	t.AssertEqual(false, t.HasClass("floatText_ctl", "error"))
	t.AssertEqual(false, t.HasClass("emailText_ctl", "error"))
	t.AssertEqual(false, t.HasClass("dateText_ctl", "error"))
	t.AssertEqual(false, t.HasClass("timeText_ctl", "error"))
	t.AssertEqual(false, t.HasClass("dateTimeText_ctl", "error"))
	t.AssertEqual("Sample instructions", t.InnerHtml("plainText_inst"))

	t.AssertEqual("plainText_lbl plainText", t.ControlAttribute("plainText", "aria-labelledby"))

	// Test SaveState
	f = t.LoadUrl(myUrl)
	plainText = f.Page().GetControl("plainText").(*Textbox)
	//multiText := f.Page().GetControl("multiText").(*Textbox)
	t.AssertEqual("me", plainText.Text())
}
