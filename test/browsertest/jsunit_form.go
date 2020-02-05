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

func (f *JsUnitForm)Init(ctx context.Context, formID string) {
	f.FormBase.Init(ctx, formID)
	f.AddRelatedFiles()

	f.Results = NewPanel(f, "results")

	f.RunButton = NewButton(f, "startButton")
	f.RunButton.SetText("Start Test")
	f.RunButton.OnSubmit(action.Javascript("goradd.jsUnit.run(goradd.testsuite, 'results')"))
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
	t.LoadUrl(JsUnitTestFormPath)
	t.Click("startButton")

	h := t.ControlInnerHtml("results")
	t.AssertEqual("Done", h)

	t.Done("Complete")
}

func init() {
	page.RegisterForm(JsUnitTestFormPath, &JsUnitForm{}, JsUnitTestFormId)
}
