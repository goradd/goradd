package panels

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/column"
)

type TablePanel struct {
	Panel
}

type TableMapData map[string]string
type TableSliceData []string

// Make the TableMapData satisfy the Getter interface so it can be used in a Getter column.
func (m TableMapData) Get(i string) string {
	return m[i]
}

var tableMapData = []TableMapData{
	{"id": "1", "name": "This"},
	{"id": "2", "name": "That"},
	{"id": "3", "name": "Other"},
	{"id": "4", "name": "Here"},
	{"id": "5", "name": "There"},
	{"id": "6", "name": "Everywhere"},
	{"id": "7", "name": "Over"},
	{"id": "8", "name": "Under"},
	{"id": "9", "name": "Near"},
	{"id": "10", "name": "Far"},
	{"id": "11", "name": "Who"},
	{"id": "12", "name": "What"},
	{"id": "13", "name": "Why"},
	{"id": "14", "name": "When"},
	{"id": "15", "name": "How"},
	{"id": "16", "name": "Which"},
	{"id": "17", "name": "If"},
	{"id": "18", "name": "Then"},
	{"id": "19", "name": "Or"},
	{"id": "20", "name": "And"},
	{"id": "21", "name": "But"},
}

var tableSliceData = []TableSliceData{
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
	{"19", "Or"},
	{"20", "And"},
	{"21", "But"},
}

func NewTablePanel(ctx context.Context, parent page.ControlI) {
	p := &TablePanel{}
	p.Self = p
	p.Init(ctx, parent, "tablePanel")
}

func (p *TablePanel) Init(ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(parent, id)

	p.AddControls(ctx,
		PagedTableCreator{
			ID:             "table1",
			HeaderRowCount: 1,
			DataProvider:   p,
			Columns: []ColumnCreator{
				column.TexterColumnCreator{
					TexterID: "tablePanel",
					Title:    "Custom",
				},
				column.SliceColumnCreator{
					Index: 1,
					Title: "Slice",
				},
			},
			PageSize: 5,
		},
		// A DataPager can be a standalone control, which you draw manually
		DataPagerCreator{
			ID:           "pager1",
			PagedControl: "table1",
		},
		PagedTableCreator{
			ID:             "table2",
			HeaderRowCount: 1,
			DataProviderID: "tablePanel", // You can use the id of a provider too
			Columns: []ColumnCreator{
				column.TexterColumnCreator{
					TexterID: "tablePanel",
					Title:    "Custom",
				},
				column.MapColumnCreator{
					Index: "id",
					Title: "StdMap",
				},
				column.GetterColumnCreator{
					Index: "name",
					Title: "Getter",
				},
			},
			PageSize: 5,
			// A DataPager can also be a caption, and will get drawn for you as part of the table
			Caption: DataPagerCreator{
				ID:           "pager2",
				PagedControl: "table2",
			},
		},
	)
}

// BindData satisfies the data provider interface so that the parent panel of the table
// is the one that is providing the table.
func (p *TablePanel) BindData(ctx context.Context, s DataManagerI) {
	switch s.ID() {
	case "table1":
		t := s.(PagedControlI)
		t.SetTotalItems(uint(len(tableSliceData)))
		start, end := t.SliceOffsets()
		s.SetData(tableSliceData[start:end])
	case "table2":
		t := s.(PagedControlI)
		t.SetTotalItems(uint(len(tableMapData)))
		start, end := t.SliceOffsets()
		s.SetData(tableMapData[start:end])

	}
}

// CellText here satisfies the CellTexter interface so that the panel can provide the text for a cell.
func (p *TablePanel) CellText(ctx context.Context, col ColumnI, info CellInfo) string {
	// Here is an example of how to figure out what table we are talking about.
	tid := col.ParentTable().ID()
	switch tid {
	case "table1":
		return fmt.Sprintf("Id: %s, Row #%d, Col #%d", info.Data.(TableSliceData)[0], info.RowNum, info.ColNum)
	case "table2":
		return fmt.Sprintf("Id: %s, Row #%d, Col #%d", info.Data.(TableMapData)["id"], info.RowNum, info.ColNum)
	}
	return ""
}

func init() {
	page.RegisterControl(&TablePanel{})
}
