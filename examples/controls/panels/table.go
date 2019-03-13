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
)

type TablePanel struct {
	Panel

	Table1	*PaginatedTable
	Pager1 *DataPager
	SelectCol *column.CheckboxColumn

	SubmitAjax      *Button
	SubmitServer    *Button
}

type Table1Data map[string]string

// This sample data is in the form of a slice of maps. Typically you would not do this, but
// some special situations may find this approach useful.
var table1Data = []Table1Data {
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

type SelectedProvider struct{
	column.DefaultCheckboxProvider
}

func (c SelectedProvider) ID(data interface{}) string {
	return data.(Table1Data)["id"]
}

func (c SelectedProvider) IsChecked(data interface{}) bool {
	return data.(Table1Data)["s"] == "1"
}

func NewTablePanel(parent page.ControlI) *TablePanel {
	p := &TablePanel{}
	p.Panel.Init(p, parent, "checkboxPanel")

	p.Table1 = NewPaginatedTable(p, "table1")
	p.Table1.SetHeaderRowCount(1)
	p.Table1.SetDataProvider(p)
	p.Table1.AddColumn(column.NewMapColumn("name").SetTitle("Name"))
	p.Table1.AddColumn(column.NewCheckboxColumn(SelectedProvider{}).SetTitle("Selected"))
	//p.Table1.AddColumn(column.NewCheckboxColumn(p).SetTitle("Completed"))

	p.Pager1 = NewDataPager(p, "", p.Table1)
	p.Table1.SetPageSize(5)

	p.SubmitAjax = NewButton(p, "ajaxButton")
	p.SubmitAjax.SetText("Submit Ajax")
	p.SubmitAjax.OnSubmit(action.Ajax(p.ID(), ButtonSubmit))

	p.SubmitServer = NewButton(p, "serverButton")
	p.SubmitServer.SetText("Submit Server")
	p.SubmitServer.OnSubmit(action.Server(p.ID(), ButtonSubmit))

	return p
}

// BindData satisfies the data provider interface so that the parent panel of the table
// is the one that is providing the table.
func (f *TablePanel) BindData(ctx context.Context, s data.DataManagerI) {
	f.Table1.SetTotalItems(uint(len(table1Data)))
	start, end := f.Pager1.SliceOffsets()
	s.SetData(table1Data[start:end])
}


func (p *TablePanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case ButtonSubmit:
	}
}



func init() {
	//browsertest.RegisterTestFunction("Table Ajax Submit", testTableAjaxSubmit)
	//browsertest.RegisterTestFunction("Table Server Submit", testTableServerSubmit)
}

// testPlain exercises the plain text box
func testTableAjaxSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).AddValue("control", "checkbox").String()
	f := t.LoadUrl(myUrl)

	testTableSubmit(t, f, f.Page().GetControl("ajaxButton"))

	t.Done("Complete")
}

func testTableServerSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).AddValue("control", "checkbox").String()
	f := t.LoadUrl(myUrl)

	testTableSubmit(t, f, f.Page().GetControl("serverButton"))

	t.Done("Complete")
}

// testTableSubmit does a variety of submits using the given button. We use this to double check the various
// results we might get after a submission, as well as nsure that the ajax and server submits produce
// the same results.
func testTableSubmit(t *browsertest.TestForm, f page.FormI, btn page.ControlI) {
//	table1 := f.Page().GetControl("table1").(*Table)

	t.Click(btn)

}
