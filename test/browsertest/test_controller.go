package browsertest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
	"path"
	"time"

	_ "github.com/goradd/goradd/test/browsertest/assets"
)

const StepTimeoutSeconds = 10

const (
	testStepAction = iota + 100
	testMarkerAction
)

func TestStepEvent() *event.Event {
	return event.NewEvent("teststep")
}

func TestMarkerEvent() *event.Event {
	return event.NewEvent("testmarker")
}

type stepItemType struct {
	Step int
	Err  string
}

type markerItemType struct {
	Marker string
	Err    string
}

type TestController struct {
	control.Panel
	pagestate        string
	stepTimeout      time.Duration     // number of seconds before a step should timeout
	stepChannel      chan stepItemType // probably will leak memory TODO: Close this before it is removed from page cache
	markerChannel    chan string       // probably will leak memory TODO: Close this before it is removed from page cache
	latestJsValue    interface{}       // A value returned for the jsValue function
	stepDescriptions []string
}

func NewTestController(parent page.ControlI, id string) *TestController {
	p := new(TestController)
	p.Self = p
	p.Init(parent, id)
	return p
}

func (p *TestController) Init(parent page.ControlI, id string) {
	p.Panel.Init(parent, id)
	p.Tag = "pre"
	p.stepChannel = make(chan stepItemType, 1)
	p.markerChannel = make(chan string, 1000)
	p.ParentForm().AddJavaScriptFile(path.Join(config.AssetPrefix, "goradd", "test", "js", "test_controller.js"), false, nil)
	// Use declarative attribute to attach javascript to the control
	p.SetDataAttribute("grWidget", "goradd.TestController")

	p.On(TestStepEvent(), action.Ajax(p.ID(), testStepAction))
	p.On(TestMarkerEvent(), action.Ajax(p.ID(), testMarkerAction))
	p.stepTimeout = StepTimeoutSeconds
}

func (p *TestController) logLine(line string) {
	p.ExecuteWidgetFunction("logLine", line)
}

// loadUrl loads the url and returns the pagestate of the new form, if a goradd form got loaded.
func (p *TestController) loadUrl(url string, description string) {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.ExecuteWidgetFunction("loadUrl", len(p.stepDescriptions), url)
	p.waitStep() // load function will wait until window is loaded before firing
}

func (p *TestController) Action(ctx context.Context, a action.Params) {
	switch a.ID {
	case testStepAction:
		stepItem := new(stepItemType)
		ok, err := a.EventValue(stepItem)
		if err != nil {
			panic(err)
		}
		if !ok {
			panic("no step data found")
		}

		p.stepChannel <- *stepItem
	case testMarkerAction:
		marker := a.EventValueString()
		select {
		case p.markerChannel <- marker:
		default: // If we overflow the marker channel, we just move on and assume nobody is listening
			// Maybe we should close the channel at this point and stop sending?
		}
	}

}

func (p *TestController) UpdateFormValues(ctx context.Context) {
	id := p.ID()
	grctx := page.GetContext(ctx)

	if v := grctx.CustomControlValue(id, "pagestate"); v != nil {
		p.pagestate = v.(string)
	}
	if v := grctx.CustomControlValue(id, "jsvalue"); v != nil {
		p.latestJsValue = v
	}

}

func (p *TestController) waitStep() {
	log.FrameworkDebugf("Waiting for step %d: %s", len(p.stepDescriptions), p.stepDescriptions[len(p.stepDescriptions)-1])
	p.ParentForm().(*TestForm).PushRedraw()
	for {
		select {
		case stepItem := <-p.stepChannel:
			if stepItem.Step == -1 {
				log.FrameworkDebugf("Received form open")
			} else if stepItem.Step < len(p.stepDescriptions) {
				log.FrameworkDebugf("Received old step: %d, wanted %d", stepItem.Step, len(p.stepDescriptions))
				continue // this is a return from a previous step that timed out. We want to ignore it.
			} else if stepItem.Err != "" {
				panic(stepItem.Err)
			}
		case <-time.After(p.stepTimeout * time.Second):
			panic(fmt.Errorf("test step timed out: %s", p.stepDescriptions[len(p.stepDescriptions)-1]))
		}
		log.FrameworkDebugf("Completed step %d: %s", len(p.stepDescriptions), p.stepDescriptions[len(p.stepDescriptions)-1])
		break // we successfully returned from the step
	}
}

func (p *TestController) waitMarker(desc string, expectedMarker string) {
	log.FrameworkDebugf("Waiting for marker %s: %s", expectedMarker, desc)
	p.ParentForm().(*TestForm).PushRedraw()
	for {
		select {
		case marker := <-p.markerChannel:
			if marker == expectedMarker {
				log.FrameworkDebug("Received marker ", marker)
			} else {
				log.FrameworkDebugf("Received unexpected marker: %s, wanted %s", marker, expectedMarker)
				continue
			}
		case <-time.After(p.stepTimeout * time.Second):
			panic(fmt.Errorf("test marker  %s timed out: %s", expectedMarker, desc))
		}
		break // we successfully received the marker
	}
}

func (p *TestController) changeVal(id string, val interface{}, description string) {
	p.stepDescriptions = append(p.stepDescriptions, description)
	s, _ := json.Marshal(val)
	p.ExecuteWidgetFunction("changeVal", len(p.stepDescriptions), id, string(s))
	p.waitStep()
}

func (p *TestController) checkControl(id string, val bool, description string) {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.ExecuteWidgetFunction("checkControl", len(p.stepDescriptions), id, val)
	p.waitStep()
}

// checks a control or controls from a control group, specifically for checkbox and radio groups
func (p *TestController) checkGroup(name string, vals []string, description string) {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.ExecuteWidgetFunction("checkGroup", len(p.stepDescriptions), name, vals)
	p.waitStep()
}

func (p *TestController) click(id string, description string) {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.ExecuteWidgetFunction("click", len(p.stepDescriptions), id)
	p.waitStep()
}

func (p *TestController) waitSubmit(desc string) {
	p.stepDescriptions = append(p.stepDescriptions, desc)
	p.waitStep()
}

func (p *TestController) callWidgetFunction(id string, funcName string, params []interface{}, description string) interface{} {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.latestJsValue = nil
	p.ExecuteWidgetFunction("callWidgetFunction", len(p.stepDescriptions), id, funcName, params)
	p.waitStep()
	return p.latestJsValue
}

func (p *TestController) getHtmlElementInfo(selector string, attribute string, description string) interface{} {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.latestJsValue = nil
	p.ExecuteWidgetFunction("getHtmlElementInfo", len(p.stepDescriptions), selector, attribute)
	p.waitStep()
	return p.latestJsValue
}

func (p *TestController) typeChars(id string, chars string, description string) {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.ExecuteWidgetFunction("typeChars", len(p.stepDescriptions), id, chars)
	p.waitStep()
}

func (p *TestController) focus(id string, description string) {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.ExecuteWidgetFunction("focus", len(p.stepDescriptions), id)
	p.waitStep()
}

func (p *TestController) closeWindow(description string) {
	p.stepDescriptions = append(p.stepDescriptions, description)
	p.ExecuteWidgetFunction("closeWindow", len(p.stepDescriptions))
	p.waitStep()
}
