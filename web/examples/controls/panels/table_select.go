package panels

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/column"
	"github.com/goradd/goradd/pkg/page/control/data"
	"strconv"
)

const (
	rowSelectedEvent = iota + 1
)

type TableSelectPanel struct {
	Panel
}


func NewTableSelectPanel(ctx context.Context, parent page.ControlI) {
	// In this first example, we create a small table pre-filled with data
	var items = []map[string]string {
		{"id": "1", "col1": "Row 1, Col 1", "col2": "Row 1, Col 2"},
		{"id": "2", "col1": "Row 2, Col 1", "col2": "Row 2, Col 2"},
		{"id": "3", "col1": "Row 3, Col 1", "col2": "Row 3, Col 2"},
	}

	p := &TableSelectPanel{}
	p.Panel.Init(p, parent, "tableSelectPanel")
	p.AddControls(ctx,
		SelectTableCreator{
			ID: "table1",
			Columns:[]ColumnCreator {
				column.MapColumnCreator{
					Index: "col1",
				},
				column.MapColumnCreator{
					Index: "col2",
				},
			},
			ControlOptions: page.ControlOptions{
				Class: "gr-table-rows",
			},
			SaveState: true,
			OnRowSelected: action.Ajax(p.ID(), rowSelectedEvent),
			Data: items,
		},
		SelectTableCreator{
			ID: "table2",
			DataProvider: p, // The data provider can be a predefined control, including the parent of the table.
			Columns:[]ColumnCreator {
				column.MapColumnCreator{
					Index: "col1",
				},
				column.MapColumnCreator{
					Index: "col2",
				},
			},
			ControlOptions: page.ControlOptions{
				Class: "gr-table-rows",
				DataAttributes: page.DataAttributeMap {
					"grOptScrollable": true, // make it scrollable
				},
			},
			SaveState: true,
			OnRowSelected: action.Ajax(p.ID(), rowSelectedEvent),
		},
		PanelCreator{
			ID: "infoPanel",
		},
		ButtonCreator{
			ID:       "showButton",
			Text:     "Show Selected Item",
			OnClick: action.Javascript("g$('table2').showSelectedItem()"),
		},
	)

	// In this example, the lines below only take affect when you refresh the page.
	// The above SaveState call not only tells the table to remember its selection, but it also serves
	// as the place for the table to recall a previously saved selection. So, SelectedID() will return
	// a remembered value after SaveState is called above.
	if t2 := GetSelectTable(p, "table2"); t2.SelectedID() != "" {
		GetPanel(p, "infoPanel").SetText(fmt.Sprintf("Row %s was selected.", t2.SelectedID()))
	}
}

// BindData satisfies the data provider interface so that the parent panel of the table
// is the one that is providing the table.
func (p *TableSelectPanel) BindData(ctx context.Context, s data.DataManagerI) {
	var items []map[string]string
	for i := 0; i < 50; i++ {
		item := map[string]string{"id": strconv.Itoa(i), "col1": fmt.Sprintf("Row %d, Col 0", i), "col2":fmt.Sprintf("Row %d, Col 1", i)}
		items = append(items, item)
	}

	s.(*SelectTable).SetData(items)
}


func (p *TableSelectPanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case rowSelectedEvent:
		rowID := a.EventValueString()
		GetPanel(p, "infoPanel").SetText(fmt.Sprintf("Row %s was selected.", rowID))
	}
}


func init() {
}
