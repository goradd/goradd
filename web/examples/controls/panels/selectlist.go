package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
	"strings"
)

type SelectListPanel struct {
	Panel
	SingleSelect   *SelectList
	SingleSelectWithSize   *SelectList
	RadioList1   *RadioList
	RadioList2   *RadioList
	RadioList3   *RadioList

	MultiSelect   *MultiselectList
	CheckboxList1   *CheckboxList

	SubmitAjax      *Button
	SubmitServer    *Button
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

	p.SingleSelect = NewSelectList(p, "singleSelectList")
	p.SingleSelect.SetLabel("Standard SelectList")
	p.SingleSelect.SetData(itemList)
	p.SingleSelect.AddItemAt(0, "- Select One -", nil)
	p.SingleSelect.SetIsRequired(true)

	p.SingleSelectWithSize = NewSelectList(p, "selectListWithSize")
	p.SingleSelectWithSize.SetLabel("SelectList With Size")
	p.SingleSelectWithSize.SetAttribute("size", 4)
	p.SingleSelectWithSize.AddListItems(itemList)

	p.RadioList1 = NewRadioList(p, "radioList1")
	p.RadioList1.SetLabel("Rows Radio List")
	p.RadioList1.AddListItems(itemList)
	p.RadioList1.SetColumnCount(2)

	p.RadioList2 = NewRadioList(p, "radioList2")
	p.RadioList2.SetLabel("Columns Radio List")
	p.RadioList2.AddListItems(itemList)
	p.RadioList2.SetColumnCount(2)
	p.RadioList2.SetDirection(LayoutColumn)

	p.RadioList3 = NewRadioList(p, "radioList3")
	p.RadioList3.SetLabel("Scrolling Radio List")
	p.RadioList3.AddListItems(itemList)
	p.RadioList3.SetIsScrolling(true)
	p.RadioList3.SetHeightStyle(80) // Limit the height to see the scrolling effect

	p.MultiSelect = NewMultiselectList(p, "multiselectList")
	p.MultiSelect.SetLabel("Multiselect List")
	p.MultiSelect.AddListItems(itemList)
	p.MultiSelect.SetIsRequired(true)

	p.CheckboxList1 = NewCheckboxList(p, "checklist1")
	p.CheckboxList1.SetLabel("Checkbox List")
	p.CheckboxList1.AddListItems(itemList)
	p.CheckboxList1.SetColumnCount(2)

	// TODO: Make radio list settings into functions
	// TODO: Test dynamic data setting
	p.SubmitAjax = NewButton(p, "ajaxButton")
	p.SubmitAjax.SetText("Submit Ajax")
	p.SubmitAjax.OnSubmit(action.Ajax(p.ID(), ButtonSubmit))

	p.SubmitServer = NewButton(p, "serverButton")
	p.SubmitServer.SetText("Submit Server")
	p.SubmitServer.OnSubmit(action.Server(p.ID(), ButtonSubmit))

}

func (p *SelectListPanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case ButtonSubmit:
		p.CheckboxList1.SetInstructions(strings.Join(p.CheckboxList1.SelectedIds(), ","))
		p.CheckboxList1.Refresh()
	}
}


func init() {
	browsertest.RegisterTestFunction("Select List Ajax Submit", testSelectListAjaxSubmit)
	browsertest.RegisterTestFunction("Select List Server Submit", testSelectListServerSubmit)
}

// testPlain exercises the plain text box
func testSelectListAjaxSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).AddValue("control", "selectlist").AddValue("testing", 1).String()
	f := t.LoadUrl(myUrl)

	testSelectListSubmit(t, f, f.Page().GetControl("ajaxButton"))

	t.Done("Complete")
}

func testSelectListServerSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).AddValue("control", "selectlist").AddValue("testing", 1).String()
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
	select1 := f.Page().GetControl("singleSelectList").(*SelectList)
	select2 := f.Page().GetControl("selectListWithSize").(*SelectList)
	radio1 := f.Page().GetControl("radioList1").(*RadioList)
	radio2 := f.Page().GetControl("radioList2").(*RadioList)

	id,_ := select2.GetItemByValue(2)
	t.ChangeVal("selectListWithSize", id)

	t.Click(btn)

	t.AssertEqual(true, t.HasClass("singleSelectList_ctl", "error"))


	t.AssertEqual(2, select2.IntValue())

	id,_ = select1.GetItemByValue(1)
	t.ChangeVal("singleSelectList", id)
	id,_ = select2.GetItemByValue(2)
	t.ChangeVal("selectListWithSize", id)
	id,_ = radio1.GetItemByValue(3)
	t.CheckGroup("radioList1", id)
	id,_ = radio2.GetItemByValue(4)
	t.CheckGroup("radioList2", id)

	t.Click(btn)

	t.AssertEqual(1, select1.IntValue())
	t.AssertEqual(2, select2.IntValue())
	t.AssertEqual(3, radio1.IntValue())
	t.AssertEqual(4, radio2.IntValue())
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