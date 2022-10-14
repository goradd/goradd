package widget

import (
	"context"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/html5tag"
	"path"
)

const (
	filterChanged = iota + 4000
)

type ItemListPanelI interface {
	PanelI
}

type ItemListPanel struct {
	Panel
	FilterPanel *Panel
	ScrollPanel *Panel
	ButtonPanel *Panel

	FilterText *Textbox
	ItemTable  *SelectTable
	NewButton  *Button
	EditButton *Button
}

func NewItemListPanel(parent page.ControlI, id string) *ItemListPanel {
	p := &ItemListPanel{}
	p.Self = p
	p.Init(parent, id)
	return p
}

func (p *ItemListPanel) Init(parent page.ControlI, id string) {
	p.Panel.Init(parent, id)
	p.ParentForm().AddStyleSheetFile(path.Join(config.AssetPrefix, "goradd", "css", "item-list-panel.css"), nil)
	p.FilterPanel = NewPanel(p, p.ID()+"-filter")
	p.ScrollPanel = NewPanel(p, p.ID()+"-scroller")
	p.ButtonPanel = NewPanel(p, p.ID()+"-btnpnl")

	p.FilterPanel.AddClass("filter")
	p.ScrollPanel.AddClass("scroller")
	p.ButtonPanel.AddClass("buttons")

	p.FilterText = NewTextbox(p.FilterPanel, p.ID()+"-filtertxt")
	p.FilterText.SetPlaceholder(p.ParentForm().GT("Search"))
	p.FilterText.SetType(TextboxTypeSearch)

	p.ItemTable = NewSelectTable(p.ScrollPanel, p.ID()+"-table")

	p.NewButton = NewButton(p.ButtonPanel, p.ID()+"-newbtn")
	p.NewButton.SetText(p.ParentForm().GT("New"))

	p.EditButton = NewButton(p.ButtonPanel, p.ID()+"-editbtn")
	p.EditButton.SetText(p.ParentForm().GT("Edit"))

	p.FilterText.On(event.Input().Delay(300), action.Ajax(p.ID(), filterChanged))
	p.FilterText.On(event.EnterKey().Terminating(), action.Ajax(p.ID(), filterChanged))

}

func (p *ItemListPanel) Load(ctx context.Context) {
	p.ItemTable.SaveState(ctx, true)
}

func (p *ItemListPanel) this() ItemListPanelI {
	return p.Self.(ItemListPanelI)
}

func (c *ItemListPanel) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := c.Panel.DrawingAttributes(ctx)
	a.SetData("grctl", "itemlistpnl")
	return a
}

func (f *ItemListPanel) DoAction(ctx context.Context, a action.Params) {
	switch a.ID {
	case filterChanged:
		f.ItemTable.Refresh() // TODO: Change this to some kind of data only refresh so that when control is redrawn the scroll position is maintained
	}
}
