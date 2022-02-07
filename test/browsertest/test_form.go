// Package test contains the test harness, which controls browser based tests.
// Tests should call RegisterTestFunction to register a particular test. These tests get presented to the user
// in the test form available at the address "/test", and the user can select a test and execute it.
// The form is also a repository for operations you can perform on the form being tested. A test generally should
// start with a call to LoadURL. Follow that with calls to control the form and check for expected results.
package browsertest

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/datetime"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/messageServer"
	event2 "github.com/goradd/goradd/pkg/messageServer/event"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
	log2 "log"
	"os"
	"reflect"
	"runtime"
	"strings"
)

var testFormPageState string

const TestFormPath = "/goradd/Test.g"
const TestFormId = "TestForm"

const (
	TestButtonAction = iota + 1
	TestAllAction
)

type TestForm struct {
	FormBase
	Controller      *TestController
	currentLog      string
	failed          bool
	currentFailed   bool
	currentTestName string
	callerInfo      string
	usingForm		bool
}

func (form *TestForm) Init(ctx context.Context, formID string) {
	form.FormBase.Init(ctx, formID)
	//f.Page().SetDrawFunction(LoginPageTmpl)
	form.AddRelatedFiles()
	form.createControls(ctx)
	form.WatchChannel(ctx, "redraw")
	testFormPageState = form.Page().StateID()

	grctx := page.GetContext(ctx)

	if _, ok := grctx.FormValue("all"); ok {
		form.On(event2.MessengerReady(), action.Ajax(form.ID(), TestAllAction))
	}
}

func (form *TestForm) createControls(ctx context.Context) {
	form.Controller = NewTestController(form, "controller")

	NewSelectList(form, "test-list").
		SetAttribute("size", 10)

	NewSpan(form, "running-label")

	NewButton(form, "run-button").
		SetValidationType(page.ValidateNone).
		SetText("Run Test").
		On(event.Click(), action.Ajax(form.ID(), TestButtonAction))


	NewButton(form, "run-all-button").
		SetText("Run All Tests").
		SetValidationType(page.ValidateNone).
		On(event.Click(), action.Redirect(TestFormPath + "?all=1"))
}

