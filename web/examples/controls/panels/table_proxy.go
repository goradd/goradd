package panels

import (
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/crypt"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/column"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
	. "github.com/goradd/goradd/web/examples/model"
)

type TableProxyPanel struct {
	Panel
}

func NewTableProxyPanel(ctx context.Context, parent page.ControlI) {
	p := &TableProxyPanel{}
	p.Panel.Init(p, parent, "tableProxyPanel")
	p.AddControls(ctx,
		ProxyCreator{
			ID: "pxy",
			On: On {
				Event:  event.Click(),
				Action: action.Ajax(p.ID(), ProxyClick),
			},
		},
		PagedTableCreator{
			ID: "table1",
			DataProvider:p,
			Columns:[]ColumnCreator {
				column.TexterColumnCreator{
					Texter: p,
					ColumnOptions: ColumnOptions{
						IsHtml:true,
					},
				},
			},
			SaveState: true,
			Caption:DataPagerCreator{
				PagedControl: "table1",
			},
			PageSize:5,
		},
		ProjectPanelCreator{},
	)

	log.Debug("Proxy Table Created")

}

// BindData satisfies the data provider interface so that the parent panel of the table
// is the one that is providing the table.
func (p *TableProxyPanel) BindData(ctx context.Context, s DataManagerI) {
	t := GetPagedTable(p, "table1")
	t.SetTotalItems(QueryProjects(ctx).Count(ctx, false))

	projects := QueryProjects(ctx).
		Limit(t.SqlLimits()).
		Load(ctx)
	t.SetData(projects)

	log.Debug("Binding Data - ", projects)

}

func (p *TableProxyPanel) CellText(ctx context.Context, col ColumnI, rowNum int, colNum int, data interface{}) string {
	// Since we only have one custom column, we know what we are getting.
	project := data.(*Project)
	id := crypt.SessionEncryptUrlValue(ctx, project.ID()) // Since this is a database id, lets encrypt it for extra security

	// This is just to assign an id for click testing. You don't normally need to assign an id.
	attr := html.NewAttributes()
	attr.SetID("pxy" + project.ID())

	pxy := GetProxy(p, "pxy")
	v := pxy.LinkHtml(ctx, project.Name(),
		id,
		attr)
	return v
}

func (p *TableProxyPanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case ProxyClick:
		id := a.ControlValueString()
		id = crypt.SessionDecryptUrlValue(ctx, id)
		if id != "" {
			project := LoadProject(ctx, id)
			GetProjectPanel(p).SetProject(project)
		}
	}
}

type ProjectPanel struct {
	Panel
	project *Project
}

func NewProjectPanel(parent page.ControlI) *ProjectPanel {
	p := &ProjectPanel{}
	p.Init(p, parent, "personPanel")

	return p
}

func (p *ProjectPanel) SetProject(project *Project) {
	p.project = project
	p.Refresh()
}

func init() {
	browsertest.RegisterTestFunction("Table - Proxy Column", testTableProxyCol)
	gob.Register(ProjectPanel{})
	gob.Register(TableProxyPanel{})
}

func testTableProxyCol(t *browsertest.TestForm) {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "tableproxy").SetValue("testing", 1).String()
	t.LoadUrl(myUrl)

	t.ClickHtmlItem("pxy1")
/*	h := t.InnerHtml("nameItem")
	t.AssertEqual("<label>Name</label>ACME Website Redesign", h)
*/	t.Done("Complete")

}

// PanelCreator creates a div control with child controls.
// Pass it to AddControls or as a child of a parent control.
type ProjectPanelCreator struct {
}

// Create is called by the framework to create the panel. You do not normally need to call this.
func (c ProjectPanelCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewProjectPanel(parent)
	return ctrl
}

// GetProjectPanel is a convenience method to return the panel with the given id from the page.
func GetProjectPanel(c page.ControlI) *ProjectPanel {
	return c.Page().GetControl("personPanel").(*ProjectPanel)
}