// Package test contains the test harness, which controls browser based tests.
// Tests should call RegisterTestFunction to register a particular test. These tests get presented to the user
// in the test form available at the address "/test", and the user can select a test and execute it.
// The form is also a repository for operations you can perform on the form being tested. A test generally should
// start with a call to LoadURL. Follow that with calls to control the form and check for expected results.
// page/test contains a variety of tests that serve to unit test the form framework.
package browsertest

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/datetime"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
	"log"
	"os"
	"runtime"
)


var testFormPageState string

const TestFormPath = "/test"
const TestFormId = "TestForm"

const (
	TestButtonAction = iota + 1
)

type TestForm struct {
	page.ΩFormBase
	TestList     *SelectList
	RunningLabel *Span
	RunButton    *Button
	Controller   *TestController
	currentLog   string
	failed		 bool
	currentFailed bool
	currentTestName string
}

func NewTestForm(ctx context.Context) page.FormI {
	f := &TestForm{}
	f.Init(ctx, f, TestFormPath, TestFormId)
	//f.Page().SetDrawFunction(LoginPageTmpl)
	f.AddRelatedFiles()
	f.createControls(ctx)
	testFormPageState = f.Page().StateID()

	grctx := page.GetContext(ctx)

	if _,ok := grctx.FormValue("all"); ok {
		f.testAllAndExit()
	}
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
	for name,_ := range tests {
		f.TestList.AddItem(name, name)
	}
}

func (f *TestForm) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case TestButtonAction:
		f.runSelectedTest()
	}
}

func (f *TestForm) runSelectedTest() {
	f.RunningLabel.SetText(f.TestList.SelectedItem().Label())
	name := f.TestList.SelectedItem().Value().(string)
	f.runTest(name)
}

func (f *TestForm) runTest(name string) (result string) {
	var testF testRunnerFunction
	var ok bool
	var done = make(chan string)

	f.currentLog = ""

	if testF, ok = tests[name]; !ok {
		return fmt.Sprintf("Test %s does not exist", name)
	}

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
			done <- f.currentLog
			close(done)
		}()
		testF(f)
		done <- f.currentLog
		close(done)
	} ()


	return <- done
}

// Log will send a message to the log. The message might not draw right away.
func (f *TestForm) Log(s string) {
	d := datetime.Now()
	s = d.Format(datetime.StampMicro) + ": " + s
	f.currentLog += s + "\n"
	f.Controller.logLine(s)
	//log.Debugf("Log line %s", s)
}

// Mark the successful end of testing with a message.
func (f *TestForm) Done(s string) {
	f.Log(s)
	f.Page().PushRedraw()
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
	if page.GetPageCache().Has(f.Controller.pagestate) {
		return page.GetPageCache().Get(f.Controller.pagestate).Form()
	}
	return nil
}

// AssertEqual will test that the two values are equal, and will error if they are not equal.
// The test will continue after this.
func (f *TestForm) AssertEqual(expected, actual interface{}) {
	if expected != actual {
		f.Error(fmt.Sprintf("AssertEqual failed. %v != %v.", expected, actual))
	}
}

// Error will cause the test to error, but will continue performing the test.
func (f *TestForm) Error(message string) {
	f.Log(fmt.Sprintf("*** Test %s erred: %s", f.currentTestName, message))
	f.failed = true
	f.currentFailed = true
}

// Fail will cause a test to stop with the given messages.
func (f *TestForm) Fatal(message string) {
	panic(fmt.Sprint(message))
}

func (f *TestForm) fail(message string, testName string) {
	f.Log(fmt.Sprintf("*** Test %s failed: %s", testName, message))
	f.Page().PushRedraw()
	f.failed = true
	f.currentFailed = true
}

// ChangeVal will change the value of a form object. It essentially calls the jQuery .val() function on
// the html object with the given id, followed by sending a change event to the object. This is not quite
// the same thing as what happens when a user changes a value, as text boxes may send input events, and change
// is fired on some objects only when losing focus. However, this will simulate changing a value adequately for most
// situations.
func (f *TestForm) ChangeVal(id string, val interface{}) {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d ChangeVal(%q, %q)`, file, line, id, val)
	f.Controller.changeVal(id, val, desc);
}

// Click sends a click event to the html object with the given id. Note that this is not the same as simulating a click
// but for buttons, it will essentially be the same thing. More complex web objects will need a different mechanism
// for clicking, likely a chromium driver or something similar.
func (f *TestForm) Click(id string) {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d Click(%q)`, file, line, id)
	f.Controller.click(id, desc);
}

//CallJqueryFunction will call the given function with the given parameters on the jQuery object
// specified by the id. It will return the javascript result of the function call.
func (f *TestForm) CallJqueryFunction(id string, funcName string, params []string) string {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d CallJqueryFunction(%q, %q, %q)`, file, line, id, funcName, params)
	return f.Controller.callJqueryFunction(id, funcName, params, desc)
}

/*
func (f *TestForm) TypeValue(id string, chars string) {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d CallJqueryFunction(%q, %q, %q)`, file, line, id, funcName, params)
	f.Controller.typeChars(id, chars)
}*/

func (f *TestForm) Focus(id string) {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d Focus(%q)`, file, line, id)
	f.Controller.focus(id, desc)
}

func (f *TestForm) CloseWindow() {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d CloseWindow()`, file, line)
	f.Controller.closeWindow(desc)
}



// GetTestForm returns the test form itself, if its loaded
func GetTestForm() page.FormI {
	if page.GetPageCache().Has(testFormPageState) {
		return page.GetPageCache().Get(testFormPageState).Form()
	}
	return nil
}

func (f *TestForm) testAllAndExit() {
	var done = make(chan int)

	f.currentLog = ""

	go func() {
		for testName,testF := range tests {
			go func() {
				defer func() {
					f.CloseWindow()
					if r := recover(); r != nil {
						switch v := r.(type) {
						case error:
							f.fail(v.Error(), testName)
						case string:
							f.fail(v, testName)
						default:
							f.fail("Unknown error", testName)
						}
					}
					done <- 1
				}()
				f.Log("Starting test: " + testName)
				f.currentTestName = testName
				f.currentFailed = true
				testF(f)
				if !f.currentFailed {
					f.Log(fmt.Sprintf("Test %s completed successfully.", testName))
				}
			}()

			<- done
		}
		close(done)
		if f.failed {
			f.Log("Failed.")
		} else {
			f.Log("All tests passed.")
		}
		log.Print(f.currentLog)

		if f.failed {
			log.Fatal("Test failed.")
		} else {
			os.Exit(0)
		}
	} ()
}

func init() {
	page.RegisterPage(TestFormPath, NewTestForm, TestFormId)
}
