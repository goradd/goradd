package page


import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/test/browser"

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
	f.AddRelatedFiles()

	f.PlainText = NewTextbox(f, "plain")
	return f
}


func TestPlain(t *browser.TestForm)  {
	t.LoadUrl(TestTextboxPath)
	f := t.GetForm().(*TestTextboxForm)
	t.AssertEqual(TestTextboxId, f.ID())
	t.AssertEqual("plain", f.PlainText.ID())
	t.Error("Bad boy")

	/*
	t.ChangeVal(f.UserName.ID(), "me")
	t.Click(f.Submit.ID())
	t.AssertEqual("me", f.UserName.Text())
	t.AssertEqual("A value is required", f.Password.ValidationMessage())
	t.AssertEqual("A value is required", t.CallJqueryFunction(f.Password.ID() + "_err", "text", nil))

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
	page.RegisterPage(TestTextboxPath, NewTestTextboxForm, TestTextboxId)

	browser.RegisterTestFunction("Plain Textbox", TestPlain)
}