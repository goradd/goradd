package panels

import (
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/orm/op"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/column"
	"github.com/goradd/goradd/web/examples/model"
	"github.com/goradd/goradd/web/examples/model/node"
)

type TableDbPanel struct {
	Panel
}

func NewTableDbPanel(ctx context.Context, parent page.ControlI) {
	p := &TableDbPanel{}
	p.Panel.Init(p, parent, "tableDbPanel")
	p.AddControls(ctx,
		PagedTableCreator{
			ID: "table1",
			HeaderRowCount: 1,
			DataProvider: p, // The data provider can be a predefined control, including the parent of the table.
			Sortable: true,
			Columns:[]ColumnCreator {
				column.NodeColumnCreator{
					Node: node.Person().FirstName(),
					Title:"First Name",
					Sortable:true,
				},
				column.NodeColumnCreator{
					Node: node.Person().LastName(),
					Title:"Last Name",
					Sortable:true,
				},

				column.TexterColumnCreator{
					Texter: p,
					Title:"Combined",
				},
				column.AliasColumnCreator{
					Alias: "manager_count",
					Title: "Project Count",
				},
			},
			PageSize:5,
			SaveState: true,
			Caption:DataPagerCreator{
				ID:           "pager",
				PagedControl: "table1",
			},
			SortHistoryLimit: 3,
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
func (p *TableDbPanel) BindData(ctx context.Context, s DataManagerI) {
	t := s.(*PagedTable)
	t.SetTotalItems(model.QueryPeople(ctx).Count(ctx, false))

	// figure out how to sort the columns. This could be a simple process, or complex, depending on your data

	// Since we are asking the database to do the sort, we have to make a slice of nodes
	sortNodes := column.MakeNodeSlice(t.SortColumns())

	people := model.QueryPeople(ctx).
		Alias("manager_count",
			model.QueryProjects(ctx).
				Alias("", op.Count(node.Project().ManagerID())).
				Where(op.Equal(node.Project().ManagerID(), node.Person().ID())).
				Subquery()).
		Limit(t.SqlLimits()).
		OrderBy(sortNodes...).
		Load(ctx)
	t.SetData(people)
}

func (p *TableDbPanel) CellText(ctx context.Context, col ColumnI, rowNum int, colNum int, data interface{}) string {
	// Since we only have one custom column, we know what we are getting.
	person := data.(*model.Person)
	return person.FirstName() + " " + person.LastName()
}

func init() {
	gob.Register(TableDbPanel{})
}
