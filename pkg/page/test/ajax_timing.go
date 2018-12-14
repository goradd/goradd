package test

import (
	"context"
	"fmt"
	"github.com/spekary/goradd/pkg/log"
	"github.com/spekary/goradd/pkg/page"
	"github.com/spekary/goradd/pkg/page/action"
	. "github.com/spekary/goradd/pkg/page/control"
	"github.com/spekary/goradd/pkg/page/event"
	"github.com/spekary/goradd/test/browser"
)


const AjaxTimingPath = "/page/test/AjaxTiming.g"
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
 * in javascript, that would cause the Click to be lost. Our javascript tries to accomodate this by queueing all the
 * event responses before
 */
type AjaxTimingForm struct {
	FormBase
	Txt1            *Textbox
	Txt1ChangeLabel *Span
	Txt1KeyUpLabel  *Span
	Chk             *Checkbox
	ChkLabel        *Span
	Txt2            *Textbox
	Btn 			*Button
}

func NewAjaxTimingForm(ctx context.Context) page.FormI {
	f := &AjaxTimingForm{}
	f.Init(ctx, f, AjaxTimingPath, AjaxTimingId)
	f.AddRelatedFiles()
	f.createControls(ctx)

	return f
}

func (f *AjaxTimingForm) createControls(ctx context.Context) {
	f.Txt1 = NewTextbox(f, "")
	f.Txt1.SetColumnCount(30)
	f.Txt1.SetLabel("TextBox KeyUp Test")
	f.Txt1.SetText("Change Me")
	f.Txt1.On(event.Change(), action.Ajax(f.ID(), Txt1ChangeAction))
	f.Txt1.On(event.KeyUp(), action.Ajax(f.ID(), Txt1KeyUpAction))

	f.Txt1ChangeLabel = NewSpan(f, "")
	f.Txt1ChangeLabel.SetLabel("Value after Change: ")

	f.Txt1KeyUpLabel = NewSpan(f, "")
	f.Txt1KeyUpLabel.SetLabel("Value after Key Up: ")

	f.Chk = NewCheckbox(f, "")
	f.Chk.SetLabel("Checkbox Test")
	f.Chk.On(event.Click(), action.Ajax(f.ID(), ChkChangeAction))

	f.ChkLabel = NewSpan(f, "")
	f.ChkLabel.SetLabel("Value after Click: ")

	f.Txt2 = NewTextbox(f, "")
	f.Txt2.SetLabel("TextBox Refocus Test")
	f.Txt2.SetText("Change Me")
	f.Txt2.On(event.Change(), action.Focus(f.Txt1.ID()))

	f.Btn = NewButton(f, "")
	f.Btn.SetLabel("Submit")
	f.Btn.On(event.Click(), action.Ajax(f.ID(), BtnClickAction))
	f.Btn.SetValidationType(page.ValidateNone)
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

func TestForm(t *browser.TestForm)  {
	log.Debug("AjaxTiming test")

	t.LoadUrl(AjaxTimingPath)
	f := t.GetForm().(*AjaxTimingForm)
	t.AssertEqual(AjaxTimingId, f.ID())
	t.Focus(f.Txt1.ID())
	//t.TypeValue(f.Txt1.ID(), "m")

	t.Log("Complete")
	/*
		t.AssertEquals("A value is required", t.SelectorInnerText("#user-name_err"))*/
}



func init() {
	page.RegisterPage(AjaxTimingPath, NewAjaxTimingForm, AjaxTimingId)
	browser.RegisterTestFunction("AjaxTiming", TestForm)

}
