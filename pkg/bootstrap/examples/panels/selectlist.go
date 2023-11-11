package panels

import (
	"context"
	. "github.com/goradd/goradd/pkg/bootstrap/control"
	"github.com/goradd/goradd/pkg/bootstrap/examples"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/list"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
	"github.com/goradd/html5tag"
	"strings"
)

type SelectListPanel struct {
	control.Panel
}

func NewSelectListPanel(ctx context.Context, parent page.ControlI) {
	p := new(SelectListPanel)
	p.Init(p, ctx, parent, "selectListPanel")

}

func (p *SelectListPanel) Init(self any, ctx context.Context, parent page.ControlI, id string) {
	itemList := []list.ListValue{
		{"First", 1},
		{"Second", 2},
		{"Third", 3},
		{"Fourth", 4},
		{"Fifth", 5},
		{"Sixth", 6},
		{"Seventh", 7},
		{"Eighth", 8},
	}
	p.Panel.Init(self, parent, id)

	p.AddControls(ctx,
		FormGroupCreator{
			Label: "Standard SelectList",
			Child: SelectListCreator{
				ID:      "singleSelectList",
				NilItem: "- Select One -",
				Items:   itemList,
				ControlOptions: page.ControlOptions{
					IsRequired: true,
				},
			},
		},
		FormGroupCreator{
			Label: "Dropdown",
			Child: DropdownCreator{
				ID:    "dropdown",
				Text:  "Select One",
				Items: itemList,
				ButtonAttributes: html5tag.Attributes{
					"class": "btn-outline-primary",
				},
			},
		},
		FormGroupCreator{
			Label: "SelectList With Size",
			Child: SelectListCreator{
				ID:    "selectListWithSize",
				Items: itemList,
				Size:  4,
				ControlOptions: page.ControlOptions{
					IsRequired: true,
				},
			},
		},
		FormFieldsetCreator{
			Legend:       "Radio List",
			Instructions: "A radio list",
			Child: RadioListCreator{
				ID:    "radioList1",
				Items: itemList,
			},
		},
		FormFieldsetCreator{
			Legend: "CheckboxList List",
			Child: CheckboxListCreator{
				ID:    "checklist1",
				Items: itemList,
			},
		},
		ButtonCreator{
			ID:       "ajaxButton",
			Text:     "Submit Ajax",
			OnSubmit: action.Do("selectListPanel", ButtonSubmit),
		},
		ButtonCreator{
			ID:       "serverButton",
			Text:     "Submit Server",
			OnSubmit: action.Do("selectListPanel", ButtonSubmit),
		},
	)

}

func (p *SelectListPanel) DoAction(ctx context.Context, a action.Params) {
	switch a.ID {
	case ButtonSubmit:
		GetFormFieldset(p, "checklist1-ff").SetInstructions(
			strings.Join(GetCheckboxList(p, "checklist1").SelectedIds(), ","))
	}
}

func testSelectListAjaxSubmit(t *browsertest.TestForm) {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "lists").String()
	t.LoadUrl(myUrl)

	testSelectListSubmit(t, "ajaxButton")

	t.Done("Complete")
}

func testSelectListServerSubmit(t *browsertest.TestForm) {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "lists").String()
	t.LoadUrl(myUrl)

	testSelectListSubmit(t, "serverButton")

	t.Done("Complete")
}

// testSelectListSubmit does a variety of submits using the given button. We use this to double check the various
// results we might get after a submission, as well as ensure that the ajax and server submits produce
// the same results.
func testSelectListSubmit(t *browsertest.TestForm, btnID string) {
	t.ChooseListValue("selectListWithSize", "2")
	t.Click(btnID)
	t.AssertEqual(true, t.ControlHasClass("singleSelectList-ff", "error"))

	t.WithForm(func(f page.FormI) {
		t.AssertEqual(2, GetSelectList(f, "selectListWithSize").IntValue())
	})

	t.ChooseListValue("singleSelectList", "1")
	t.ChooseListValue("selectListWithSize", "2")

	t.CheckGroup("radioList1", "3")

	t.Click(btnID)

	t.WithForm(func(f page.FormI) {
		t.AssertEqual(1, GetSelectList(f, "singleSelectList").IntValue())
		t.AssertEqual(2, GetSelectList(f, "selectListWithSize").IntValue())
		t.AssertEqual(3, GetRadioList(f, "radioList1").IntValue())
	})
}

func init() {
	examples.RegisterPanel("lists", "Lists", NewSelectListPanel, 3)
	page.RegisterControl(&SelectListPanel{})

	browsertest.RegisterTestFunction("Bootstrap Select List Do Submit", testSelectListAjaxSubmit)
	browsertest.RegisterTestFunction("Bootstrap Select List Server Submit", testSelectListServerSubmit)
}
