package test

import (
	"context"
	"fmt"
	"github.com/spekary/goradd/log"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/action"
	. "github.com/spekary/goradd/page/control"
)


type testRunnerFunction func(*TestForm)

var tests = make(map[string]testRunnerFunction)

const TestFormPath = "/test"
const TestFormId = "TestForm"

const (
	TestButtonAction = iota + 1
)

type TestForm struct {
	page.FormBase
	TestList     *SelectList
	RunningLabel *Span
	RunButton    *Button
	Controller   *TestController
	currentLog   string
}

func NewTestForm(ctx context.Context) page.FormI {
	f := &TestForm{}
	f.Init(ctx, f, TestFormPath, TestFormId)
	//f.Page().SetDrawFunction(LoginPageTmpl)
	f.AddRelatedFiles()
	f.createControls(ctx)
	return f
}

func (f *TestForm) createControls(ctx context.Context) {
	f.Controller = NewTestController(f, "controller")


	f.TestList = NewSelectList(f, "test-list")
	f.TestList.SetLabel("Tests")
	f.TestList.SetAttribute("size", 10)

	f.RunningLabel = NewSpan(f, "running-label")

	f.RunButton = NewButton(f, "run-button")
	f.RunButton.SetText("Run Test")
	f.RunButton.SetIsPrimary(true)
	f.RunButton.OnClick(action.Ajax(f.ID(), TestButtonAction))
}

func (f *TestForm) LoadControls(ctx context.Context) {
	for name,testF := range tests {
		f.TestList.AddItem(name, testF)
	}
}

func (f *TestForm) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case TestButtonAction:
		f.runTest()
	}
}

func (f *TestForm) runTest() {

	f.RunningLabel.SetText(f.TestList.SelectedItem().Label())
	testF := f.TestList.SelectedItem().Value().(testRunnerFunction)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case error:
					f.Log(v.Error())
				case string:
					f.Log(v)
				default:
					f.Log("Unknown error")
				}
			}
		}()
		testF(f)
	} ()
}

// Log will send a message to the log. The message might not draw right away.
func (f *TestForm) Log(s string) {
	f.currentLog += s + "\n"
	f.Controller.LogLine(s)
	log.Debugf("Log line %s", s)
	f.Page().PushRedraw()
}

func RegisterTestFunction (name string, f testRunnerFunction) {
	tests[name] = f
}


func init() {
	page.RegisterPage(TestFormPath, NewTestForm, TestFormId)
}

// loadUrl will launch a new window controlled by the test form. It will wait for the
// new url to be loaded in the window, and if the new url contains a goradd form, it will prepare
// to return the form if you call GetForm.
func (f *TestForm) LoadUrl(url string) {
	f.Log("Loading url: " + url)
	f.Controller.loadUrl(url)
}

// GetForm returns the currently loaded form.
func (f *TestForm) GetForm() page.FormI {
	if page.GetPageCache().Has(f.Controller.formstate) {
		return page.GetPageCache().Get(f.Controller.formstate).Form()
	}
	return nil
}

func (f *TestForm) AssertEqual(expected, actual interface{}) {
	if expected != actual {
		f.Controller.LogLine(fmt.Sprintf("AssertEqual failed. %v != %v.", expected, actual))
	}
}

func (f *TestForm) ChangeVal(id string, val interface{}) {
	f.Controller.changeVal(id, val);
}

func (f *TestForm) Click(id string) {
	f.Controller.click(id);
}


