package panels

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/column"
	"github.com/goradd/goradd/pkg/page/control/data"
)

type TablePanel struct {
	Panel

	Table1	*PaginatedTable
	Pager1 *DataPager
	Column1 *column.SliceColumn
	Column2 *column.CustomColumn

	Table2	*PaginatedTable
	Pager2 *DataPager
	Column3 *column.MapColumn
	Column4 *column.GetterColumn
}

type TableMapData map[string]string
type TableSliceData []string

// Make the TableMapData satisfy the Getter interface so it can be used in a Getter column.
func (m TableMapData) Get(i string) string {
	return m[i]
}

var tableMapData = []TableMapData {
	{"id":"1", "name":"This"},
	{"id":"2", "name":"That"},
	{"id":"3", "name":"Other"},
	{"id":"4", "name":"Here"},
	{"id":"5", "name":"There"},
	{"id":"6", "name":"Everywhere"},
	{"id":"7", "name":"Over"},
	{"id":"8", "name":"Under"},
	{"id":"9", "name":"Near"},
	{"id":"10", "name":"Far"},
	{"id":"11", "name":"Who"},
	{"id":"12", "name":"What"},
	{"id":"13", "name":"Why"},
	{"id":"14", "name":"When"},
	{"id":"15", "name":"How"},
	{"id":"16", "name":"Which"},
	{"id":"17", "name":"If"},
	{"id":"18", "name":"Then"},
	{"id":"19", "name":"Or"},
	{"id":"20", "name":"And"},
	{"id":"21", "name":"But"},
}

var tableSliceData = []TableSliceData {
	{"1", "This"},
	{"2", "That"},
	{"3", "Other"},
	{"4", "Here"},
	{"5", "There"},
	{"6", "Everywhere"},
	{"7", "Over"},
	{"8", "Under"},
	{"9", "Near"},
	{"10", "Far"},
	{"11", "Who"},
	{"12", "What"},
	{"13", "Why"},
	{"14", "When"},
	{"15", "How"},
	{"16", "Which"},
	{"17", "If"},
	{"18", "Then"},
	{"19","Or"},
	{"20", "And"},
	{"21", "But"},
}


func NewTablePanel(ctx context.Context, parent page.ControlI) {
	p := &TablePanel{}
	p.Panel.Init(p, parent, "tablePanel")

	// The two tables here just demonstrate a variety of columns available to use in a data table.
	// Be sure to consider the Node column and Alias column which are not listed below, as they work directly with databases.
	p.Table1 = NewPaginatedTable(p, "table1")
	p.Table1.SetHeaderRowCount(1)
	p.Table1.SetDataProvider(p)
	p.Table1.AddColumn(column.NewCustomColumn(p).SetTitle("Custom"))
	p.Table1.AddColumn(column.NewSliceColumn(1).SetTitle("Slice"))
	p.Pager1 = NewDataPager(p, "pager1", p.Table1)
	p.Table1.SetPageSize(5)

	p.Table2 = NewPaginatedTable(p, "table2")
	p.Table2.SetHeaderRowCount(1)
	p.Table2.SetDataProvider(p)
	p.Table2.AddColumn(column.NewCustomColumn(p).SetTitle("Custom"))
	p.Table2.AddColumn(column.NewMapColumn("id").SetTitle("Map"))
	p.Table2.AddColumn(column.NewGetterColumn("name").SetTitle("Getter"))

	// The lines below put the pager into the caption of the table. The caption of a table can accept text or
	// a data pager type object. Note that the parent of the DataPager is the table, NOT the form, and the form
	// template does NOT draw the pager.
	p.Pager2 = NewDataPager(p.Table2, "pager2", p.Table2)
	p.Table2.SetCaption(p.Pager2)
	p.Table2.SetPageSize(5)

}

// BindData satisfies the data provider interface so that the parent panel of the table
// is the one that is providing the table.
func (f *TablePanel) BindData(ctx context.Context, s data.DataManagerI) {
	switch s.ID() {
	case "table1":
		f.Table1.SetTotalItems(uint(len(tableSliceData)))
		start, end := f.Pager1.SliceOffsets()
		s.SetData(tableSliceData[start:end])
	case "table2":
		f.Table2.SetTotalItems(uint(len(tableMapData)))
		start, end := f.Pager2.SliceOffsets()
		s.SetData(tableMapData[start:end])

	}
}

// CellText here satisfies the CellTexter interface so that the panel can provide the text for a cell.
func (f *TablePanel) 	CellText(ctx context.Context, col ColumnI, rowNum int, colNum int, data interface{}) string {
	// Here is an example of how to figure out what table we are talking about.
	tid := col.ParentTable().ID()
	switch tid {
	case "table1":
		return fmt.Sprintf("Id: %s, Row #%d, Col #%d", data.(TableSliceData)[0], rowNum, colNum)
	case "table2":
		return fmt.Sprintf("Id: %s, Row #%d, Col #%d", data.(TableMapData)["id"], rowNum, colNum)
	}
	return""
}

func init() {
}


