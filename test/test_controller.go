package test

import (
	"context"
	"fmt"
	"github.com/spekary/goradd/pkg/html"
	"github.com/spekary/goradd/pkg/javascript"
	"github.com/spekary/goradd/pkg/log"
	"github.com/spekary/goradd/pkg/page"
	"github.com/spekary/goradd/pkg/page/action"
	"github.com/spekary/goradd/pkg/page/control"
	"goradd-project/config"
	"time"
)

const (
	TestStepAction = iota + 100
)

// rowSelectedEvent indicates that a row was selected from the SelectTable
type testStepEvent struct {
	page.Event
}

// RowSelected
func TestStepEvent() *testStepEvent {
	e := &testStepEvent{page.Event{JsEvent: "goradd.teststep"}}
	e.ActionValue(javascript.JsCode("ui")) // the error string and step
	return e
}

type stepItemType struct {
	Step int
	Err string
}

type  TestController struct {
	control.Panel
	pagestate         string
	stepTimeout       time.Duration	// number of seconds before a step should timeout
	stepChannel chan stepItemType	// probably will leak memory TODO: Close this before it is removed from page cache
	latestJsValue string // A valure returned for the jsValue function
	stepDescriptions []string
}

func NewTestController(parent page.ControlI, id string) *TestController {
	p := new(TestController)
	p.Init(parent, id)
	p.Tag = "pre"
	p.stepChannel = make(chan stepItemType, 1)
	return p
}

func (p *TestController) Init(parent page.ControlI, id string) {
	p.Panel.Init(p, parent, id)
	path, attr := config.JQueryUIPath()
	p.ParentForm().AddJavaScriptFile(path, false, html.NewAttributesFromMap(attr))

	p.ParentForm().AddJavaScriptFile(config.GoraddDir + "/test/assets/js/test_controller.js", false, nil)
	p.On(TestStepEvent(), action.Ajax(p.ID(), TestStepAction))
	p.stepTimeout = 3
}

func (p *TestController) PutCustomScript(ctx context.Context, response *page.Response) {

	script := fmt.Sprintf (`$j("#%s").testController();`, p.ID())
	response.ExecuteJavaScript(script, page.PriorityStandard)
}

func (p *TestController) logLine(line string) {
	script := fmt.Sprintf (`$j("#%s").testController("logLine", %q);`, p.ID(), line)
	p.ParentForm().Response().ExecuteJavaScript(script, page.PriorityStandard)
}

func (p *TestController) loadUrl(url string, description string) {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.ExecuteJqueryFunction("testController", "loadUrl", len(p.stepDescriptions), url)
	p.waitStep(); // load function will wait until window is loaded before firing
}

func (p *TestController) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case TestStepAction:
		stepItem := new(stepItemType)
		ok,err := a.EventValue(stepItem)
		if err != nil {panic(err)}
		if !ok {panic("no step data found")}

		p.stepChannel<-*stepItem
	}
}

func (p *TestController) UpdateFormValues(ctx *page.Context) {
	id := p.ID()

	if v := ctx.CustomControlValue(id, "pagestate"); v != nil {
		p.pagestate = v.(string)
	}
	if v := ctx.CustomControlValue(id, "jsvalue"); v != nil {
		p.latestJsValue = v.(string)
	}

}

func (p *TestController) waitStep() {
	log.FrameworkDebugf("Waiting for step %d: %s", len(p.stepDescriptions), p.stepDescriptions[len(p.stepDescriptions)-1])
	p.Page().PushRedraw()
	for {
		select {
		case stepItem := <-p.stepChannel:
			if stepItem.Step < len(p.stepDescriptions) {
				log.FrameworkDebugf("Received old step: %d, wanted %d", stepItem.Step, len(p.stepDescriptions))
				continue // this is a return from a previous step that timed out. We want to ignore it.
			}
			if stepItem.Err != "" {
				panic (stepItem.Err)
			}
	//	case <-time.After(p.stepTimeout * time.Second):
	//		panic (fmt.Errorf("test step timed out: %s", p.stepDescriptions[len(p.stepDescriptions) - 1] ))
		}
		log.FrameworkDebugf("Completed step %d: %s", len(p.stepDescriptions), p.stepDescriptions[len(p.stepDescriptions)-1])
		break // we successfully returned from the step
	}
}

func (p *TestController) changeVal(id string, val interface{}, description string) {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.ExecuteJqueryFunction("testController", "changeVal", len(p.stepDescriptions), id, fmt.Sprintf("%v", val))
	p.waitStep()
}

func (p *TestController) click(id string, description string) {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.ExecuteJqueryFunction("testController", "click", len(p.stepDescriptions), id)
	p.waitStep()
}

func (p *TestController) jqValue(id string, funcName string, params []string, description string) string {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.ExecuteJqueryFunction("testController", "jqValue", len(p.stepDescriptions), id, funcName, params)
	p.waitStep()
	return p.latestJsValue
}

func (p *TestController) typeChars(id string, chars string, description string) {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.ExecuteJqueryFunction("testController", "typeChars", len(p.stepDescriptions), id, chars)
	p.waitStep()
}

func (p *TestController) focus(id string, description string) {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.ExecuteJqueryFunction("testController", "focus", len(p.stepDescriptions), id)
	p.waitStep()
}





func init() {
	page.RegisterAssetDirectory(config.GoraddDir + "/test/assets", config.AssetPrefix + "test")
}
