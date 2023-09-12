package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	. "github.com/goradd/goradd/pkg/page/control/button"
	. "github.com/goradd/goradd/pkg/page/control/textbox"
	time2 "github.com/goradd/goradd/pkg/time"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
	"time"
)

type TextboxPanel struct {
	Panel
}

func NewTextboxPanel(ctx context.Context, parent page.ControlI) {
	p := new(TextboxPanel)
	p.Init(p, ctx, parent, "textboxPanel")
}

func (p *TextboxPanel) Init(self any, ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, id)

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
			Label: "IntegerTextbox Text",
			Child: IntegerTextboxCreator{
				ID: "intText",
			},
		},
		FormFieldWrapperCreator{
			ID:    "floatText-ff",
			Label: "FloatTextbox Text",
			Child: FloatTextboxCreator{
				ID: "floatText",
			},
		},
		FormFieldWrapperCreator{
			ID:    "emailText-ff",
			Label: "EmailTextbox Text",
			Child: EmailTextboxCreator{
				ID: "emailText",
			},
		},
		FormFieldWrapperCreator{
			ID:    "passwordText-ff",
			Label: "PasswordTextbox",
			Child: PasswordTextboxCreator{
				ID: "passwordText",
			},
		},
		FormFieldWrapperCreator{
			ID:    "searchText-ff",
			Label: "Search",
			Child: TextboxCreator{
				ID:   "searchText",
				Type: SearchType,
			},
		},
		FormFieldWrapperCreator{
			ID:    "dateTimeText-ff",
			Label: "U.S. DateTextbox-time",
			Child: DateTextboxCreator{
				ID: "dateTimeText",
			},
		},
		FormFieldWrapperCreator{
			ID:    "dateText-ff",
			Label: "Euro DateTextbox",
			Child: DateTextboxCreator{
				ID:      "dateText",
				Formats: []string{time2.EuroDate},
			},
		},
		FormFieldWrapperCreator{
			ID:    "timeText-ff",
			Label: "U.S. Time",
			Child: DateTextboxCreator{
				ID:      "timeText",
				Formats: []string{time2.UsTime},
			},
		},
		ButtonCreator{
			ID:       "ajaxButton",
			Text:     "Submit Ajax",
			OnSubmit: action.Do("textboxPanel", ButtonSubmit),
		},
		ButtonCreator{
			ID:       "serverButton",
			Text:     "Submit Post",
			OnSubmit: action.Do("textboxPanel", ButtonSubmit).Post(),
		},
	)
}

func (p *TextboxPanel) DoAction(ctx context.Context, a action.Params) {
	switch a.ID {
	case ButtonSubmit:
	}
}

func init() {
	browsertest.RegisterTestFunction("Textbox Do Submit", testTextboxAjaxSubmit)
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
	t.ChangeVal("timeText", "4:59AM")
	t.ChangeVal("dateTimeText", "2/19/2018 4:23 pm")

	t.Click(btnID)

	t.WithForm(func(f page.FormI) {
		t.AssertEqual(5, GetIntegerTextbox(f, "intText").Int())
		t.AssertEqual(6.7, GetFloatTextbox(f, "floatText").Float64())
		t.AssertEqual("me@you.com", GetEmailTextbox(f, "emailText").Text())

		v, _ := time.Parse(time2.EuroDate, "19/2/2018")
		t2 := time2.As(GetDateTextbox(f, "dateText").Date(), time.FixedZone("", 0))
		t.AssertEqual(true, v.Equal(t2))

		v2, _ := time.Parse(time2.UsTime, "4:59 AM")
		t3 := time2.As(GetDateTextbox(f, "timeText").Date(), time.FixedZone("", 0))
		t.AssertEqual(true, v2.Equal(t3))

		v3, _ := time.Parse(time2.UsDateTime, "2/19/2018 4:23 PM")
		t4 := time2.As(GetDateTextbox(f, "dateTimeText").Date(), time.FixedZone("", 0))
		t.AssertEqual(true, v3.Equal(t4))
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
