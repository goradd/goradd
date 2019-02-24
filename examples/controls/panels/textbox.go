package panels

import (
	"github.com/goradd/goradd/pkg/datetime"
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
	MultiText   *Textbox
	IntegerText *IntegerTextbox
	FloatText   *FloatTextbox
	EmailText 	*EmailTextbox
	PasswordText *Textbox
	SearchText *Textbox
	DateTimeText *DateTextbox
	DateText 	*DateTextbox
	TimeText 	*DateTextbox

	SubmitAjax      *Button
	SubmitServer    *Button
}

func NewTextboxPanel(parent page.ControlI) *TextboxPanel {
	p := &TextboxPanel{}
	p.Panel.Init(p, parent, "textboxPanel")

	p.PlainText = NewTextbox(p, "plainText")
	p.PlainText.SetLabel("Plain Text")

	p.MultiText = NewTextbox(p, "multiText")
	p.MultiText.SetLabel("Multi Text")
	p.MultiText.SetRowCount(2)

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
	p.SubmitAjax.OnSubmit(action.Ajax(p.ID(), AjaxSubmit))

	p.SubmitServer = NewButton(p, "serverButton")
	p.SubmitServer.SetText("Submit Server")
	p.SubmitServer.OnSubmit(action.Server(p.ID(), ServerSubmit))

	return p
}


func init() {
	browsertest.RegisterTestFunction("Textbox Ajax Submit", testAjaxSubmit)
	browsertest.RegisterTestFunction("Textbox Server Submit", testServerSubmit)
}

// testPlain exercises the plain text box
func testAjaxSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).AddValue("control", "textbox").String()
	f := t.LoadUrl(myUrl)

	testSubmit(t, f, "ajaxButton")

	t.Done("Complete")
}

func testServerSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).AddValue("control", "textbox").String()
	f := t.LoadUrl(myUrl)

	testSubmit(t, f, "serverButton")

	t.Done("Complete")
}

// testSubmit does a variety of submits using the given button. We use this to double check the various
// results we might get after a submission, as well as nsure that the ajax and server submits produce
// the same results.
func testSubmit(t *browsertest.TestForm, f page.FormI, btn string) {
	t.ChangeVal("plainText", "me")
	t.ChangeVal("multiText", "me")
	t.ChangeVal("intText", "me")
	t.ChangeVal("floatText", "me")
	t.ChangeVal("emailText", "me")
	t.ChangeVal("dateTimeText", "me")
	t.ChangeVal("dateText", "me")
	t.ChangeVal("timeText", "me")

	t.Click(btn)

	t.AssertEqual("me", t.JqueryValue("plainText"))
	t.AssertEqual("me", t.JqueryValue("multiText"))
	t.AssertEqual("me", t.JqueryValue("intText"))
	t.AssertEqual("me", t.JqueryValue("floatText"))
	t.AssertEqual("me", t.JqueryValue("emailText"))
	t.AssertEqual("me", t.JqueryValue("dateTimeText"))
	t.AssertEqual("me", t.JqueryValue("dateText"))
	t.AssertEqual("me", t.JqueryValue("timeText"))

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

}