func (form *TestForm) LoadControls(ctx context.Context) {
	tests.Range(func(k string, v interface{}) bool {
		GetSelectList(form, "test-list").AddItem(k, k)
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
	testList := GetSelectList(form, "test-list")
	GetSpan(form, "running-label").SetText(testList.SelectedItem().Label())
	name := testList.SelectedItem().Value()
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
	form.PushRedraw()
}

func (form *TestForm) PushRedraw() {
	messageServer.Messenger.Send("redraw", "U")
}

// LoadUrl will launch a new window controlled by the test form. It will wait for the
// new url to be loaded in the window, and if the new url contains a goradd form, it will return
// the form.
func (form *TestForm) LoadUrl(url string)  {
	form.Log("Loading url: " + url)
	form.Controller.loadUrl(url, form.captureCaller())
}

// getForm returns the currently loaded form.
func (form *TestForm) getForm() page.FormI {
	if page.GetPagestateCache().Has(form.Controller.pagestate) {
		pc := page.GetPagestateCache()
		/*if loader,ok := pc.(GetLoader); ok {
			p := loader.GetLoaded(form.Controller.pagestate)
			f :=  p.Form()
			return f
		}*/
		return pc.Get(form.Controller.pagestate).Form()
	}
	return nil
}

// WithForm gives you access to the current form so that you can set or get values in the form.
// Call it with a function that will receive the form.
// Do not call test functions that might cause an ajax or server call to fire from within the function.
func (form *TestForm) WithForm(f func(page.FormI) ) {
	pc := page.GetPagestateCache()
	testForm := pc.Get(form.Controller.pagestate).Form()
	{
		form.usingForm = true
		defer func(){ form.usingForm = false}()
		f(testForm)
	}
	pc.Set(form.Controller.pagestate, testForm.Page())
}


func (form *TestForm) AssertNil(v interface{}) {
	isNil := v == nil ||
		reflect.ValueOf(v).IsNil() // this will panic if value is not nillable, so be careful
	if !isNil{
		form.error(fmt.Sprintf("*** AssertNil failed. (%s)", form.captureCaller()))
	}
}

func (form *TestForm) AssertNotNil(v interface{}) {
	isNil := v == nil ||
		reflect.ValueOf(v).IsNil() // this will panic if value is not nillable, so be careful
	if isNil {
		form.error(fmt.Sprintf("*** AssertNotNil failed. (%s)", form.captureCaller()))
	}
}

// AssertEqual will test that the two values are equal, and will error if they are not equal.
// The test will continue after this.
func (form *TestForm) AssertEqual(expected, actual interface{}) {
	if expected != actual {
		form.error(fmt.Sprintf("*** AssertEqual failed. %v != %v. (%s)", expected, actual, form.captureCaller()))
	}
}

// Error will cause the test to error, but will continue performing the test.
func (form *TestForm) Error(message string) {
	form.error(fmt.Sprintf("*** Test %s erred: %s, (%s)", form.currentTestName, message, form.captureCaller()))
}

func (form *TestForm) error(message string) {
	form.Log(message)
	form.failed = true
	form.currentFailed = true
}

// Fatal will cause a test to stop with the given messages.
func (form *TestForm) Fatal(message string) {
	panic(fmt.Sprint(message))
}

func (form *TestForm) panicked(message string, testName string) {
	var panickingLine string
	if _, file, line, ok := runtime.Caller(4); ok {
		panickingLine = fmt.Sprintf("%s:%d", file, line)
	}
	msg := fmt.Sprintf("\n*** Test %s panicked: %s\n*** Last test step: %s\n*** Panicking line: %s", testName, message, form.callerInfo, panickingLine)
	log.Debug(msg)
	form.Log(msg)
	form.PushRedraw()
	form.failed = true
	form.currentFailed = true
}

// ChangeVal will change the value of a form object. It essentially calls the javascript .val() function on
// the html object with the given id, followed by sending a change event to the object. This is not quite
// the same thing as what happens when a user changes a value, as text boxes may send input events, and change
// is fired on some objects only when losing focus. However, this will simulate changing a value adequately for most
// situations.
func (form *TestForm) ChangeVal(id string, val interface{}) {
	if form.usingForm {
		panic("do not call ChangeVal from inside the F() function")
	}

	// TODO: Make sure that you don't call this from within an F() call if it has a change handler
	// attached to it.

	form.Controller.changeVal(id, val, form.captureCaller())
}

// SetCheckbox sets the given checkbox control to the given value. Use this instead of ChangeVal on checkboxes.
func (form *TestForm) SetCheckbox(id string, val bool) {
	if form.usingForm {
		panic("do not call SetCheckbox from inside the F() function")
	}
	// TODO: Make sure that you don't call this from within an F() call if there is a click
	//  or change handler attached to it. We should also see if there is a server click or change handler
	// attached to the control and wait for a reload if so. See Click() for example.
	form.Controller.checkControl(id, val, form.captureCaller())
}

func (form *TestForm) ChooseListValue(id string, value string) {
	if form.usingForm {
		panic("do not call SetListVal from inside the F() function")
	}

	form.ChangeVal(id, value)
}

func (form *TestForm) ChooseListValues(id string, values ...string) {
	if form.usingForm {
		panic("do not call SetListVal from inside the F() function")
	}

	form.ChangeVal(id, values)
}

// CheckGroup sets the checkbox group to a list of values. Radio groups should only be given one value
// to check. Will uncheck anything checked in the group before checking the given values. Specify nil
// to uncheck everything.
func (form *TestForm) CheckGroup(id string, values ...string) {
	if form.usingForm {
		panic("do not call CheckGroup from inside the F() function")
	}
	form.Controller.checkGroup(id, values, form.captureCaller())
}

// Click sends a click to the goradd control.
// Note that the act of clicking often causes an action, and an action will change the form
func (form *TestForm) Click(id string) {
	if form.usingForm {
		panic("do not call Click from inside the F() function")
	}
	form.Controller.click(id, form.captureCaller())
	f := form.getForm()
	c := f.Page().GetControl(id)
	if c.HasServerAction("click") {
		// wait for the new page to load
		form.Controller.waitSubmit(form.callerInfo)
	}
}

// ClickSubItem sends a click to the html object with the given sub-id inside the given control.
func (form *TestForm) ClickSubItem(id string, subId string) {
	form.Controller.click(id+"_"+subId, form.captureCaller())
	f := form.getForm()
	c := f.Page().GetControl(id)
	if c.HasServerAction("click") {
		// wait for the new page to load
		form.Controller.waitSubmit(form.callerInfo)
	}
}

func (form *TestForm) ClickHtmlItem(id string) {
	if form.usingForm {
		panic("do not call ClickHtmlItem from inside the F() function")
	}
	form.Controller.click(id, form.captureCaller())
}

func (form *TestForm) WaitSubmit() {
	form.Controller.waitSubmit(form.captureCaller())
}

// WaitMarker waits for a marker event before proceeding. You can use this to signal that your code has reached a specific
// place. To signal the marker, call FireTestMarker from a control in your application.
func (form *TestForm) WaitMarker(expectedMarker string) {
	form.Controller.waitMarker(form.captureCaller(), expectedMarker)
}


// CallControlFunction will call the given function with the given parameters on the goradd object
// specified by the id. It will return the javascript result of the function call.
func (form *TestForm) CallControlFunction(id string, funcName string, params ...interface{}) interface{} {
	return form.Controller.callWidgetFunction(id, funcName, params, form.captureCaller())
}

// ControlValue will call the .val() function on the given goradd object and return the result.
func (form *TestForm) ControlValue(id string) interface{} {
	return form.Controller.callWidgetFunction(id, "val", nil, form.captureCaller())

}

// ControlAttribute will call the .attr("attribute") function on the given html object looking for the given
// attribute name and will return the value.
func (form *TestForm) ControlAttribute(id string, attribute string) interface{} {
	return form.Controller.callWidgetFunction(id, "attr", []interface{}{attribute}, form.captureCaller())
}

// ControlHasClass will return true if the given goradd control has the given class attribute
func (form *TestForm) ControlHasClass(id string, needle string) bool {
	res := form.Controller.callWidgetFunction(id, "hasClass", []interface{}{needle}, form.captureCaller())
	return res.(bool)
}

// ControlInnerHtml will return the inner html drawn by a goradd control
func (form *TestForm) ControlInnerHtml(id string) string {
	res := form.Controller.callWidgetFunction(id, "html", nil, form.captureCaller())
	return strings.TrimSpace(res.(string)) // html can have a variety of inconsequential spaces
}

// HtmlElementInfo will return the inner html drawn by a goradd control
func (form *TestForm) HtmlElementInfo(selector string, attribute string) string {
	res := form.Controller.getHtmlElementInfo(selector, attribute, form.captureCaller())
	if res == nil {return ""}
	return strings.TrimSpace(res.(string))
}


/*
func (f *TestForm) TypeValue(id string, chars string) {
	_, file, line, _ := runtime.Caller(1)
	desc := fmt.Sprintf(`%s:%d CallControlFunction(%q, %q, %q)`, file, line, id, funcName, params)
	f.Controller.typeChars(id, chars)
}*/

func (form *TestForm) Focus(id string) {
	form.Controller.focus(id, form.captureCaller())
}

func (form *TestForm) CloseWindow() {
	form.Controller.closeWindow(form.captureCaller())
}

func (form *TestForm) captureCaller() string {
	if _, file, line, ok := runtime.Caller(2); ok {
		form.callerInfo = fmt.Sprintf(`%s:%d`, file, line)
	} else {
		form.callerInfo = "Unknown caller"
	}
	return form.callerInfo
}

// GetTestForm returns the test form itself, if its loaded
func GetTestForm() page.FormI {
	if page.GetPagestateCache().Has(testFormPageState) {
		return page.GetPagestateCache().Get(testFormPageState).Form()
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
					if r := recover(); r != nil {
						switch v := r.(type) {
						case error:
							form.panicked(v.Error(), testName)
						case string:
							form.panicked(v, testName)
						default:
							form.panicked("Unknown error", testName)
						}
					}
					done <- 1
				}()
				form.Log("Starting test: " + testName)
				form.currentTestName = testName
				form.currentFailed = false
				testF(form)
				if !form.currentFailed {
					form.Log(fmt.Sprintf("Test %s completed successfully.", testName))
				}
			}()

			<-done
			return true
		})
		close(done)
		if form.failed {
			form.Log("Failed.")
		} else {
			form.Log("All tests passed.")
		}
		log.Debug(form.currentLog)

		if form.failed {
			log2.Fatal("Test failed.")
		}
		os.Exit(0)
	}()
}

func (form *TestForm) testOne(testName string) {
	var done = make(chan int)

	form.currentLog = ""

	go func() {
		if i := tests.Get(testName); i != nil {
			testF := i.(testRunnerFunction)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						switch v := r.(type) {
						case error:
							form.panicked(v.Error(), testName)
						case string:
							form.panicked(v, testName)
						default:
							form.panicked("Unknown error", testName)
						}
					}
					done <- 1
				}()
				form.Log("Starting test: " + testName)
				form.currentTestName = testName
				form.currentFailed = false
				testF(form)
				//form.CloseWindow()
			}()
		}

		<-done

		close(done)
		if form.currentFailed {
			form.Log("Failed.")
		} else {
			form.Log("Succeeded.")
		}
		log.Debug(form.currentLog)

	}()
}

func (form *TestForm) NoSerialize() bool {
	return true
}

// Makes a sql context for the purpose of doing database tests
// TODO: This really needs to be a generalized context so NoSQL can be included
func MakeSqlContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, goradd.SqlContext, &db.SqlContext{}) // needed for transactions
	return ctx
}

func init() {
	page.RegisterForm(TestFormPath, &TestForm{}, TestFormId)
}
