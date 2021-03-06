package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/datetime"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
	"time"
)

type TextboxPanel struct {
	Panel
}

func NewTextboxPanel(ctx context.Context, parent page.ControlI) {
	p := &TextboxPanel{}
	p.Self = p
	p.Init(ctx, parent, "textboxPanel")
}

func (p *TextboxPanel) Init(ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(parent, "textboxPanel")

	p.Panel.AddControls(ctx,
		FormFieldWrapperCreator{
			ID:    "plainText-ff",
			Label: "Plain Text",
			Child: TextboxCreator{
				ID:        "plainText",
				SaveState: true,
			},
		},
		FormFieldWrapperCreator{
			ID:    "multiText-ff",
			Label: "Multi Text",
			Child: TextboxCreator{
				ID:        "multiText",
				SaveState: true,
				RowCount:  2,
			},
		},
		FormFieldWrapperCreator{
			ID:    "intText-ff",
			Label: "Integer Text",
			Child: IntegerTextboxCreator{
				ID: "intText",
			},
		},
		FormFieldWrapperCreator{
			ID:    "floatText-ff",
			Label: "Float Text",
			Child: FloatTextboxCreator{
				ID: "floatText",
			},
		},
		FormFieldWrapperCreator{
			ID:    "emailText-ff",
			Label: "Email Text",
			Child: EmailTextboxCreator{
				ID: "emailText",
			},
		},
		FormFieldWrapperCreator{
			ID:    "passwordText-ff",
			Label: "Password",
			Child: TextboxCreator{
				ID:   "passwordText",
				Type: TextboxTypePassword,
			},
		},
		FormFieldWrapperCreator{
			ID:    "searchText-ff",
			Label: "Search",
			Child: TextboxCreator{
				ID:   "searchText",
				Type: TextboxTypeSearch,
			},
		},
		FormFieldWrapperCreator{
			ID:    "dateTimeText-ff",
			Label: "U.S. Date-time",
			Child: DateTextboxCreator{
				ID: "dateTimeText",
			},
		},
		FormFieldWrapperCreator{
			ID:    "dateText-ff",
			Label: "Euro Date",
			Child: DateTextboxCreator{
				ID:     "dateText",
				Formats: []string{datetime.EuroDate},
			},
		},
		FormFieldWrapperCreator{
			ID:    "timeText-ff",
			Label: "U.S. Time",
			Child: DateTextboxCreator{
				ID:     "timeText",
				Formats: []string{datetime.UsTime},
			},
		},
		ButtonCreator{
			ID:       "ajaxButton",
			Text:     "Submit Ajax",
			OnSubmit: action.Ajax("textboxPanel", ButtonSubmit),
		},
		ButtonCreator{
			ID:       "serverButton",
			Text:     "Submit Server",
			OnSubmit: action.Server("textboxPanel", ButtonSubmit),
		},
	)
}

func (p *TextboxPanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case ButtonSubmit:
	}
}

func init() {
	browsertest.RegisterTestFunction("Textbox Ajax Submit", testTextboxAjaxSubmit)
	browsertest.RegisterTestFunction("Textbox Server Submit", testTextboxServerSubmit)
	page.RegisterControl(&TextboxPanel{})
}

func testTextboxAjaxSubmit(t *browsertest.TestForm) {
	testTextboxSubmit(t, "ajaxButton")

	t.Done("Complete")
}

func testTextboxServerSubmit(t *browsertest.TestForm) {
	testTextboxSubmit(t, "serverButton")

	t.Done("Complete")
}

// testTextboxSubmit does a variety of submits using the given button. We use this to double check the various
// results we might get after a submission, as well as nsure that the ajax and server submits produce
// the same results.
func testTextboxSubmit(t *browsertest.TestForm, btnID string) {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "textbox").SetValue("testing", 1).String()
	t.LoadUrl(myUrl)

	t.ChangeVal("plainText", "me")
	t.ChangeVal("multiText", "me\nyou")
	t.ChangeVal("intText", "me")
	t.ChangeVal("floatText", "me")
	t.ChangeVal("emailText", "me")
	t.ChangeVal("dateTimeText", "me")
	t.ChangeVal("dateText", "me")
	t.ChangeVal("timeText", "me")

	t.Click(btnID)

	t.AssertEqual("me", t.ControlValue("plainText"))
	t.AssertEqual("me\nyou", t.ControlValue("multiText"))
	t.AssertEqual("me", t.ControlValue("intText"))
	t.AssertEqual("me", t.ControlValue("floatText"))
	t.AssertEqual("me", t.ControlValue("emailText"))
	t.AssertEqual("me", t.ControlValue("dateTimeText"))
	t.AssertEqual("me", t.ControlValue("dateText"))
	t.AssertEqual("me", t.ControlValue("timeText"))

	t.AssertEqual(true, t.ControlHasClass("intText-ff", "error"))
	t.AssertEqual(true, t.ControlHasClass("floatText-ff", "error"))
	t.AssertEqual(true, t.ControlHasClass("emailText-ff", "error"))
	t.AssertEqual(true, t.ControlHasClass("dateText-ff", "error"))
	t.AssertEqual(true, t.ControlHasClass("timeText-ff", "error"))
	t.AssertEqual(true, t.ControlHasClass("dateTimeText-ff", "error"))

	t.WithForm(func(f page.FormI) {
		GetFormFieldWrapper(f, "plainText-ff").SetInstructions("Sample instructions")
	})
	t.ChangeVal("intText", 5)
	t.ChangeVal("floatText", 6.7)
	t.ChangeVal("emailText", "me@you.com")
	t.ChangeVal("dateText", "19/2/2018")
	t.ChangeVal("timeText", "4:59 am")
	t.ChangeVal("dateTimeText", "2/19/2018 4:23 pm")

	t.Click(btnID)

	t.WithForm(func(f page.FormI) {
		t.AssertEqual(5, GetIntegerTextbox(f, "intText").Int())
		t.AssertEqual(6.7, GetFloatTextbox(f, "floatText").Float64())
		t.AssertEqual("me@you.com", GetEmailTextbox(f, "emailText").Text())
		t.AssertEqual(true, datetime.NewDateTime("19/2/2018", datetime.EuroDate).Equal(GetDateTextbox(f, "dateText").Date().As(time.FixedZone("", 0))))
		t.AssertEqual(true, datetime.NewDateTime("4:59 am", datetime.UsTime).Equal(GetDateTextbox(f, "timeText").Date().As(time.FixedZone("", 0))))
		t.AssertEqual(true, datetime.NewDateTime("2/19/2018 4:23 pm", datetime.UsDateTime).Equal(GetDateTextbox(f, "dateTimeText").Date().As(time.FixedZone("", 0))))
	})

	t.AssertEqual(false, t.ControlHasClass("intText-ff", "error"))
	t.AssertEqual(false, t.ControlHasClass("floatText-ff", "error"))
	t.AssertEqual(false, t.ControlHasClass("emailText-ff", "error"))
	t.AssertEqual(false, t.ControlHasClass("dateText-ff", "error"))
	t.AssertEqual(false, t.ControlHasClass("timeText-ff", "error"))
	t.AssertEqual(false, t.ControlHasClass("dateTimeText-ff", "error"))
	t.AssertEqual("Sample instructions", t.ControlInnerHtml("plainText-ff_inst"))

	t.AssertEqual("plainText-ff_lbl plainText", t.ControlAttribute("plainText", "aria-labelledby"))

	// Test SaveState
	t.WithForm(func(f page.FormI) {
		t.AssertEqual("me", GetTextbox(f, "plainText").Text())
	})

}
