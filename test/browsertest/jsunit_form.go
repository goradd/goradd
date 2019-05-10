package browsertest

import (
	"context"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"path/filepath"
)

const JsUnitTestFormPath = "/goradd/test/jsunit.g"
const JsUnitTestFormId = "JsUnitTestForm"


type JsUnitForm struct {
	FormBase
	Results   *Panel
	RunButton *Button
}

func NewJsUnitForm(ctx context.Context) page.FormI {
	f := &JsUnitForm{}
	f.Init(ctx, f, JsUnitTestFormPath, JsUnitTestFormId)
	f.AddRelatedFiles()

	f.Results = NewPanel(f, "results")

	f.RunButton = NewButton(f, "startButton")
	f.RunButton.SetText("Start Test")
	f.RunButton.OnSubmit(action.Javascript("goradd.jsUnit.run(goradd.testsuite, 'results')"))
	return f
}

func (f *JsUnitForm) AddRelatedFiles() {
	f.FormBase.AddRelatedFiles()
	f.AddJavaScriptFile(filepath.Join(config.GoraddAssets(), "js", "goradd-js-unit.js"), false, nil)
	f.AddJavaScriptFile(filepath.Join(TestAssets(), "js", "goradd-js-unit-suite.js"), false, nil)
}


func init() {
	RegisterTestFunction("Goradd JavaScript Unit Tests", testJsUnit)
}

func testJsUnit(t *TestForm)  {
	f := t.LoadUrl(JsUnitTestFormPath)
	btn := f.Page().GetControl("startButton").(*Button)
	t.Click(btn)

	h := t.InnerHtml("results")
	t.AssertEqual("Done", h)

	t.Done("Complete")
}

func init() {
	page.RegisterPage(JsUnitTestFormPath, NewJsUnitForm, JsUnitTestFormId)
}
