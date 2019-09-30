package panels

import (
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/column"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
	"strconv"
)

type TableCheckboxPanel struct {
	Panel

}

type Table1Data map[string]string

// This sample data is in the form of a slice of maps. Typically you would not do this, but
// some special situations may find this approach useful.
var table1Data = getCheckTestData()

type SelectedProvider struct {
	column.DefaultCheckboxProvider
}

func (c SelectedProvider) RowID(data interface{}) string {
	return data.(Table1Data)["id"]
}

func (c SelectedProvider) IsChecked(data interface{}) bool {
	if data == nil {
		return false // since we aren't keeping track, just assume not everything is checked
	}
	return data.(Table1Data)["s"] == "1"
}

func NewTableCheckboxPanel(ctx context.Context, parent page.ControlI) {
	p := &TableCheckboxPanel{}
	p.Panel.Init(p, parent, "checkboxTablePanel")
	p.AddControls(ctx,
		PagedTableCreator{
			ID: "table1",
			HeaderRowCount: 1,
			DataProvider: p,
			Columns:[]ColumnCreator {
				column.MapColumnCreator{
					Index: "name",
					Title:"Name",
				},
				column.CheckboxColumnCreator{
					ID:"check1",
					Title:"Selected",
					ShowCheckAll:true,
					CheckboxProvider:SelectedProvider{},
				},
			},
			PageSize:5,
			SaveState: true,
		},
		// A DataPager can be a standalone control, which you draw manually
		DataPagerCreator{
			ID:           "pager",
			PagedControl: "table1",
		},
		ButtonCreator{
			ID:       "ajaxButton",
			Text:     "Submit Ajax",
			OnSubmit: action.Ajax("checkboxPanel", ButtonSubmit),
		},
		ButtonCreator{
			ID:       "serverButton",
			Text:     "Submit Server",
			OnSubmit: action.Ajax("checkboxPanel", ButtonSubmit),
		},
	)
}

// BindData satisfies the data provider interface so that the parent panel of the table
// is the one that is providing the table.
func (f *TableCheckboxPanel) BindData(ctx context.Context, s DataManagerI) {
	t := s.(PagedControlI)
	t.SetTotalItems(uint(len(table1Data)))
	start, end := t.SliceOffsets()
	s.SetData(table1Data[start:end])
}

func (p *TableCheckboxPanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case ButtonSubmit:
		col := GetPagedTable(p, "table1").GetColumnByID("check1").(*column.CheckboxColumn)
		for k, v := range col.Changes() {
			i, _ := strconv.Atoi(k)
			var s string
			if v {
				s = "1"
			}
			table1Data[i-1]["s"] = s
		}
	}
}

func init() {
	browsertest.RegisterTestFunction("Table - Checkbox Nav", testTableCheckboxNav)
	browsertest.RegisterTestFunction("Table - Checkbox Ajax Submit", testTableCheckboxAjaxSubmit)
	browsertest.RegisterTestFunction("Table - Checkbox Server Submit", testTableCheckboxServerSubmit)

	gob.Register(SelectedProvider{}) // We must register this here because we are putting the changes map into the session,

}

func testTableCheckboxNav(t *browsertest.TestForm) {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "tablecheckbox").String()
	f := t.LoadUrl(myUrl)

	t.SetCheckbox("table1_check1_1", true)
	pager := f.Page().GetControl("pager").(*DataPager)
	table := f.Page().GetControl("table1").(*PagedTable)
	col := table.GetColumnByID("check1").(*column.CheckboxColumn)
	changes := col.Changes()
	_, ok := changes["1"]
	t.AssertEqual(false, ok)

	t.ClickSubItem(pager, "page_2")
	changes = col.Changes()
	changed, _ := changes["1"]
	t.AssertEqual(true, changed)

	// restore state for other tests
	t.ClickSubItem(pager, "page_1")
	t.SetCheckbox("table1_check1_1", false)
	t.ClickSubItem(pager, "page_1")

	t.Done("Complete")
}

func testTableCheckboxAjaxSubmit(t *browsertest.TestForm) {
	testTableCheckboxSubmit(t, "ajaxButton")

	t.Done("Complete")
}

func testTableCheckboxServerSubmit(t *browsertest.TestForm) {
	testTableCheckboxSubmit(t, "serverButton")

	t.Done("Complete")
}

func testTableCheckboxSubmit(t *browsertest.TestForm, btnName string) {
	table1Data = getCheckTestData()
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "tablecheckbox").SetValue("testing", 1).String()
	f := t.LoadUrl(myUrl)
	btn := f.Page().GetControl(btnName)

	t.SetCheckbox("table1_check1_1", true)
	table := f.Page().GetControl("table1").(*PagedTable)
	col := table.GetColumnByID("check1").(*column.CheckboxColumn)
	changes := col.Changes()
	_, ok := changes["1"]
	t.AssertEqual(false, ok)

	t.Click(btn)
	changes = col.Changes()
	changed, _ := changes["1"]
	t.AssertEqual(true, changed)

	// restore state for other tests
	t.SetCheckbox("table1_check1_1", false)
	t.Click(btn)

}

func getCheckTestData() []Table1Data {
	return []Table1Data{
		{"id": "1", "name": "This", "s": "", "c": "1"},
		{"id": "2", "name": "That", "s": "1", "c": ""},
		{"id": "3", "name": "Other", "s": "", "c": ""},
		{"id": "4", "name": "Here", "s": "", "c": ""},
		{"id": "5", "name": "There", "s": "", "c": ""},
		{"id": "6", "name": "Everywhere", "s": "", "c": ""},
		{"id": "7", "name": "Over", "s": "", "c": ""},
		{"id": "8", "name": "Under", "s": "", "c": ""},
		{"id": "9", "name": "Near", "s": "", "c": ""},
		{"id": "10", "name": "Far", "s": "", "c": ""},
		{"id": "11", "name": "Who", "s": "", "c": ""},
		{"id": "12", "name": "What", "s": "", "c": ""},
		{"id": "13", "name": "Why", "s": "", "c": ""},
		{"id": "14", "name": "When", "s": "", "c": ""},
		{"id": "15", "name": "How", "s": "", "c": ""},
		{"id": "16", "name": "Which", "s": "", "c": ""},
		{"id": "17", "name": "If", "s": "", "c": ""},
		{"id": "18", "name": "Then", "s": "", "c": ""},
		{"id": "19", "name": "Or", "s": "", "c": ""},
		{"id": "20", "name": "And", "s": "", "c": "1"},
		{"id": "21", "name": "But", "s": "1", "c": ""},
	}
}
