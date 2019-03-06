package panels

import (
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
)

type SelectListPanel struct {
	Panel
	SingleSelect   *SelectList
	SingleSelectWithSize   *SelectList
	RadioList1   *RadioList
	RadioList2   *RadioList
	RadioList3   *RadioList

	SubmitAjax      *Button
	SubmitServer    *Button
}

func NewSelectListPanel(parent page.ControlI) *SelectListPanel {
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
	p.RadioList1.ColumnCount = 2

	p.RadioList2 = NewRadioList(p, "radioList2")
	p.RadioList2.SetLabel("Columns Radio List")
	p.RadioList2.AddListItems(itemList)
	p.RadioList2.ColumnCount = 2
	p.RadioList2.Placement = NextItemCrossAxis

	p.RadioList3 = NewRadioList(p, "radioList3")
	p.RadioList3.SetLabel("Scrolling Radio List")
	p.RadioList3.AddListItems(itemList)
	p.RadioList3.IsScrolling = true
	p.RadioList3.SetHeightStyle(80) // Limit the height to see the scrolling effect


	p.SubmitAjax = NewButton(p, "ajaxButton")
	p.SubmitAjax.SetText("Submit Ajax")
	p.SubmitAjax.OnSubmit(action.Ajax(p.ID(), AjaxSubmit))

	p.SubmitServer = NewButton(p, "serverButton")
	p.SubmitServer.SetText("Submit Server")
	p.SubmitServer.OnSubmit(action.Server(p.ID(), ServerSubmit))

	return p
}


func init() {
	browsertest.RegisterTestFunction("Select List Ajax Submit", testSelectListAjaxSubmit)
	browsertest.RegisterTestFunction("Select List Server Submit", testSelectListServerSubmit)
}

// testPlain exercises the plain text box
func testSelectListAjaxSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).AddValue("control", "selectlist").String()
	f := t.LoadUrl(myUrl)

	testSelectListSubmit(t, f, "ajaxButton")

	t.Done("Complete")
}

func testSelectListServerSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).AddValue("control", "selectlist").String()
	f := t.LoadUrl(myUrl)

	testSelectListSubmit(t, f, "serverButton")

	t.Done("Complete")
}

// testCheckboxSubmit does a variety of submits using the given button. We use this to double check the various
// results we might get after a submission, as well as nsure that the ajax and server submits produce
// the same results.
func testSelectListSubmit(t *browsertest.TestForm, f page.FormI, btn string) {

	t.ChangeVal("selectListWithSize", 2)

	t.Click(btn)

	t.AssertEqual(true, t.HasClass("selectlist_ctl", "error"))

	select1 := f.Page().GetControl("singleSelectList").(*SelectList)
	select2 := f.Page().GetControl("selectListWithSize").(*SelectList)
	radio1 := f.Page().GetControl("radioList1").(*RadioList)
	radio2 := f.Page().GetControl("radioList2").(*RadioList)

	t.AssertEqual(2, select2.IntValue())

	t.ChangeVal("selectList", 1)
	t.ChangeVal("selectListWithSize", 2)
	t.ChangeVal("radioList1", 3)
	t.ChangeVal("radioList2", 4)

	t.Click(btn)

	t.AssertEqual(1, select1.IntValue())
	t.AssertEqual(2, select2.IntValue())
	t.AssertEqual(3, radio1.IntValue())
	t.AssertEqual(4, radio2.IntValue())
}

