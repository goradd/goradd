package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/column"
	"github.com/goradd/goradd/pkg/page/control/data"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
	"strconv"
)

type TableCheckboxPanel struct {
	Panel

	Table1	*PaginatedTable
	Pager1 *DataPager
	CheckboxColumn1 *column.CheckboxColumn
	SelectCol *column.CheckboxColumn

	SubmitAjax      *Button
	SubmitServer    *Button
}

type Table1Data map[string]string

// This sample data is in the form of a slice of maps. Typically you would not do this, but
// some special situations may find this approach useful.
var table1Data = getCheckTestData()

type SelectedProvider struct{
	column.DefaultCheckboxProvider
}

func (c SelectedProvider) RowID(data interface{}) string {
	return data.(Table1Data)["id"]
}

func (c SelectedProvider) IsChecked(data interface{}) bool {
	return data.(Table1Data)["s"] == "1"
}

func NewTableCheckboxPanel(ctx context.Context, parent page.ControlI) {
	p := &TableCheckboxPanel{}
	p.Panel.Init(p, parent, "checkboxPanel")

	p.Table1 = NewPaginatedTable(p, "table1")
	p.Table1.SetHeaderRowCount(1)
	p.Table1.SetDataProvider(p)
	p.Table1.AddColumn(column.NewMapColumn("name").SetTitle("Name"))

	// get a copy of the column since we have to refer to it later
	p.CheckboxColumn1 = column.NewCheckboxColumn(SelectedProvider{})
	p.CheckboxColumn1.SetID("check1")
	p.CheckboxColumn1.SetTitle("Selected")
	p.Table1.AddColumn(p.CheckboxColumn1)
	//p.Table1.AddColumn(column.NewCheckboxColumn(p).SetTitle("Completed"))

	p.Pager1 = NewDataPager(p, "pager", p.Table1)
	p.Table1.SetPageSize(5)

	p.Table1.SaveState(ctx, true)

	p.SubmitAjax = NewButton(p, "ajaxButton")
	p.SubmitAjax.SetText("Submit Ajax")
	p.SubmitAjax.OnSubmit(action.Ajax(p.ID(), ButtonSubmit))

	p.SubmitServer = NewButton(p, "serverButton")
	p.SubmitServer.SetText("Submit Server")
	p.SubmitServer.OnSubmit(action.Server(p.ID(), ButtonSubmit))

}

// BindData satisfies the data provider interface so that the parent panel of the table
// is the one that is providing the table.
func (f *TableCheckboxPanel) BindData(ctx context.Context, s data.DataManagerI) {
	f.Table1.SetTotalItems(uint(len(table1Data)))
	start, end := f.Pager1.SliceOffsets()
	s.SetData(table1Data[start:end])
}


func (p *TableCheckboxPanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case ButtonSubmit:
		for k,v := range p.CheckboxColumn1.Changes() {
			i,_ := strconv.Atoi(k)
			var s string
			if v {
				s = "1"
			}
			table1Data[i - 1]["s"] = s
		}
	}
}



func init() {
	browsertest.RegisterTestFunction("Table - Checkbox Nav", testTableCheckboxNav)
	browsertest.RegisterTestFunction("Table - Checkbox Ajax Submit", testTableCheckboxAjaxSubmit)
	browsertest.RegisterTestFunction("Table - Checkbox Server Submit", testTableCheckboxServerSubmit)
}


func testTableCheckboxNav(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).AddValue("control", "tablecheckbox").String()
	f := t.LoadUrl(myUrl)

	t.SetCheckbox("table1_check1_1", true)
	pager := f.Page().GetControl("pager").(*DataPager)
	table := f.Page().GetControl("table1").(*PaginatedTable)
	col := table.GetColumnByID("check1").(*column.CheckboxColumn)
	changes := col.Changes()
	_,ok := changes["1"]
	t.AssertEqual(false, ok)

	t.ClickSubItem(pager, "page_2")
	changes = col.Changes()
	changed,_ := changes["1"]
	t.AssertEqual(true, changed)

	// restore state for other tests
	t.ClickSubItem(pager, "page_1")
	t.SetCheckbox("table1_check1_1", false)
	t.ClickSubItem(pager, "page_1")


	t.Done("Complete")
}


func testTableCheckboxAjaxSubmit(t *browsertest.TestForm)  {
	testTableCheckboxSubmit(t, "ajaxButton", )

	t.Done("Complete")
}

func testTableCheckboxServerSubmit(t *browsertest.TestForm)  {
	testTableCheckboxSubmit(t, "serverButton")

	t.Done("Complete")
}

func testTableCheckboxSubmit(t *browsertest.TestForm, btnName string) {
	table1Data = getCheckTestData()
	var myUrl = url.NewBuilder(controlsFormPath).AddValue("control", "tablecheckbox").String()
	f := t.LoadUrl(myUrl)
	btn := f.Page().GetControl(btnName)

	t.SetCheckbox("table1_check1_1", true)
	table := f.Page().GetControl("table1").(*PaginatedTable)
	col := table.GetColumnByID("check1").(*column.CheckboxColumn)
	changes := col.Changes()
	_,ok := changes["1"]
	t.AssertEqual(false, ok)

	t.Click(btn)
	changes = col.Changes()
	changed,_ := changes["1"]
	t.AssertEqual(true, changed)

	// restore state for other tests
	t.SetCheckbox("table1_check1_1", false)
	t.Click(btn)

}

func getCheckTestData() []Table1Data {
	return []Table1Data {
		{"id":"1", "name":"This","s":"","c":"1"},
		{"id":"2", "name":"That","s":"1","c":""},
		{"id":"3", "name":"Other","s":"","c":""},
		{"id":"4", "name":"Here","s":"","c":""},
		{"id":"5", "name":"There","s":"","c":""},
		{"id":"6", "name":"Everywhere","s":"","c":""},
		{"id":"7", "name":"Over","s":"","c":""},
		{"id":"8", "name":"Under","s":"","c":""},
		{"id":"9", "name":"Near","s":"","c":""},
		{"id":"10", "name":"Far","s":"","c":""},
		{"id":"11", "name":"Who","s":"","c":""},
		{"id":"12", "name":"What","s":"","c":""},
		{"id":"13", "name":"Why","s":"","c":""},
		{"id":"14", "name":"When","s":"","c":""},
		{"id":"15", "name":"How","s":"","c":""},
		{"id":"16", "name":"Which","s":"","c":""},
		{"id":"17", "name":"If","s":"","c":""},
		{"id":"18", "name":"Then","s":"","c":""},
		{"id":"19", "name":"Or","s":"","c":""},
		{"id":"20", "name":"And","s":"","c":"1"},
		{"id":"21", "name":"But","s":"1","c":""},
	}
}