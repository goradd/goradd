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
	TestAllAction
)

type TestForm struct {
	FormBase
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
		f.ExecuteJqueryFunction("trigger", "testall", page.PriorityLow)
		f.On(event.Event("testall"), action.Ajax(f.ID(), TestAllAction))
	}
	return f
}

func (form *TestForm) createControls(ctx context.Context) {
	form.Controller = NewTestController(form, "controller")


	form.TestList = NewSelectList(form, "test-list")
	form.TestList.SetLabel("Tests")
	form.TestList.SetAttribute("size", 10)

	form.RunningLabel = NewSpan(form, "running-label")

	form.RunButton = NewButton(form, "run-button")
	form.RunButton.SetText("Run Test")
	form.RunButton.On(event.Click(), action.Ajax(form.ID(), TestButtonAction))
	form.RunButton.SetValidationType(page.ValidateNone)
}

func (form *TestForm) LoadControls(ctx context.Context) {
	tests.Range(func(k string,v interface{}) bool {
		form.TestList.AddItem(k, k)
		return true
	})
}

func (form *TestForm) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case TestButtonAction:
		form.runSelectedTest()
	case TestAllAction:
		form.testAllAndExit()
	}
}

func (form *TestForm) runSelectedTest() {
	form.RunningLabel.SetText(form.TestList.SelectedItem().Label())
	name := form.TestList.SelectedItem().Value().(string)
	form.testOne(name)
}


// Log will send a message to the log. The message might not draw right away.
func (form *TestForm) Log(s string) {
	d := datetime.Now()
	s = d.Format(datetime.StampMicro) + ": " + s
	form.currentLog += s + "\n"
	form.Controller.logLine(s)
	//log.Debugf("Log line %s", s)
}

// Mark the successful end of testing with a message.
func (form *TestForm) Done(s string) {
	form.Log(s)
	form.Page().PushRedraw()
}


// LoadUrl will launch a new window controlled by the test form. It will wait for the
// new url to be loaded in the window, and if the new url contains a goradd form, it will return
// the form.
func (form *TestForm) LoadUrl(url string) page.FormI {
	form.Log("Loading url: " + url)
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d LoadUrl(%q)`, file, line, url)
	form.Controller.loadUrl(url, desc)
	return form.GetForm()
}

// GetForm returns the currently loaded form.
func (form *TestForm) GetForm() page.FormI {
	if page.GetPageCache().Has(form.Controller.pagestate) {
		return page.GetPageCache().Get(form.Controller.pagestate).Form()
	}
	return nil
}

func (form *TestForm) AssertNotNil(v interface{}) {
	_, file, line, _ := runtime.Caller(1)

	if v == nil { // TODO: Check for a nil in the value
		form.error(fmt.Sprintf("*** AssertNotNil failed. File: %s, Line: %d", file, line))
	}
}


// AssertEqual will test that the two values are equal, and will error if they are not equal.
// The test will continue after this.
func (form *TestForm) AssertEqual(expected, actual interface{}) {
	_, file, line, _ := runtime.Caller(1)

	if expected != actual {
		form.error(fmt.Sprintf("*** AssertEqual failed. %v != %v. File: %s, Line: %d", expected, actual, file, line))
	}
}

// Error will cause the test to error, but will continue performing the test.
func (form *TestForm) Error(message string) {
	_, file, line, _ := runtime.Caller(1)
	form.error(fmt.Sprintf("*** Test %s erred: %s, File: %s, Line: %s", form.currentTestName, message, file, line))
}


func (form *TestForm) error(message string) {
	form.Log(message)
	form.failed = true
	form.currentFailed = true
}

// Fail will cause a test to stop with the given messages.
func (form *TestForm) Fatal(message string) {
	panic(fmt.Sprint(message))
}

func (form *TestForm) fail(message string, testName string) {
	form.Log(fmt.Sprintf("*** Test %s failed: %s", testName, message))
	form.Page().PushRedraw()
	form.failed = true
	form.currentFailed = true
}

// ChangeVal will change the value of a form object. It essentially calls the jQuery .val() function on
// the html object with the given id, followed by sending a change event to the object. This is not quite
// the same thing as what happens when a user changes a value, as text boxes may send input events, and change
// is fired on some objects only when losing focus. However, this will simulate changing a value adequately for most
// situations.
func (form *TestForm) ChangeVal(id string, val interface{}) {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d ChangeVal(%q, %q)`, file, line, id, val)
	form.Controller.changeVal(id, val, desc)
}

