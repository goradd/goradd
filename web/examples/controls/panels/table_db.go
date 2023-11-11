package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/orm/op"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	. "github.com/goradd/goradd/pkg/page/control/button"
	. "github.com/goradd/goradd/pkg/page/control/table"
	column2 "github.com/goradd/goradd/pkg/page/control/table/column"
	"github.com/goradd/goradd/web/examples/gen/goradd/model"
	"github.com/goradd/goradd/web/examples/gen/goradd/model/node"
	"strconv"
)

type TableDbPanel struct {
	Panel
}

func NewTableDbPanel(ctx context.Context, parent page.ControlI) {
	p := new(TableDbPanel)
	p.Init(p, ctx, parent, "tableDbPanel")
}

func (p *TableDbPanel) Init(self any, ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, id)
	p.AddControls(ctx,
		PagedTableCreator{
			ID:             "table1",
			HeaderRowCount: 1,
			DataProvider:   p, // The data provider can be a predefined control, including the parent of the table.
			Sortable:       true,
			Columns: []ColumnCreator{
				column2.TexterColumnCreator{
					ID:     "num",
					Texter: p,
					Title:  "#",
				},
				column2.NodeColumnCreator{
					Node:     node.Person().FirstName(),
					Title:    "First Name",
					Sortable: true,
				},
				column2.NodeColumnCreator{
					Node:     node.Person().LastName(),
					Title:    "Last Name",
					Sortable: true,
				},

				column2.TexterColumnCreator{
					ID:     "combined",
					Texter: p,
					Title:  "Combined",
				},
				column2.AliasColumnCreator{
					Alias: "manager_count",
					Title: "Project Count",
				},
			},
			PageSize:  5,
			SaveState: true,
			Caption: DataPagerCreator{
				ID:             "pager",
				PagedControlID: "table1",
			},
			SortHistoryLimit: 3,
		},
		ButtonCreator{
			ID:       "ajaxButton",
			Text:     "Submit Ajax",
			OnSubmit: action.Do("checkboxPanel", ButtonSubmit),
		},
		ButtonCreator{
			ID:       "serverButton",
			Text:     "Submit Server",
			OnSubmit: action.Do("checkboxPanel", ButtonSubmit),
		},
	)

}

// BindData satisfies the data provider interface so that the parent panel of the table
// is the one that is providing the table.
func (p *TableDbPanel) BindData(ctx context.Context, s DataManagerI) {
	t := s.(*PagedTable)
	t.SetTotalItems(model.QueryPeople(ctx).Count(false))

	// figure out how to sort the columns. This could be a simple process, or complex, depending on your data

	// Since we are asking the database to do the sort, we have to make a slice of nodes
	sortNodes := column2.MakeNodeSlice(t.SortColumns())
	maxRowCount, offset := t.SqlLimits()
	people := model.QueryPeople(ctx).
		Alias("manager_count",
			model.QueryProjects(ctx).
				Alias("", op.Count(node.Project().ManagerID())).
				Where(op.Equal(node.Project().ManagerID(), node.Person().ID())).
				Subquery()).
		Limit(maxRowCount, offset).
		OrderBy(sortNodes...).
		Load()
	t.SetDataWithOffset(people, offset)
}

func (p *TableDbPanel) CellText(ctx context.Context, col ColumnI, info CellInfo) string {
	switch col.ID() {
	case "num":
		return strconv.Itoa(info.RowNum + 1)
	case "combined":
		person := info.Data.(*model.Person)
		return person.FirstName() + " " + person.LastName()
	}
	return ""
}

func init() {
	page.RegisterControl(&TableDbPanel{})
}
