package panels

import (
	"context"
	"github.com/goradd/goradd/examples/model"
	"github.com/goradd/goradd/examples/model/node"
	"github.com/goradd/goradd/pkg/orm/op"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/column"
	"github.com/goradd/goradd/pkg/page/control/data"
)

type TableDbPanel struct {
	Panel

	Table1	*PaginatedTable
	Pager1 *DataPager
}


func NewTableDbPanel(ctx context.Context, parent page.ControlI) {
	p := &TableDbPanel{}
	p.Panel.Init(p, parent, "tableDbPanel")

	p.Table1 = NewPaginatedTable(p, "table1")
	p.Table1.SetHeaderRowCount(1)
	p.Table1.SetDataProvider(p)
	p.Table1.SetSortable()
	p.Table1.AddColumn(column.NewNodeColumn(node.Person().FirstName()).
		SetTitle("First Name").
		SetSortable())
	p.Table1.AddColumn(column.NewNodeColumn(node.Person().LastName()).
		SetTitle("Last Name").
		SetSortable())
	p.Table1.AddColumn(column.NewCustomColumn(p).SetTitle("Combined"))
	p.Table1.AddColumn(column.NewAliasColumn("manager_count").SetTitle("Project Count"))
	p.Pager1 = NewDataPager(p, "pager1", p.Table1)
	p.Table1.SetPageSize(5)
	p.Table1.SetSortHistoryLimit(3) // Sets the depth of how many sort columns we will track
}

// BindData satisfies the data provider interface so that the parent panel of the table
// is the one that is providing the table.
func (p *TableDbPanel) BindData(ctx context.Context, s data.DataManagerI) {
	p.Table1.SetTotalItems(model.QueryPeople().Count(ctx, false))

	// figure out how to sort the columns. This could be a simple process, or complex, depending on your data

	// Since we are asking the database to do the sort, we have to make a slice of nodes
	sortNodes := column.MakeNodeSlice(p.Table1.SortColumns())

	people := model.QueryPeople().
		Alias("manager_count",
			model.QueryProjects().
				Alias("", op.Count(node.Project().ManagerID())).
				Where(op.Equal(node.Project().ManagerID(), node.Person().ID())).
				Subquery()).
		Limit(p.Pager1.SqlLimits()).
		OrderBy(sortNodes...).
		Load(ctx)
	p.Table1.SetData(people)
}

func (f *TableDbPanel) 	CellText(ctx context.Context, col ColumnI, rowNum int, colNum int, data interface{}) string {
	// Since we only have one custom column, we know what we are getting.
	p := data.(*model.Person)
	return p.FirstName() + " " + p.LastName()
}

func init() {
}


