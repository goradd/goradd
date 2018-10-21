package test

import (
	"context"
	"fmt"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/javascript"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/action"
	"github.com/spekary/goradd/page/control"
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
	formstate         string
	currentStepNumber int
	stepTimeout       time.Duration	// number of seconds before a step should timeout
	stepChannel chan stepItemType	// probably will leak memory TODO: Close this before it is removed from page cache
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

	p.ParentForm().AddJavaScriptFile(config.GoraddDir() + "/test/assets/js/test_controller.js", false, nil)
	p.On(TestStepEvent(), action.Ajax(p.ID(), TestStepAction))
	p.stepTimeout = 3
}

func (p *TestController) PutCustomScript(ctx context.Context, response *page.Response) {

	script := fmt.Sprintf (`$j("#%s").testController();`, p.ID())
	response.ExecuteJavaScript(script, page.PriorityStandard)
}

func (p *TestController) LogLine(line string) {
	script := fmt.Sprintf (`$j("#%s").testController("logLine", %q);`, p.ID(), line)
	p.ParentForm().Response().ExecuteJavaScript(script, page.PriorityStandard)
}

func (p *TestController) LoadUrl(url string) {
	p.ExecuteJqueryFunction("testController", "loadUrl", p.currentStepNumber, url)
	p.waitStep()
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

	if v := ctx.CustomControlValue(id, "formstate"); v != nil {
		p.formstate = v.(string)
	}
}

func (p *TestController) waitStep() {
	p.currentStepNumber++
	p.ExecuteJqueryFunction("testController", "waitStep", p.currentStepNumber, page.PriorityFinal)
	for {
		select {
		case stepItem := <-p.stepChannel:
			if stepItem.Step != p.currentStepNumber {
				continue // this is a return from a previous step that timed out. We want to ignore it.
			}
			if stepItem.Err != "" {
				panic (stepItem.Err)
			}
		case <-time.After(p.stepTimeout * time.Second):
			panic ("test step timed out")
		}
		break // we successfully returned from the step
	}
}

func (p *TestController) changeVal(id string, val interface{}) {
	p.ExecuteJqueryFunction("testController", "changeVal", p.currentStepNumber, fmt.Sprintf("%v", val))
	p.waitStep()
}


func init() {
	page.RegisterAssetDirectory(config.GoraddDir() + "/test/assets", config.AssetPrefix + "test")
}
