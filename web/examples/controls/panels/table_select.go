package panels

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/column"
	"github.com/goradd/goradd/pkg/page/control/data"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/goradd/web/examples/controls"
	"strconv"
)

const (
	rowSelectedEvent = iota + 1
)

type TableSelectPanel struct {
	Panel

	Table1	*SelectTable
	InfoPanel *Panel
	ShowButton *Button
}


func NewTableSelectPanel(ctx context.Context, parent page.ControlI) {
	p := &TableSelectPanel{}
	p.Panel.Init(p, parent, "tableSelectPanel")

	p.Table1 = NewSelectTable(p, "table1")
	p.Table1.SetDataProvider(p)
	p.Table1.AddColumn(column.NewMapColumn("col1"))
	p.Table1.AddColumn(column.NewMapColumn("col2"))
	p.Table1.AddClass("gr-table-rows") // Add default table styling
	p.Table1.SaveState(ctx, true)
	p.Table1.On(event.RowSelected(), action.Ajax(p.ID(), rowSelectedEvent))

	p.InfoPanel = NewPanel(p, "infoPanel")

	// In this example, the lines below only take affect when you refresh the page.
	// The above SaveState call not only tells the table to remember its selection, but it also serves
	// as the place for the table to recall a previously saved selection. So, SelectedID() will return
	// a remembered value after SaveState is called above.
	if p.Table1.SelectedID() != "" {
		p.InfoPanel.SetText(fmt.Sprintf("Row %s was selected.", p.Table1.SelectedID()))
	}

	p.ShowButton = NewButton(p, "")
	p.ShowButton.SetLabel("Show Selected Item")
	p.ShowButton.On(event.Click(), action.Javascript("g$('table1').showSelectedItem()"))
}

// BindData satisfies the data provider interface so that the parent panel of the table
// is the one that is providing the table.
func (p *TableSelectPanel) BindData(ctx context.Context, s data.DataManagerI) {
	var items []map[string]string
	for i := 0; i < 50; i++ {
		item := map[string]string{"id": strconv.Itoa(i), "col1": fmt.Sprintf("Row %d, Col 0", i), "col2":fmt.Sprintf("Row %d, Col 1", i)}
		items = append(items, item)
	}

	p.Table1.SetData(items)
}


func (p *TableSelectPanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case rowSelectedEvent:
		rowID := a.EventValueString()
		p.InfoPanel.SetText(fmt.Sprintf("Row %s was selected.", rowID))
	}
}


func init() {
	controls.RegisterPanel("tableselect", "Tables - Select Row", NewTableSelectPanel, 10)
}
