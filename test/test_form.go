// Package test contains the test harness, which controls browser based tests.
// Tests should call RegisterTestFunction to register a particular test. These tests get presented to the user
// in the test form available at the address "/test", and the user can select a test and execute it.
// The form is also a repository for operations you can perform on the form being tested. A test generally should
// start with a call to LoadURL. Follow that with calls to control the form and check for expected results.
// page/test contains a variety of tests that serve to unit test the form framework.
package test

import (
	"context"
	"fmt"
	"github.com/spekary/goradd/datetime"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/action"
	. "github.com/spekary/goradd/page/control"
	"github.com/spekary/goradd/page/event"
	"runtime"
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
	f.RunButton.On(event.Click(), action.Ajax(f.ID(), TestButtonAction))
	f.RunButton.SetValidationType(page.ValidateNone)
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
					f.Done(v.Error())
				case string:
					f.Done(v)
				default:
					f.Done("Unknown error")
				}
			}
		}()
		testF(f)
	} ()
}

// Log will send a message to the log. The message might not draw right away.
func (f *TestForm) Log(s string) {
	d := datetime.Now()
	s = d.Format(datetime.StampMicro) + ": " + s
	f.currentLog += s + "\n"
	f.Controller.logLine(s)
	//log.Debugf("Log line %s", s)
}

// Mark the end of testing with a message.
func (f *TestForm) Done(s string) {
	f.Log(s)
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
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d LoadUrl(%q)`, file, line, url)
	f.Controller.loadUrl(url, desc)
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
		f.Controller.logLine(fmt.Sprintf("AssertEqual failed. %v != %v.", expected, actual))
	}
}

func (f *TestForm) ChangeVal(id string, val interface{}) {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d ChangeVal(%q, %q)`, file, line, id, val)
	f.Controller.changeVal(id, val, desc);
}

func (f *TestForm) Click(id string) {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d Click(%q)`, file, line, id)
	f.Controller.click(id, desc);
}

func (f *TestForm) JqueryValue(id string, funcName string, params []string) string {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d JqueryValue(%q, %q, %q)`, file, line, id, funcName, params)
	return f.Controller.jqValue(id, funcName, params, desc)
}

/*
func (f *TestForm) TypeValue(id string, chars string) {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d JqueryValue(%q, %q, %q)`, file, line, id, funcName, params)
	f.Controller.typeChars(id, chars)
}*/

func (f *TestForm) Focus(id string) {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d Focus(%q)`, file, line, id)
	f.Controller.focus(id, desc)
}


