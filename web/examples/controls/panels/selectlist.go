package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/button"
	"github.com/goradd/goradd/pkg/page/control/list"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
	"github.com/goradd/html5tag"
	"strings"
)

type SelectListPanel struct {
	Panel
}

func NewSelectListPanel(ctx context.Context, parent page.ControlI) {
	p := &SelectListPanel{}
	p.Self = p
	p.Init(ctx, parent, "selectListPanel")
}

func (p *SelectListPanel) Init(ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(parent, id)

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

	p.AddControls(ctx,
		FormFieldWrapperCreator{
			Label: "Standard SelectList",
			Child: list.SelectListCreator{
				ID:      "singleSelectList",
				NilItem: "- Select One -",
				Items:   itemList,
				ControlOptions: page.ControlOptions{
					IsRequired: true,
				},
			},
		},
		FormFieldWrapperCreator{
			Label: "SelectList With Size",
			Child: list.SelectListCreator{
				ID:    "selectListWithSize",
				Items: itemList,
				Size:  4,
				ControlOptions: page.ControlOptions{
					IsRequired: true,
				},
			},
		},
		FormFieldWrapperCreator{
			Label: "Rows Radio List",
			Child: list.RadioListCreator{
				ID:          "radioList1",
				Items:       itemList,
				ColumnCount: 2,
			},
		},
		FormFieldWrapperCreator{
			Label: "Columns Radio List",
			Child: list.RadioListCreator{
				ID:              "radioList2",
				Items:           itemList,
				ColumnCount:     2,
				LayoutDirection: LayoutColumn,
			},
		},
		FormFieldWrapperCreator{
			Label: "Scrolling Radio List",
			Child: list.RadioListCreator{
				ID:          "radioList3",
				Items:       itemList,
				IsScrolling: true,
				ControlOptions: page.ControlOptions{
					Styles: html5tag.Style{
						"height": "80px",
					},
				},
			},
		},
		FormFieldWrapperCreator{
			Label: "Multiselect List",
			Child: list.MultiselectListCreator{
				ID:    "multiselectList",
				Items: itemList,
				ControlOptions: page.ControlOptions{
					IsRequired: true,
				},
			},
		},
		FormFieldWrapperCreator{
			Label: "CheckboxList List",
			Child: list.CheckboxListCreator{
				ID:          "checklist1",
				Items:       itemList,
				ColumnCount: 2,
			},
		},
		button.ButtonCreator{
			ID:       "ajaxButton",
			Text:     "Submit Ajax",
			OnSubmit: action.Ajax("selectListPanel", ButtonSubmit),
		},
		button.ButtonCreator{
			ID:       "serverButton",
			Text:     "Submit Server",
			OnSubmit: action.Server("selectListPanel", ButtonSubmit),
		},
	)
}

func (p *SelectListPanel) DoAction(ctx context.Context, a action.Params) {
	switch a.ID {
	case ButtonSubmit:
		checklist1 := list.GetCheckboxList(p, "checklist1")
		checklistWrapper := GetFormFieldWrapper(p, "checklist1-ff")
		checklistWrapper.SetInstructions(strings.Join(checklist1.SelectedIds(), ","))
	}
}

func init() {
	browsertest.RegisterTestFunction("Select List Ajax Submit", testSelectListAjaxSubmit)
	browsertest.RegisterTestFunction("Select List Server Submit", testSelectListServerSubmit)
}

// testPlain exercises the plain text box
func testSelectListAjaxSubmit(t *browsertest.TestForm) {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "selectlist").SetValue("testing", 1).String()
	t.LoadUrl(myUrl)

	testSelectListSubmit(t, "ajaxButton")

	t.Done("Complete")
}

func testSelectListServerSubmit(t *browsertest.TestForm) {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "selectlist").SetValue("testing", 1).String()
	t.LoadUrl(myUrl)

	testSelectListSubmit(t, "serverButton")

	t.Done("Complete")
}

// testSelectListSubmit does a variety of submits using the given button. We use this to double check the various
// results we might get after a submission, as well as nsure that the ajax and server submits produce
// the same results.
func testSelectListSubmit(t *browsertest.TestForm, btnID string) {
	t.ChooseListValue("selectListWithSize", "2")
	t.Click(btnID)

	t.AssertEqual(true, t.ControlHasClass("singleSelectList-ff", "error"))
	t.AssertEqual(true, t.ControlHasClass("multiselectList-ff", "error"))

	t.WithForm(func(f page.FormI) {
		t.AssertEqual(2, list.GetSelectList(f, "selectListWithSize").IntValue())
	})
	t.ChooseListValue("singleSelectList", "1")
	t.ChooseListValue("selectListWithSize", "2")
	t.CheckGroup("radioList1", "3")
	t.CheckGroup("radioList2", "4")
	t.ChooseListValues("multiselectList", "5")

	t.Click(btnID)

	t.WithForm(func(f page.FormI) {
		t.AssertEqual(1, list.GetSelectList(f, "singleSelectList").IntValue())
		t.AssertEqual("2", list.GetSelectList(f, "selectListWithSize").Value())
		t.AssertEqual(3, list.GetRadioList(f, "radioList1").IntValue())
		t.AssertEqual("4", list.GetRadioList(f, "radioList2").Value())
		v := list.GetMultiselectList(f, "multiselectList").Value().([]string)
		t.AssertEqual("5", v[0])
	})
}

/*
	select1 := f.Page().GetControl("multiselectList").(*MultiselectList)
	checklist1 := f.Page().GetControl("checklist1").(*CheckboxList)

	t.Click(btn)

	t.AssertEqual(true, t.ControlHasClass("multiselectList_ctl", "error"))

	t.AssertNotNil(select1.Value())

	id1,_ := select1.GetItemByValue(1)
	id2,_ := select1.GetItemByValue(3)

	t.ChangeVal("multiselectList", []string{id1, id2})
	id1,_ = select1.GetItemByValue(2)
	id2,_ = select1.GetItemByValue(3)
	t.CheckGroup("checklist1", id1, id2)

	t.Click(btn)

*/

func init() {
	page.RegisterControl(&SelectListPanel{})
}