func (form *TestForm) CheckControl(id string, val bool) {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d CheckControl(%q, %q)`, file, line, id, val)
	form.Controller.checkControl(id, val, desc)
}


// Click sends a click event to the html object with the given id. Note that this is not the same as simulating a click
// but for buttons, it will essentially be the same thing. More complex web objects will need a different mechanism
// for clicking, likely a chromium driver or something similar.
func (form *TestForm) Click(id string) {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d Click(%q)`, file, line, id)
	form.Controller.click(id, desc)
}

// CallJqueryFunction will call the given function with the given parameters on the jQuery object
// specified by the id. It will return the javascript result of the function call.
func (form *TestForm) CallJqueryFunction(id string, funcName string, params ...interface{}) interface{} {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d CallJqueryFunction(%q, %q, %q)`, file, line, id, funcName, params)
	return form.Controller.callJqueryFunction(id, funcName, params, desc)
}

// Value will call the jquery .val() function on the given html object and return the result.
func (form *TestForm) JqueryValue(id string) interface{} {
	return form.CallJqueryFunction(id, "val")
}

// Attribute will call the jquery .attr("attribute") function on the given html object looking for the given
// attribute name and will return the value.
func (form *TestForm) JqueryAttribute(id string, attribute string) interface{} {
	return form.CallJqueryFunction(id, "attr", attribute)
}

func (form *TestForm) HasClass(id string, needle string) bool {
	res := form.CallJqueryFunction(id, "hasClass", needle)
	return res.(bool)
}

func (form *TestForm) InnerHtml(id string) string {
	res := form.CallJqueryFunction(id, "html")
	return res.(string)
}


/*
func (f *TestForm) TypeValue(id string, chars string) {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d CallJqueryFunction(%q, %q, %q)`, file, line, id, funcName, params)
	f.Controller.typeChars(id, chars)
}*/

func (form *TestForm) Focus(id string) {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d Focus(%q)`, file, line, id)
	form.Controller.focus(id, desc)
}

func (form *TestForm) CloseWindow() {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d CloseWindow()`, file, line)
	form.Controller.closeWindow(desc)
}



// GetTestForm returns the test form itself, if its loaded
func GetTestForm() page.FormI {
	if page.GetPageCache().Has(testFormPageState) {
		return page.GetPageCache().Get(testFormPageState).Form()
	}
	return nil
}

func (form *TestForm) testAllAndExit() {
	var done = make(chan int)

	form.currentLog = ""

	go func() {
		tests.Range(func(testName string, v interface{}) bool {
			testF := v.(testRunnerFunction)
			go func() {
				defer func() {
					form.CloseWindow()
					if r := recover(); r != nil {
						switch v := r.(type) {
						case error:
							form.fail(v.Error(), testName)
						case string:
							form.fail(v, testName)
						default:
							form.fail("Unknown error", testName)
						}
					}
					done <- 1
				}()
				form.Log("Starting test: " + testName)
				form.currentTestName = testName
				form.currentFailed = true
				testF(form)
				if !form.currentFailed {
					form.Log(fmt.Sprintf("Test %s completed successfully.", testName))
				}
			}()

			<- done
			return true
		})
		close(done)
		if form.failed {
			form.Log("Failed.")
		} else {
			form.Log("All tests passed.")
		}
		log.Print(form.currentLog)

		if form.failed {
			log.Fatal("Test failed.")
		} else {
			os.Exit(0)
		}
	} ()
}

func (form *TestForm) testOne(testName string) {
	var done = make(chan int)

	form.currentLog = ""

	go func() {
		if i := tests.Get(testName); i != nil {
			testF := i.(testRunnerFunction)
			go func() {
				defer func() {
					form.CloseWindow()
					if r := recover(); r != nil {
						switch v := r.(type) {
						case error:
							form.fail(v.Error(), testName)
						case string:
							form.fail(v, testName)
						default:
							form.fail("Unknown error", testName)
						}
					}
					done <- 1
				}()
				form.Log("Starting test: " + testName)
				form.currentTestName = testName
				form.currentFailed = true
				testF(form)
			}()
		}

		<- done

		close(done)
		if form.failed {
			form.Log("Failed.")
		} else {
			form.Log("Succeeded.")
		}
		log.Print(form.currentLog)

	} ()
}


func init() {
	page.RegisterPage(TestFormPath, NewTestForm, TestFormId)
}
