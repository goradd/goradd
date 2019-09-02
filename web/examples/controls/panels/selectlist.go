package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
	"strings"
)

type SelectListPanel struct {
	Panel
}

func NewSelectListPanel(ctx context.Context, parent page.ControlI) {
	itemList := []ListValue{
		{"First", 1},
		{"Second", 2},
		{"Third", 3},
		{"Fourth", 4},
		{"Fifth", 5},
		{"Sixth", 6},
		{"Seventh", 7},
		{"Eighth", 8},
	}

	p := &SelectListPanel{}
	p.Panel.Init(p, parent, "selectListPanel")
	p.AddControls(ctx,
		FormFieldWrapperCreator{
			Label: "Standard SelectList",
			Child: SelectListCreator{
				ID: "singleSelectList",
				NilItem: "- Select One -",
				Items: itemList,
				ControlOptions:page.ControlOptions{
					IsRequired: true,
				},
			},
		},
		FormFieldWrapperCreator{
			Label: "SelectList With Size",
			Child: SelectListCreator{
				ID: "selectListWithSize",
				Items: itemList,
				Size: 4,
				ControlOptions:page.ControlOptions{
					IsRequired: true,
				},
			},
		},
		FormFieldWrapperCreator{
			Label: "Rows Radio List",
			Child: RadioListCreator{
				ID: "radioList1",
				Items: itemList,
				ColumnCount: 2,
			},
		},
		FormFieldWrapperCreator{
			Label: "Columns Radio List",
			Child: RadioListCreator{
				ID: "radioList2",
				Items: itemList,
				ColumnCount: 2,
				LayoutDirection: LayoutColumn,
			},
		},
		FormFieldWrapperCreator{
			Label: "Scrolling Radio List",
			Child: RadioListCreator{
				ID: "radioList3",
				Items: itemList,
				IsScrolling: true,
				ControlOptions:page.ControlOptions{
					Styles:html.Style {
						"height": "80px",
					},
				},
			},
		},
		FormFieldWrapperCreator{
			Label: "Multiselect List",
			Child: MultiselectListCreator{
				ID: "multiselectList",
				Items: itemList,
				ControlOptions:page.ControlOptions{
					IsRequired: true,
				},
			},
		},
		FormFieldWrapperCreator{
			Label: "Checkbox List",
			Child: CheckboxListCreator{
				ID: "checklist1",
				Items: itemList,
				ColumnCount:2,
			},
		},
		ButtonCreator{
			ID:       "ajaxButton",
			Text:     "Submit Ajax",
			OnSubmit: action.Ajax("selectListPanel", ButtonSubmit),
		},
		ButtonCreator{
			ID:       "serverButton",
			Text:     "Submit Server",
			OnSubmit: action.Ajax("selectListPanel", ButtonSubmit),
		},

	)
}

func (p *SelectListPanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case ButtonSubmit:
		checklist1 := GetCheckboxList(p, "checklist1")
		checklistWrapper := GetFormFieldWrapper(p, "checklist1-ff")
		checklistWrapper.SetInstructions(strings.Join(checklist1.SelectedIds(), ","))
	}
}


func init() {
	browsertest.RegisterTestFunction("Select List Ajax Submit", testSelectListAjaxSubmit)
	browsertest.RegisterTestFunction("Select List Server Submit", testSelectListServerSubmit)
}

// testPlain exercises the plain text box
func testSelectListAjaxSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "selectlist").SetValue("testing", 1).String()
	f := t.LoadUrl(myUrl)

	testSelectListSubmit(t, f, f.Page().GetControl("ajaxButton"))

	t.Done("Complete")
}

func testSelectListServerSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "selectlist").SetValue("testing", 1).String()
	f := t.LoadUrl(myUrl)

	testSelectListSubmit(t, f, f.Page().GetControl("serverButton"))

	t.Done("Complete")
}

// testSelectListSubmit does a variety of submits using the given button. We use this to double check the various
// results we might get after a submission, as well as nsure that the ajax and server submits produce
// the same results.
func testSelectListSubmit(t *browsertest.TestForm, f page.FormI, btn page.ControlI) {

	// For testing purposes, we need to use the id of the list item, rather than the value of the list item,
	// since that is what is presented in the html.
	select1 := GetSelectList(f, "singleSelectList")
	select2 := GetSelectList(f, "selectListWithSize")
	radio1 := GetRadioList(f, "radioList1")
	radio2 := GetRadioList(f, "radioList2")
	multi := GetMultiselectList(f, "multiselectList")

	id, _ := select2.GetItemByValue(2)
	t.ChangeVal("selectListWithSize", id)

	t.Click(btn)

	t.AssertEqual(true, t.HasClass("singleSelectList-ff", "error"))
	t.AssertEqual(true, t.HasClass("multiselectList-ff", "error"))

	t.AssertEqual(2, select2.IntValue())

	id, _ = select1.GetItemByValue(1)
	t.ChangeVal("singleSelectList", id)
	id, _ = select2.GetItemByValue(2)
	t.ChangeVal("selectListWithSize", id)
	id, _ = radio1.GetItemByValue(3)
	t.CheckGroup("radioList1", id)
	id, _ = radio2.GetItemByValue(4)
	t.CheckGroup("radioList2", id)
	id, _ = multi.GetItemByValue(5)
	t.ChangeVal("multiselectList", []string{id})

	t.Click(btn)

	t.AssertEqual(1, select1.IntValue())
	t.AssertEqual(2, select2.IntValue())
	t.AssertEqual(3, radio1.IntValue())
	t.AssertEqual(4, radio2.IntValue())
	v := multi.Value().([]interface{})
	t.AssertEqual(5, v[0])
}

/*
	select1 := f.Page().GetControl("multiselectList").(*MultiselectList)
	checklist1 := f.Page().GetControl("checklist1").(*CheckboxList)

	t.Click(btn)

	t.AssertEqual(true, t.HasClass("multiselectList_ctl", "error"))

	t.AssertNotNil(select1.Value())

	id1,_ := select1.GetItemByValue(1)
	id2,_ := select1.GetItemByValue(3)

	t.ChangeVal("multiselectList", []string{id1, id2})
	id1,_ = select1.GetItemByValue(2)
	id2,_ = select1.GetItemByValue(3)
	t.CheckGroup("checklist1", id1, id2)

	t.Click(btn)

*/
