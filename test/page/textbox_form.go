package page


import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/test/browser"
	"site"

	. "github.com/goradd/goradd/pkg/page/control"
)

const TestTextboxPath = "/test/textbox.g"
const TestTextboxId = "TestTextboxForm"

const (
	TestButtonAction = iota + 1
)

type TestTextboxForm struct {
	FormBase
	PlainText     *Textbox
	IntegerText     *IntegerTextbox
	FloatText 	*FloatTextbox
	Submit       *Button
}

func NewTestTextboxForm(ctx context.Context) page.FormI {
	f := &TestTextboxForm{}
	f.Init(ctx, f, TestTextboxPath, TestTextboxId)

	return f
}


/*
func RunLoginSuite() {
	TestLoginLaunch(t)
	Test1()
}

func TestLoginLaunch(t *test.TestForm) {
	t.LaunchBrowser("/")
	f := t.getForm("LoginForm")
}
*/


// Wrap this in panic catcher
// Possibly turn this into saved commands. Would help with logging.
func TestPasswordBlank(t *browser.TestForm)  {
	t.Log("Start TestPasswordBlank")
	t.LoadUrl("/")
	f := t.GetForm().(*site.LoginForm)
	t.AssertEqual("LoginForm", f.ID())
	t.ChangeVal(f.UserName.ID(), "me")
	t.Click(f.Submit.ID())
	t.AssertEqual("me", f.UserName.Text())
	t.AssertEqual("A value is required", f.Password.ValidationMessage())
	t.AssertEqual("A value is required", t.JqueryValue(f.Password.ID() + "_err", "text", nil))

	t.ChangeVal(f.Password.ID(), "me")
	t.Click(f.Submit.ID())
	t.AssertEqual("User not found", f.UserName.ValidationMessage())

	t.ChangeVal(f.UserName.ID(), "bob")
	t.Click(f.Submit.ID())
	t.AssertEqual("Password does not match", f.Password.ValidationMessage())

	t.Done("Complete")
	/*
		t.AssertEquals("A value is required", t.SelectorInnerText("#user-name_err"))*/
}

func init() {
	browser.RegisterTestFunction("PasswordBlank", TestPasswordBlank)
}