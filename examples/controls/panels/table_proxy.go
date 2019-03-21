package panels

import (
	"context"
	"github.com/goradd/goradd/examples/model"
	"github.com/goradd/goradd/pkg/crypt"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/column"
	"github.com/goradd/goradd/pkg/page/control/data"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
)

type TableProxyPanel struct {
	Panel

	Table1	*PaginatedTable
	Pager1 *DataPager
	Pxy *Proxy
	ProjectPanel *ProjectPanel
}


func NewTableProxyPanel(ctx context.Context, parent page.ControlI) {
	p := &TableProxyPanel{}
	p.Panel.Init(p, parent, "tableProxyPanel")

	p.Pxy = NewProxy(p)
	p.Pxy.On(event.Click(), action.Ajax(p.ID(), ProxyClick))

	p.Table1 = NewPaginatedTable(p, "table1")
	p.Table1.SetDataProvider(p)
	p.Table1.AddColumn(column.NewCustomColumn(p).SetIsHtml(true))
	p.Pager1 = NewDataPager(p, "pager1", p.Table1)
	p.Table1.SetPageSize(5)

	p.ProjectPanel = NewProjectPanel(p)

	log.Debug("Proxy Table Created")

}

// BindData satisfies the data provider interface so that the parent panel of the table
// is the one that is providing the table.
func (p *TableProxyPanel) BindData(ctx context.Context, s data.DataManagerI) {
	p.Table1.SetTotalItems(model.QueryProjects().Count(ctx, false))

	projects := model.QueryProjects().
		Limit(p.Pager1.SqlLimits()).
		Load(ctx)
	p.Table1.SetData(projects)

	log.Debug("Binding Data - ", projects)

}

func (f *TableProxyPanel) 	CellText(ctx context.Context, col ColumnI, rowNum int, colNum int, data interface{}) string {
	// Since we only have one custom column, we know what we are getting.
	project := data.(*model.Project)
	id := crypt.SessionEncryptUrlValue(ctx, project.ID()) // Since this is a database id, lets encrypt it for extra security

	// This is just to assign an id for click testing. You don't normally need to assign an id.
	attr := html.NewAttributes()
	attr.SetID("pxy" + project.ID())

	v := f.Pxy.LinkHtml(ctx, project.Name(),
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
			project := model.LoadProject(ctx, id)
			p.ProjectPanel.SetProject(project)
		}
	}
}


type ProjectPanel struct {
	Panel
	project *model.Project
}

func NewProjectPanel(parent page.ControlI) *ProjectPanel {
	p := &ProjectPanel{}
	p.Init(p, parent, "personPanel")

	return p
}

func (p *ProjectPanel) SetProject(project *model.Project) {
	p.project = project
	p.Refresh()
}

func init() {
	browsertest.RegisterTestFunction("Table - Proxy Column", testTableProxyCol)
}

func testTableProxyCol(t *browsertest.TestForm) {
	var myUrl = url.NewBuilder(controlsFormPath).AddValue("control", "tableproxy").String()
	t.LoadUrl(myUrl)

	t.ClickHtmlItem("pxy1")
	h := t.InnerHtml("nameItem")
	t.AssertEqual("<label>Name</label>ACME Website Redesign", h)
	t.Done("Complete")

}
