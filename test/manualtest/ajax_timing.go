package manualtest

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/goradd/test/browsertest"
)

const AjaxTimingPath = "/goradd/test/AjaxTiming.g"
const AjaxTimingId = "AjaxTimingForm"

const (
	Txt1ChangeAction = iota + 1
	Txt1KeyUpAction
	ChkChangeAction
	Txt2ChangeAction
	BtnClickAction
)

/**
 * AjaxTimingForm
 *
 * Tests the timing of ajax events and our ability to record a change to a control. There is a bit of a race condition
 * that we need to get under control. For example, if the user clicks a checkbox that also has a "click" handler, and
 * inside that handler, tests the value of the checkbox, the user should see the new value and not the one before the
 * click. Other problems come up when a Click causes a Change or FocusOut, which in turn changes the focus. Typically
 * in javascript, that would cause the Click to be lost. Our javascript tries to accommodate this by queueing all the
 * event responses before processing them.
 *
 * This is currently a manual test. Using a browser tester like Selenium, we could maybe make this part of the regular
 * continuous integration test.
 */
type AjaxTimingForm struct {
	FormBase
	Txt1            *Textbox
	Txt1ChangeLabel *Span
	Txt1KeyUpLabel  *Span
	Chk             *Checkbox
	ChkLabel        *Span
	Txt2            *Textbox
	Btn             *Button
}

func (f *AjaxTimingForm) Init(ctx context.Context, id string) {
	f.FormBase.Init(ctx, id)
	f.createControls(ctx)
}

func (f *AjaxTimingForm) createControls(ctx context.Context) {
	f.Txt1 = NewTextbox(f, "changer")
	f.Txt1.SetColumnCount(30)
	f.Txt1.SetPlaceholder("TextBox KeyUp Test")
	f.Txt1.SetText("Change Me")
	f.Txt1.On(event.Change(), action.Ajax(f.ID(), Txt1ChangeAction))
	f.Txt1.On(event.KeyUp(), action.Ajax(f.ID(), Txt1KeyUpAction))

	f.Txt1ChangeLabel = NewSpan(f, "vc")
	f.Txt1ChangeLabel.SetText("Value after Change: ")

	f.Txt1KeyUpLabel = NewSpan(f, "vu")
	f.Txt1KeyUpLabel.SetText("Value after Key Up: ")

	f.Chk = NewCheckbox(f, "cb")
	f.Chk.SetText("Checkbox Test")
	f.Chk.On(event.Click(), action.Ajax(f.ID(), ChkChangeAction))

	f.ChkLabel = NewSpan(f, "cbv")
	f.ChkLabel.SetText("Value after Click: ")

	f.Txt2 = NewTextbox(f, "tb")
	f.Txt2.SetText("Change Me")
	f.Txt2.On(event.Change(), action.Focus(f.Txt1.ID()))

	f.Btn = NewButton(f, "submit")
	f.Btn.SetLabel("Submit")
	f.Btn.On(event.Click(), action.Ajax(f.ID(), BtnClickAction))
	f.Btn.SetValidationType(event.ValidateNone)
}

func (f *AjaxTimingForm) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case Txt1ChangeAction:
		f.Txt1ChangeLabel.SetText(f.Txt1.Text())
	case Txt1KeyUpAction:
		f.Txt1KeyUpLabel.SetText(f.Txt1.Text())
	case ChkChangeAction:
		f.ChkLabel.SetText(fmt.Sprintf("%t", f.Chk.Checked()))
	case BtnClickAction:
		f.Txt1ChangeLabel.SetText("Button Was Clicked")
	}

}

func TestForm(t *browsertest.TestForm) {
	log.Debug("AjaxTiming test")

	t.LoadUrl(AjaxTimingPath)
	f := t.ParentForm().(*AjaxTimingForm)
	t.AssertEqual(AjaxTimingId, f.ID())
	t.Focus(f.Txt1.ID())
	//t.TypeValue(f.Txt1.ID(), "m")

	t.Log("Complete")
	/*
		t.AssertEquals("A value is required", t.SelectorInnerText("#user-name_err"))*/
}

func init() {
	page.RegisterForm(AjaxTimingPath, &AjaxTimingForm{}, AjaxTimingId)
}
